package datastructs

type Container interface {
	Equal(Container) bool
	Contains(interface{}) bool
	List() []interface{}
	Len() int
}

type Hashable interface {
	Hash() string
}
