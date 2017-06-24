// Package semaphore implements a basic semaphore object widely inspired from
// source: http://www.golangpatterns.info/concurrency/semaphores
package semaphore

type empty struct{}
type Semaphore chan empty

// New : new Semaphore object
func New(capacity uint64) Semaphore {
	return make(Semaphore, capacity)
}

// Acquire : increment Semaphore by one
func (s Semaphore) Acquire() {
	s.P(1)
}

// Release : decrement Semaphore by one
func (s Semaphore) Release() {
	s.V(1)
}

// P : acquire n resources
func (s Semaphore) P(n int) {
	e := empty{}
	for i := 0; i < n; i++ {
		s <- e
	}
}

// V : release n resources
func (s Semaphore) V(n int) {
	for i := 0; i < n; i++ {
		<-s
	}
}
