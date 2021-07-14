package sessionlib

import (
	"container/list"
	"context"
	"errors"
	"sync"
	"time"
)

var (
	_store *InMemorySessionStore
	once   sync.Once
)

// link list
type Element struct {
	key      string
	value    string
	deadline time.Time
}

// InMemorySessionStore stores the session in local memory.
type InMemorySessionStore struct {
	mutex sync.RWMutex
	data  map[string]Element
	// link list is used to keep the elements order by its deadline so that gc can focus only on the first element of list
	list *list.List
}

func NewInMemorySessionStore() SessionStore {
	once.Do(func() {
		_store = &InMemorySessionStore{
			data: make(map[string]Element),
			list: list.New(),
		}
		go _store.gc()
	})
	return _store
}

// gc deletes the expired elements
func (s *InMemorySessionStore) gc() {
	defaultPeriod := 1 * time.Second
	for {
		func() {
			if s.list.Front() == nil {
				time.Sleep(defaultPeriod)
				return
			}
			head := s.list.Front().Value.(Element)
			if head.deadline.After(time.Now()) {
				time.Sleep(head.deadline.Sub(time.Now()))
				return
			}
			s.mutex.Lock()
			defer s.mutex.Unlock()
			// delete head
			delete(s.data, head.key)
			s.list.Remove(s.list.Front())
		}()
	}
}

func (s *InMemorySessionStore) Get(ctx context.Context, key string) (string, error) {
	s.mutex.RLock()
	v, ok := s.data[key]
	s.mutex.RUnlock()
	if !ok {
		return "", errors.New("not found")
	}
	return v.value, nil
}

func (s *InMemorySessionStore) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	if _, ok := s.data[key]; ok {
		return nil
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	ele := Element{
		value:    value,
		key:      key,
		deadline: time.Now().Add(expiration),
	}
	s.data[key] = ele
	// insert into list ordered by deadline
	newEle := s.list.PushFront(ele)
	var start *list.Element
	for start = newEle; start.Next() != nil; start = start.Next() {
		if start.Next().Value.(Element).deadline.Before(newEle.Value.(Element).deadline) {
			continue
		}
		break
	}
	s.list.MoveAfter(newEle, start)
	return nil
}

func (s *InMemorySessionStore) Delete(ctx context.Context, key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.data, key)
	// delete from list
	var start *list.Element
	for start = s.list.Front(); start != nil; start = start.Next() {
		if start.Value.(Element).key == key {
			break
		}
	}
	if start != nil && start.Value.(Element).key == key {
		s.list.Remove(start)
	}
	return nil
}
