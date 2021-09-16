package container

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

type ExpiringSet struct {
	elements       Item
	existence      map[interface{}]bool
	mutex          sync.RWMutex
	expirationTime time.Duration
}

func NewExpiringSet(expirationTime time.Duration) *ExpiringSet {
	s := &ExpiringSet{
		elements:       make([]ItemValue, 0),
		existence:      make(map[interface{}]bool),
		expirationTime: expirationTime,
	}
	go func(s *ExpiringSet) {
		recycle(s)
	}(s)
	return s
}

func (s *ExpiringSet) Exists(value interface{}) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	e, ok := s.existence[value]
	return ok && e
}

func (s *ExpiringSet) Add(value interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.existence[value] = true
	heap.Push(&s.elements, ItemValue{value: value, deadline: time.Now().Add(s.expirationTime)})
}

func recycle(s *ExpiringSet) {
	for {
		top := func() *ItemValue {
			s.mutex.Lock()
			defer s.mutex.Unlock()
			if s.elements.Len() == 0 {
				return nil
			}
			now := time.Now()
			for s.elements.Peek() != nil && s.elements.Peek().deadline.Before(now) {
				p := heap.Pop(&s.elements).(ItemValue)
				delete(s.existence, p.value)
				fmt.Printf("poping %s, deadline: %s, time: %s\n", p.value, p.deadline, now)
			}
			return s.elements.Peek()
		}()
		if top != nil {
			time.Sleep(top.deadline.Sub(time.Now()))
		} else {
			time.Sleep(defaultPeriod)
		}
	}
}

type ItemValue struct {
	value    interface{}
	deadline time.Time
}

type Item []ItemValue

func (item Item) Len() int {
	return len(item)
}

func (item Item) Less(i, j int) bool {
	return item[i].deadline.Before(item[j].deadline)
}

func (item Item) Swap(i, j int) {
	item[i], item[j] = item[j], item[i]
}

func (item *Item) Push(x interface{}) {
	*item = append(*item, x.(ItemValue))
}

func (item *Item) Pop() interface{} {
	p := (*item)[len(*item)-1]
	*item = (*item)[:len(*item)-1]
	return p
}

func (item Item) Peek() *ItemValue {
	if item.Len() == 0 {
		return nil
	}
	return &item[0]
}

var (
	defaultPeriod = 1 * time.Second
)
