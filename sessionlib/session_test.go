package sessionlib

import (
	"testing"
	"time"
)

func TestSession(t *testing.T) {
	var sessionId = "1234"
	session, err := CreateSession(func() string {
		return sessionId
	}, func(sid string) {
		sessionId = sid
	}, []SessionOptions{
		WithRedisPsm("bytedance.redis.dolphin"),
		WithRedisTimeout(1 * time.Second),
	})

	if err != nil {
		t.Fatal(err)
	}

	v, ok := session.Get("test_key")
	if !ok {
		t.Error("key not exist")
	} else {
		t.Logf("value: %s", v)
	}

	session.Set("test_key", "test_value")
	//_ = session.Save(context.Background())

	session, err = CreateSession(func() string {
		return sessionId
	}, func(sid string) {
		sessionId = sid
	}, []SessionOptions{
		WithRedisPsm("bytedance.redis.dolphin"),
		WithRedisTimeout(1 * time.Second),
	})
	if err != nil {
		t.Fatal(err)
	}
	v, ok = session.Get("test_key")
	if !ok {
		t.Error("key not exist")
	} else {
		t.Logf("value: %s", v)
	}
}
