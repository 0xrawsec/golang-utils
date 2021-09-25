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

func (s *SyncedMap) Keys() (keys []interface{}) {
	s.RLock()
	defer s.RUnlock()
	keys = make([]interface{}, 0, len(s.m))
	for k := range s.m {
		keys = append(keys, k)
	}
	return
}

func (s *SyncedMap) Values() (values []interface{}) {
	s.RLock()
	defer s.RUnlock()
	values = make([]interface{}, 0, len(s.m))
	for _, v := range s.m {
		values = append(values, v)
	}
	return
}

func (s *SyncedMap) Len() int {
	return len(s.m)
}
