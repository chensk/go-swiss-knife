package provider

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

type LockModel struct {
	// ID is primary key normally.
	ID uint `gorm:"column:id;index;primaryKey;autoIncrement"`
	// Name is identifier of the lock, user can lock or unlock with specified name.
	Name string `gorm:"column:name"`
	// LockUntil is optional. If specified, the lock with hold until LockUntil arrives.
	LockUntil DbTime `gorm:"column:lock_until"`
	// LockAt is optional. If specified,
	LockAt DbTime `gorm:"column:locked_at"`
	// LockBy records the host ip currently holds the lock
	LockBy string `gorm:"column:locked_by"`
}

func (LockModel) TableName() string {
	return "shedlock"
}

type DbTime time.Time

func (dt DbTime) MarshalText() (data []byte, err error) {
	t := time.Time(dt)
	data = []byte(t.Format("2006-01-02 15:04:05"))
	return
}

func (dt *DbTime) UnmarshalText(text []byte) (err error) {
	t := (*time.Time)(dt)
	*t, err = time.Parse("2006-01-02T15:04:05", string(text))
	return
}

func (dt *DbTime) Scan(value interface{}) error {
	t, ok := value.(time.Time)
	if !ok {
		return errors.New(fmt.Sprintf("value is not time type: %s", value))
	}

	*dt = DbTime(t)
	return nil
}

func (dt DbTime) Value() (driver.Value, error) {
	return time.Time(dt), nil
}
