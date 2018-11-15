package datastructs

// BitSet structure definition
type BitSet struct {
	size int
	set  []uint8
}

// NewBitSet creates a new bitset
func NewBitSet(size int) (bs *BitSet) {
	bs = &BitSet{}
	bs.size = size
	if size%8 == 0 {
		bs.set = make([]uint8, size/8)
	} else {
		bs.set = make([]uint8, (size/8)+1)
	}
	return
}

// Set bit at offset o
func (b *BitSet) Set(o int) {
	bucketID := o / 8
	oInBucket := uint8(o % 8)
	if o >= b.size {
		return
	}
	b.set[bucketID] = (b.set[bucketID] | 0x1<<oInBucket)
}

// Get the value of bit at offset o
func (b *BitSet) Get(o int) bool {
	bucketID := o / 8
	oInBucket := uint8(o % 8)
	if o >= b.size {
		return false
	}
	return (b.set[bucketID]&(0x1<<oInBucket))>>oInBucket == 0x1
}

// Len returns the length of the BitSet
func (b *BitSet) Len() int {
	return b.size
}
