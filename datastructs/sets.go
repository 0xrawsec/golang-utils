package datastructs

import (
	"encoding/json"
	"reflect"
	"sync"
)

// Set datastruct that represent a thread safe set
type Set struct {
	set map[interface{}]bool
}

// NewSet constructs a new SyncedSet
func NewSet(sets ...*Set) *Set {
	s := &Set{make(map[interface{}]bool)}
	for _, set := range sets {
		pDatas := set.List()
		s.Add(pDatas...)
	}
	return s
}

// NewInitSet constructs a new SyncedSet initialized with data
func NewInitSet(data ...interface{}) *Set {
	s := NewSet()
	s.Add(data...)
	return s
}

// Equal returns true if both sets are equal
func (s *Set) Equal(other *Set) bool {
	test := reflect.DeepEqual(s.set, other.set)
	return test
}

// Copy returns a copy of the current set
func (s *Set) Copy() *Set {
	return NewSet(s)
}

// Add adds data to the set
func (s *Set) Add(data ...interface{}) {
	for _, data := range data {
		s.set[data] = true
	}
}

// Del deletes data from the set
func (s *Set) Del(data ...interface{}) {
	for _, data := range data {
		delete(s.set, data)
	}
}

// Intersect returns a pointer to a new set containing the intersection of current
// set and other
func (s *Set) Intersect(other *Set) *Set {
	newSet := NewSet()
	for k := range s.set {
		if other.Contains(k) {
			newSet.Add(k)
		}
	}
	return newSet
}

// Union returns a pointer to a new set containing the union of current set and other
func (s *Set) Union(other *Set) *Set {
	newSet := NewSet()
	for elt := range s.set {
		newSet.Add(elt)
	}
	for elt := range other.set {
		newSet.Add(elt)
	}
	return newSet
}

// Contains returns true if the set contains all the data
func (s *Set) Contains(data ...interface{}) bool {
	for _, data := range data {
		if _, ok := s.set[data]; !ok {
			return false
		}
	}
	return true
}

// List returns a pointer to a new list containing the data in the set
func (s *Set) List() []interface{} {
	i := 0
	l := make([]interface{}, s.Len())
	for k := range s.set {
		l[i] = k
		i++
	}
	return l
}

// Items returns a channel with all the elements contained in the set
func (s *Set) Items() (c chan interface{}) {
	c = make(chan interface{})
	go func() {
		defer close(c)
		for k := range s.set {
			c <- k
		}
	}()
	return c

}

// Len returns the length of the syncedset
func (s *Set) Len() int {
	return len(s.set)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (s *Set) UnmarshalJSON(data []byte) (err error) {
	tmp := make([]interface{}, 0)
	s.set = make(map[interface{}]bool)
	if err = json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	for _, data := range tmp {
		s.set[data] = true
	}
	return
}

// MarshalJSON implements json.Marshaler interface
func (s *Set) MarshalJSON() (data []byte, err error) {
	return json.Marshal(s.List())
}

// SyncedSet datastruct that represent a thread safe set
type SyncedSet struct {
	sync.RWMutex
	set *Set
}

// NewSyncedSet constructs a new SyncedSet
func NewSyncedSet(sets ...*SyncedSet) *SyncedSet {
	ss := &SyncedSet{}
	ss.set = NewSet()
	for _, set := range sets {
		ss.Add(set.List()...)
	}
	return ss
}

// NewInitSyncedSet constructs a new SyncedSet initialized with data
func NewInitSyncedSet(data ...interface{}) *SyncedSet {
	ss := &SyncedSet{}
	ss.set = NewSet()
	ss.Add(data...)
	return ss
}

// Equal returns true if both sets are equal
func (s *SyncedSet) Equal(other *SyncedSet) bool {
	s.RLock()
	defer s.RUnlock()
	test := s.set.Equal(other.set)
	return test
}

// Add adds data to the set
func (s *SyncedSet) Add(data ...interface{}) {
	s.Lock()
	defer s.Unlock()
	s.set.Add(data...)
}

// Del deletes data from the set
func (s *SyncedSet) Del(data ...interface{}) {
	s.Lock()
	defer s.Unlock()
	s.set.Del(data...)
}

// Intersect returns a pointer to a new set containing the intersection of current
// set and other
func (s *SyncedSet) Intersect(other *SyncedSet) *SyncedSet {
	s.RLock()
	defer s.RUnlock()
	newSet := NewInitSyncedSet(s.set.Intersect(other.set).List()...)
	return newSet
}

// Union returns a pointer to a new set containing the union of current set and other
func (s *SyncedSet) Union(other *SyncedSet) *SyncedSet {
	s.RLock()
	defer s.RUnlock()
	newSet := NewInitSyncedSet(s.set.Union(other.set).List()...)
	return newSet
}

// Contains returns true if the syncedset contains all the data
func (s *SyncedSet) Contains(data ...interface{}) bool {
	s.RLock()
	defer s.RUnlock()
	return s.set.Contains(data...)
}

// List returns a pointer to a new list containing the data in the set
func (s *SyncedSet) List() []interface{} {
	s.RLock()
	defer s.RUnlock()
	return s.set.List()
}

// Items returns a channel with all the elements contained in the set
func (s *SyncedSet) Items() (c chan interface{}) {
	s.RLock()
	defer s.RUnlock()
	return s.set.Items()

}

// Len returns the length of the syncedset
func (s *SyncedSet) Len() int {
	s.RLock()
	defer s.RUnlock()
	return s.set.Len()
}

// UnmarshalJSON implements json.Unmarshaler interface
func (s *SyncedSet) UnmarshalJSON(data []byte) (err error) {
	s.Lock()
	defer s.Unlock()
	return s.set.UnmarshalJSON(data)
}

// MarshalJSON implements json.Marshaler interface
func (s *SyncedSet) MarshalJSON() (data []byte, err error) {
	s.RLock()
	defer s.RUnlock()
	return json.Marshal(&s.set)
}
