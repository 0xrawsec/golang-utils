package datastructs

import "sync"

type HashMap struct {
	keys   map[string]Hashable
	values map[string]interface{}
}

type Item struct {
	Key   Hashable
	Value interface{}
}

func NewHashMap() (hm HashMap) {
	hm.keys = make(map[string]Hashable)
	hm.values = make(map[string]interface{})
	return hm
}

// Contains returns true if the HashMap contains element referenced by key
func (hm *HashMap) Contains(h Hashable) bool {
	if _, ok := (*hm).keys[h.Hash()]; ok {
		return ok
	}
	return false
}

// Get the element referenced by key in the HashMap
func (hm *HashMap) Get(h Hashable) (interface{}, bool) {
	if _, ok := (*hm).keys[h.Hash()]; ok {
		v, ok := (*hm).values[h.Hash()]
		return v, ok
	}
	return nil, false
}

// Set sets key, value in the map
func (hm *HashMap) Set(key Hashable, value interface{}) {
	(*hm).keys[key.Hash()] = key
	(*hm).values[key.Hash()] = value
}

// Del deletes the key and its associated value
func (hm *HashMap) Del(key Hashable) {
	delete((*hm).keys, key.Hash())
	delete((*hm).values, key.Hash())
}

// Keys returns a channel of Keys used by the HashMap
func (hm *HashMap) Keys() (ch chan Hashable) {
	ch = make(chan Hashable)
	go func() {
		defer close(ch)
		for _, v := range hm.keys {
			ch <- v
		}
	}()
	return
}

// Values returns a channel of Values contained in the HashMap
func (hm *HashMap) Values() (ci chan interface{}) {
	ci = make(chan interface{})
	go func() {
		defer close(ci)
		for _, v := range hm.values {
			ci <- v
		}
	}()
	return
}

// Items returns a channel of Item contained in the HashMap
func (hm *HashMap) Items() (ci chan Item) {
	go func() {
		defer close(ci)
		for k := range hm.keys {
			i := Item{(*hm).keys[k], (*hm).values[k]}
			ci <- i
		}
	}()
	return
}

// Len returns the length of the HashMap
func (hm *HashMap) Len() int {
	return len(hm.keys)
}

// SyncedHashMap is a thread safe HashMap
type SyncedHashMap struct {
	sync.RWMutex
	HashMap
}

// NewSyncedHashMap SyncedHashMap constructor
func NewSyncedHashMap() (hm SyncedHashMap) {
	hm.keys = make(map[string]Hashable)
	hm.values = make(map[string]interface{})
	return hm
}

// Contains returns true if the HashMap contains element referenced by key
func (hm *SyncedHashMap) Contains(key Hashable) bool {
	hm.RLock()
	defer hm.RUnlock()
	if _, ok := (*hm).keys[key.Hash()]; ok {
		return ok
	}
	return false
}

// Get the element referenced by key in the HashMap
func (hm *SyncedHashMap) Get(key Hashable) (interface{}, bool) {
	hm.RLock()
	defer hm.RUnlock()
	if _, ok := (*hm).keys[key.Hash()]; ok {
		v, ok := (*hm).values[key.Hash()]
		return v, ok
	}
	return nil, false
}

// Set sets key, value in the map
func (hm *SyncedHashMap) Set(key Hashable, value interface{}) {
	hm.Lock()
	defer hm.Unlock()
	(*hm).keys[key.Hash()] = key
	(*hm).values[key.Hash()] = value
}

// Del deletes the key and its associated value
func (hm *SyncedHashMap) Del(key Hashable) {
	hm.Lock()
	defer hm.Unlock()
	delete((*hm).keys, key.Hash())
	delete((*hm).values, key.Hash())
}

// Keys returns a channel of Keys used by the HashMap
func (hm *SyncedHashMap) Keys() (ch chan Hashable) {
	ch = make(chan Hashable)
	go func() {
		hm.RLock()
		defer hm.RUnlock()
		defer close(ch)
		for _, v := range hm.keys {
			ch <- v
		}
	}()
	return
}

// Values returns a channel of Values contained in the HashMap
func (hm *SyncedHashMap) Values() (ci chan interface{}) {
	ci = make(chan interface{})
	go func() {
		hm.RLock()
		defer hm.RUnlock()
		defer close(ci)
		for _, v := range hm.values {
			ci <- v
		}
	}()
	return
}

// Items returns a channel of Item contained in the HashMap
func (hm *SyncedHashMap) Items() (ci chan Item) {
	go func() {
		hm.RLock()
		defer hm.RUnlock()
		defer close(ci)
		for k := range hm.keys {
			i := Item{(*hm).keys[k], (*hm).values[k]}
			ci <- i
		}
	}()
	return
}

// Len returns the length of the HashMap
func (hm *SyncedHashMap) Len() int {
	hm.RLock()
	defer hm.RUnlock()
	return len(hm.keys)
}
