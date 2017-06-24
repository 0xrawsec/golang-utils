package datastructs

type HashMap struct {
	keys   map[string]Hashable
	values map[string]interface{}
}

func NewHashMap() (hm HashMap) {
	hm.keys = make(map[string]Hashable)
	hm.values = make(map[string]interface{})
	return hm
}

func (hm *HashMap) Contains(h Hashable) bool {
	if _, ok := (*hm).keys[h.Hash()]; ok {
		return ok
	}
	return false
}

func (hm *HashMap) Get(h Hashable) (interface{}, bool) {
	if _, ok := (*hm).keys[h.Hash()]; ok {
		v, ok := (*hm).values[h.Hash()]
		return v, ok
	}
	return nil, false
}

func (hm *HashMap) Set(key Hashable, value interface{}) {
	(*hm).keys[key.Hash()] = key
	(*hm).values[key.Hash()] = value
}

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

func (hm *HashMap) Len() int {
	return len(hm.keys)
}
