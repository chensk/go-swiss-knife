package pool

import (
	"errors"
	"io"
	"sync"
)

type Pool struct {
	m         sync.Mutex
	resources chan io.Closer
	closed    bool
}

func New(fn func() (io.Closer, error), size int, multiplexing int) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("size too small")
	}
	if multiplexing <= 0 {
		multiplexing = 1
	}
	res := make(chan io.Closer, size*multiplexing)
	for i := 0; i < size; i++ {
		c, err := fn()
		if err != nil {
			return nil, err
		}
		for j := 0; j < multiplexing; j++ {
			res <- c
		}
	}
	return &Pool{
		resources: res,
	}, nil
}

func (p *Pool) Acquire() (io.Closer, error) {
	// blocking wait
	r, ok := <-p.resources
	if !ok {
		return nil, errors.New("pool has been closed")
	}
	return r, nil
}

func (p *Pool) Release(r io.Closer) {
	p.m.Lock()
	defer p.m.Unlock()

	if p.closed {
		_ = r.Close()
		return
	}

	select {
	case p.resources <- r:
	default:
		// pool is full , just close
		_ = r.Close()
	}
}

func (p *Pool) Close() error {
	p.m.Lock()
	defer p.m.Unlock()
	if p.closed {
		return nil
	}
	p.closed = true
	close(p.resources)
	for r := range p.resources {
		if err := r.Close(); err != nil {
			return err
		}
	}
	return nil
}
