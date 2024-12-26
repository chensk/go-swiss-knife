package lock

import (
	"errors"
	"github.com/chensk/go-swiss-knife/net"
	"time"
)

func init() {
	ip := net.GetCurrentIpv6()
	CurrentIp = ip
}

// util methods

// LockTask try to get lock and launch task if succeed. Block if failed to get the lock.
// If lockTimeout is specified with and positive value, LockTask would wait for at most timeout.
func LockTask(name string, options ...Options) (Lock, error) {
	conf := Configuration{Name: name}
	for _, opt := range options {
		opt.Apply(&conf)
	}
	if conf.Provider == nil {
		return nil, ErrProviderNotFound
	}
	var lock Lock
	conf.LockBy = CurrentIp
	if l, err := conf.Provider.Lock(conf); err == nil {
		lock = l
	} else {
		return nil, err
	}
	return lock, nil
}

// TryLockTask try to get lock and launch task if succeed. It returns error immediately if failed to get the lock.
// lockTimeout is useless for the method
func TryLockTask(name string, options ...Options) (Lock, error) {
	conf := Configuration{Name: name}
	for _, opt := range options {
		opt.Apply(&conf)
	}
	if conf.Provider == nil {
		return nil, ErrProviderNotFound
	}
	var lock Lock
	conf.LockBy = CurrentIp
	if l, err := conf.Provider.TryLock(conf); err == nil {
		lock = l
	} else {
		return nil, err
	}
	return lock, nil
}

// domain interfaces
type Options interface {
	Apply(opt *Configuration)
}

// WithLockAtMost specifies the maximum holding time of the lock. If atMost added by now arrives and the task has not finished,
// the lock would be release by force. If atMost is not positive, the lock will be always held until it's released explicitly.
func WithLockAtMost(atMost time.Duration) Options {
	return newFuncOptions(func(opt *Configuration) {
		opt.LockAtMost = atMost
	})
}

// WithLockTimeout specifies the maximum waiting time when required lock is held by others. If timeout added by now arrives
// and the lock is still unavailable, LockTask would stop waiting and an error would be returned. Timeout is useless for
// TryLockTask method.
func WithLockTimeout(timeout time.Duration) Options {
	return newFuncOptions(func(opt *Configuration) {
		opt.LockTimeout = timeout
	})
}

func WithProvider(provider LockProvider) Options {
	return newFuncOptions(func(opt *Configuration) {
		opt.Provider = provider
	})
}

type Configuration struct {
	Provider    LockProvider
	Name        string
	LockAtMost  time.Duration
	LockTimeout time.Duration
	LockBy      string
}

type LockProvider interface {
	Lock(Configuration) (Lock, error)

	TryLock(Configuration) (Lock, error)
}

type Lock interface {
	Unlock() error
}

type funcOptions struct {
	f func(opt *Configuration)
}

func (f *funcOptions) Apply(opt *Configuration) {
	f.f(opt)
}

func newFuncOptions(f func(opt *Configuration)) *funcOptions {
	return &funcOptions{f: f}
}

var (
	ErrLockFailed        = errors.New("unable to get lock")
	ErrProviderNotFound  = errors.New("provider not specified")
	ErrLockAtMostMissing = errors.New("lock time at most not specified")
	ErrTimeout           = errors.New("lock failed: waiting timeout")

	CurrentIp string
)
