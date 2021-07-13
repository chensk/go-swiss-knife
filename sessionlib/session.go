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

type SessionIdGetter func() string
type SessionIdSetter func(string)

// CreateSession gets the session store associated with sessionId which can be extracted by SessionIdGetter.
// If the session doesn't exist, create a new one and set the sessionId with SessionIdSetter
func CreateSession(sessionIdGetter SessionIdGetter, sessionIdSetter SessionIdSetter, options []SessionOptions) (Session, error) {
	opt := &sessionOptions{redisClusters: nil, expiration: -1, redisTimeout: 5 * time.Second}
	for _, o := range options {
		o.apply(opt)
	}
	var strategy SessionStore = nil
	if len(opt.redisPsm) != 0 {
		_strategy, err := NewRedisSessionStore(opt)
		if err != nil {
			return nil, err
		}
		strategy = _strategy
	} else {
		return nil, errors.New("no store strategy specified")
	}

	sid := sessionIdGetter()
	v, err := strategy.Get(context.Background(), sid)
	if err == nil {
		var vv map[string]interface{}
		if err := json.Unmarshal([]byte(v), &vv); err != nil {
			return nil, err
		}
		return &session{storeStrategy: strategy, sid: sid, expiration: opt.expiration, values: vv}, nil
	}
	sid = newUUID()
	sessionIdSetter(sid)

	return &session{
		storeStrategy: strategy,
		sid:           sid,
		values:        make(map[string]interface{}),
		expiration:    opt.expiration,
	}, nil
}

type Session interface {
	Get(key string) (interface{}, bool)

	Set(key string, value interface{})

	Save(ctx context.Context) error
}

type session struct {
	mutex         sync.RWMutex
	sid           string
	values        map[string]interface{}
	expiration    time.Duration
	storeStrategy SessionStore
}

func (s *session) Get(key string) (interface{}, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	v, ok := s.values[key]
	return v, ok
}

func (s *session) Set(key string, value interface{}) {
	s.mutex.Lock()
	s.mutex.Unlock()
	s.values[key] = value
}

func (s *session) Save(ctx context.Context) error {
	vv, err := json.Marshal(s.values)
	if err != nil {
		return err
	}
	return s.storeStrategy.Set(ctx, s.sid, string(vv), s.expiration)
}

type SessionOptions interface {
	apply(*sessionOptions)
}

func WithRedisClusters(clusters []string) SessionOptions {
	return newFuncOption(func(option *sessionOptions) {
		option.redisClusters = clusters
	})
}

func WithExpiration(expiration time.Duration) SessionOptions {
	return newFuncOption(func(option *sessionOptions) {
		option.expiration = expiration
	})
}

func WithRedisTimeout(timeout time.Duration) SessionOptions {
	return newFuncOption(func(option *sessionOptions) {
		option.redisTimeout = timeout
	})
}

func WithRedisPsm(psm string) SessionOptions {
	return newFuncOption(func(option *sessionOptions) {
		option.redisPsm = psm
	})
}

type sessionOptions struct {
	redisClusters []string
	expiration    time.Duration
	redisTimeout  time.Duration
	redisPsm      string
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
