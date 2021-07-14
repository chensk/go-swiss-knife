package sessionlib

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"sync"
	"time"
)

// custom session id getter, such as fetched from http cookies.
type SessionIdGetter func() string

// custom session id setter, such as storing in http cookies.
type SessionIdSetter func(string)

// CreateSession gets the session store associated with sessionId which can be extracted by SessionIdGetter.
// If the session doesn't exist, create a new one and set the sessionId with SessionIdSetter.
// You can specify the store strategy by specify session options. By default, if WithRedisClusters option specifies
// the redis cluster, RedisSessionStore is used to store the session. Otherwise, InMemorySessionStore is used.
// You can implement your custom store strategy and specify it by WithSessionStore options.
//
// sessionIdGetter may not be nil, while sessionIdSetter can be.
func CreateSession(sessionIdGetter SessionIdGetter, sessionIdSetter SessionIdSetter, options []SessionOptions) (Session, error) {
	if sessionIdGetter == nil {
		return nil, errors.New("invalid session id getter")
	}
	opt := &sessionOptions{RedisClusters: nil, Expiration: 24 * time.Hour, RedisTimeout: 5 * time.Second}
	for _, o := range options {
		o.apply(opt)
	}
	var strategy SessionStore = nil
	if opt.Store != nil {
		strategy = opt.Store
	} else if len(opt.RedisClusters) != 0 {
		_strategy, err := NewRedisSessionStore(opt)
		if err != nil {
			return nil, err
		}
		strategy = _strategy
	} else {
		strategy = NewInMemorySessionStore()
	}

	sid := sessionIdGetter()
	v, err := strategy.Get(context.Background(), sid)
	if err == nil {
		var vv map[string]string
		if err := json.Unmarshal([]byte(v), &vv); err != nil {
			return nil, err
		}
		return &session{storeStrategy: strategy, sid: sid, expiration: opt.Expiration, values: vv}, nil
	}
	sid = newUUID()
	if sessionIdSetter != nil {
		sessionIdSetter(sid)
	}

	return &session{
		storeStrategy: strategy,
		sid:           sid,
		values:        make(map[string]string),
		expiration:    opt.Expiration,
	}, nil
}

// Session represents session which can get and put data into. You can call Get and Set any times but nothing would be store
// until Save is called.
type Session interface {
	// Get by key, returning the value and whether the key exists.
	Get(key string) (interface{}, bool)

	// Set key-value pair
	Set(key, value string) error

	// Save saves all the key-value pairs set before
	Save(ctx context.Context) error

	// get session id
	SessionId() string
}

type session struct {
	mutex         sync.RWMutex
	sid           string
	values        map[string]string
	expiration    time.Duration
	storeStrategy SessionStore
}

func (s *session) Get(key string) (interface{}, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	v, ok := s.values[key]
	return v, ok
}

func (s *session) Set(key, value string) error {
	s.mutex.Lock()
	s.mutex.Unlock()
	s.values[key] = value
	return nil
}

func (s *session) Save(ctx context.Context) error {
	vv, err := json.Marshal(s.values)
	if err != nil {
		return err
	}
	return s.storeStrategy.Set(ctx, s.sid, string(vv), s.expiration)
}

func (s *session) SessionId() string {
	return s.sid
}

type SessionOptions interface {
	apply(*sessionOptions)
}

// WithRedisClusters specifies redis clusters. If redis only contains one node, put it in the slice whose length is 1.
func WithRedisClusters(clusters []string) SessionOptions {
	return newFuncOption(func(option *sessionOptions) {
		option.RedisClusters = clusters
	})
}

// WithExpiration specifies the session expiration time. Session would be removed when the expiration time passes.
func WithExpiration(expiration time.Duration) SessionOptions {
	return newFuncOption(func(option *sessionOptions) {
		option.Expiration = expiration
	})
}

// WithRedisTimeout specifies the redis timeout
func WithRedisTimeout(timeout time.Duration) SessionOptions {
	return newFuncOption(func(option *sessionOptions) {
		option.RedisTimeout = timeout
	})
}

// WithSessionStore specifies custom session store implementation. By default, if redis clusters are specified,
// RedisSessionStore would be used which is implemented based on redis. Otherwise, InMemorySessionStore would be used,
// which is implemented in memory. You can implement your session store using other persistent strategy such as database.
func WithSessionStore(store SessionStore) SessionOptions {
	return newFuncOption(func(option *sessionOptions) {
		option.Store = store
	})
}

type sessionOptions struct {
	RedisClusters []string
	Expiration    time.Duration
	RedisTimeout  time.Duration
	Store         SessionStore
}

type funcOption struct {
	f func(option *sessionOptions)
}

func (f *funcOption) apply(option *sessionOptions) {
	f.f(option)
}

func newFuncOption(f func(option *sessionOptions)) SessionOptions {
	return &funcOption{f: f}
}

func newUUID() string {
	var buf [16]byte
	io.ReadFull(rand.Reader, buf[:])
	buf[6] = (buf[6] & 0x0f) | 0x40
	buf[8] = (buf[8] & 0x3f) | 0x80

	dst := make([]byte, 36)
	hex.Encode(dst, buf[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], buf[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], buf[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], buf[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], buf[10:])

	return string(dst)
}
