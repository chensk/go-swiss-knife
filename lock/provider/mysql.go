/*
provides simple distributed lock implemented by Mysql, which depends on Mysql table created in advance.
User can custom table name, but should keep the same columns name as follows:
*/
package provider

import (
	"code.byted.org/ocean/go-swiss-knife/lock"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

func NewMysqlLockProvider(db *gorm.DB) lock.LockProvider {
	return mysqlLockProvider{db: db}
}

type mysqlLockProvider struct {
	db *gorm.DB
}

type mysqlLock struct {
	name string
	db   *gorm.DB
}

func (m mysqlLock) Unlock() error {
	exam := LockModel{
		LockUntil: DbTime(time.Now()),
	}
	if err := m.db.Model(exam).Where("name = ?", m.name).Updates(exam).Error; err != nil {
		return fmt.Errorf("fail to unlock: %w", err)
	}
	return nil
}

func (m mysqlLockProvider) Lock(conf lock.Configuration) (lock.Lock, error) {
	if conf.Name == "" {
		return nil, ErrEmptyName
	}
	if conf.LockAtMost == 0 {
		return nil, lock.ErrLockAtMostMissing
	}
	var loc LockModel
	err := m.db.Model(LockModel{}).
		Where("name = ? and lock_until >= ?", conf.Name, time.Now()).First(&loc).Error
	timeoutWatcher := context.Background()
	if conf.LockTimeout > 0 {
		tw, f := context.WithTimeout(context.Background(), conf.LockTimeout)
		defer f()
		timeoutWatcher = tw
	}
	for {
		for err == nil {
			// keep waiting until lock is released
			elapse := time.Time(loc.LockUntil).Sub(time.Now())

			select {
			case <-time.After(elapse):
				err = m.db.Model(LockModel{}).
					Where("name = ? and lock_until >= ?", conf.Name, time.Now()).First(&loc).Error
			case <-timeoutWatcher.Done():
				return nil, lock.ErrTimeout
			}
		}

		var lb []byte

		lu := DbTime(time.Now().Add(conf.LockAtMost))
		lb, _ = lu.MarshalText()

		var e error
		if lockRegistry[conf.Name] {
			element := map[string]interface{}{
				"locked_at": DbTime(time.Now()), "locked_by": conf.LockBy, "lock_until": string(lb),
			}
			if res := m.db.Model(LockModel{}).Where("name = ? and lock_until < ?", conf.Name, time.Now()).Updates(element); res.Error != nil || res.RowsAffected == 0 {
				e = errors.New("")
			}
		} else {
			lockRegistry[conf.Name] = true
			element := map[string]interface{}{
				"name": conf.Name, "locked_at": DbTime(time.Now()), "locked_by": conf.LockBy, "lock_until": string(lb),
			}
			e = m.db.Model(LockModel{}).Create(element).Error
		}

		if e != nil {
			// if insert failed, meaning perhaps another goroutine get the lock first, keep waiting
			err = m.db.Model(LockModel{}).
				Where("name = ? and lock_until >= ?", conf.Name, time.Now()).First(&loc).Error
			continue
		} else {
			return mysqlLock{name: conf.Name, db: m.db}, nil
		}
	}
}

func (m mysqlLockProvider) TryLock(conf lock.Configuration) (lock.Lock, error) {
	if conf.Name == "" {
		return nil, ErrEmptyName
	}
	if conf.LockAtMost == 0 {
		return nil, lock.ErrLockAtMostMissing
	}
	element := map[string]interface{}{
		"name":       conf.Name,
		"locked_at":  DbTime(time.Now()),
		"locked_by":  conf.LockBy,
		"lock_until": DbTime(time.Now().Add(conf.LockAtMost)),
	}

	if lockRegistry[conf.Name] {
		if res := m.db.Model(LockModel{}).Where("name = ? and lock_until < ?", conf.Name, time.Now()).Updates(element); res.Error != nil || res.RowsAffected == 0 {
			return nil, lock.ErrLockFailed
		}
		return mysqlLock{name: conf.Name, db: m.db}, nil
	}
	lockRegistry[conf.Name] = true
	if err := m.db.Model(LockModel{}).Create(element).Error; err != nil {
		return nil, lock.ErrLockFailed
	}
	return mysqlLock{name: conf.Name, db: m.db}, nil
}

var (
	ErrEmptyName = errors.New("lock name empty")
	lockRegistry = make(map[string]bool)
)
