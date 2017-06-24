package datastructs

import "sync"

type SyncedMap struct {
	sync.RWMutex
	m map[interface{}]interface{}
}

func NewSyncedMap() (s SyncedMap) {
	s.m = make(map[interface{}]interface{})
	return
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

func (s *SyncedMap) Contains(key interface{}) (value interface{}, ok bool) {
	s.RLock()
	defer s.RUnlock()
	value, ok = s.m[key]
	return
}
