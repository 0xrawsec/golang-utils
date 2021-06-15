package datastructs

import "sync"

type SyncedMap struct {
	sync.RWMutex
	m map[interface{}]interface{}
}

func NewSyncedMap() *SyncedMap {
	return &SyncedMap{m: make(map[interface{}]interface{})}
}

func (s *SyncedMap) Add(key, value interface{}) {
	s.Lock()
	defer s.Unlock()
	s.m[key] = value
}

func (s *SyncedMap) Del(key interface{}) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, key)
}

func (s *SyncedMap) Contains(key interface{}) (ok bool) {
	s.RLock()
	defer s.RUnlock()
	_, ok = s.m[key]
	return
}

func (s *SyncedMap) Get(key interface{}) (value interface{}, ok bool) {
	s.RLock()
	defer s.RUnlock()
	value, ok = s.m[key]
	return
}

func (s *SyncedMap) Keys() chan interface{} {
	ci := make(chan interface{})
	go func() {
		defer close(ci)
		s.RLock()
		defer s.RUnlock()
		for k := range s.m {
			ci <- k
		}
	}()
	return ci
}

func (s *SyncedMap) Values() chan interface{} {
	ci := make(chan interface{})
	go func() {
		defer close(ci)
		s.RLock()
		defer s.RUnlock()
		for _, v := range s.m {
			ci <- v
		}
	}()
	return ci
}

func (s *SyncedMap) Len() int {
	return len(s.m)
}
