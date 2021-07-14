package sessionlib

import (
	"context"
	"testing"
	"time"
)

func TestSession(t *testing.T) {
	var sid = "1234"
	s, err := CreateSession(func() string {
		return sid
	}, func(s string) {
		sid = s
	}, []SessionOptions{
		WithExpiration(2 * time.Second),
	})
	if err != nil {
		t.Fatal(err)
	}
	v, ok := s.Get("test_key")
	if !ok {
		t.Log("key not found")
	} else {
		t.Logf("value: %s", v.(string))
	}
	s, err = CreateSession(func() string {
		return sid
	}, func(s string) {
		sid = s
	}, []SessionOptions{
		WithExpiration(2 * time.Second),
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = s.Set("test_key", "test_value")
	if err := s.Save(context.Background()); err != nil {
		t.Fatal("fail to save session")
	}

	s, err = CreateSession(func() string {
		return sid
	}, func(s string) {
		sid = s
	}, []SessionOptions{
		WithExpiration(2 * time.Second),
	})
	if err != nil {
		t.Fatal(err)
	}
	v, ok = s.Get("test_key")
	if !ok {
		t.Log("key not found")
	} else {
		t.Logf("value: %s", v.(string))
	}

	time.Sleep(3 * time.Second)
	s, err = CreateSession(func() string {
		return sid
	}, func(s string) {
		sid = s
	}, []SessionOptions{
		WithExpiration(2 * time.Second),
	})
	if err != nil {
		t.Fatal(err)
	}
	v, ok = s.Get("test_key")
	if !ok {
		t.Log("key not found")
	} else {
		t.Logf("value: %s", v.(string))
	}
}
