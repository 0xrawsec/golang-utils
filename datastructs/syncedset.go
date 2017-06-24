package datastructs

import (
	"reflect"
	"sync"
)

// SyncedSet datastruct that represent a thread safe set
type SyncedSet struct {
	sync.RWMutex
	set map[interface{}]bool
}

// NewSyncedSet constructs a new SyncedSet
func NewSyncedSet(sets ...*SyncedSet) (ss SyncedSet) {
	ss.set = make(map[interface{}]bool)
	if len(sets) > 0 {
		for _, s := range sets {
			pDatas := s.List()
			ss.Add(*pDatas...)
		}
	}
	return
}

// Equal returns true if both sets are equal
func (s *SyncedSet) Equal(other *SyncedSet) bool {
	s.RLock()
	defer s.RUnlock()
	test := reflect.DeepEqual(s.set, other.set)
	return test
}

// Add adds data to the set
func (s *SyncedSet) Add(datas ...interface{}) {
	s.Lock()
	defer s.Unlock()
	for _, data := range datas {
		s.set[data] = true
	}
}

// Del deletes data from the set
func (s *SyncedSet) Del(datas ...interface{}) {
	s.Lock()
	defer s.Unlock()
	for _, data := range datas {
		delete(s.set, data)
	}
}

// Intersect returns a pointer to a new set containing the intersection of current
// set and other
func (s *SyncedSet) Intersect(other *SyncedSet) *SyncedSet {
	newSet := NewSyncedSet()
	for k := range s.set {
		if other.Contains(k) {
			newSet.Add(k)
		}
	}
	return &newSet
}

// Union returns a pointer to a new set containing the union of current set and other
func (s *SyncedSet) Union(other *SyncedSet) *SyncedSet {
	newSet := NewSyncedSet()
	for elt := range s.set {
		newSet.Add(elt)
	}
	for elt := range other.set {
		newSet.Add(elt)
	}
	return &newSet
}

// Contains returns true if the syncedset contains all the data
func (s *SyncedSet) Contains(datas ...interface{}) bool {
	s.RLock()
	defer s.RUnlock()
	for _, data := range datas {
		if _, ok := s.set[data]; !ok {
			return false
		}
	}
	return true
}

// List returns a pointer to a new list containing the data in the set
func (s *SyncedSet) List() *[]interface{} {
	s.RLock()
	defer s.RUnlock()
	i := 0
	l := make([]interface{}, s.Len())
	for k := range s.set {
		l[i] = k
		i++
	}
	return &l
}

// Len returns the length of the syncedset
func (s *SyncedSet) Len() int {
	return len(s.set)
}
