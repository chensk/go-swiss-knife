package lock

import (
	"code.byted.org/ocean/swiss-knife/lock/provider"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"testing"
	"time"
)

var (
	cnt int
)

func TestMysqlLock(t *testing.T) {
	db, err := gorm.Open(mysql.Open("root:root1234@tcp(127.0.0.1:3306)/dolphin_admin?charset=utf8mb4&parseTime=True&loc=Local"),
		&gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		t.Fatal(err)
	}
	d, _ := db.DB()
	d.SetMaxIdleConns(1)
	d.SetMaxOpenConns(100)
	total := 1000
	wg := sync.WaitGroup{}
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			loc, err := TryLockTask("test-lock", []Options{
				WithProvider(provider.NewMysqlLockProvider(db)), WithLockAtMost(10 * time.Millisecond),
			}...)
			if err != nil {
				t.Logf("unexpected: %s", err)
			} else {
				cnt += 2
				_ = loc.Unlock()
			}
		}()
	}
	wg.Wait()
	t.Log("cnt: ", cnt)
}
