package crawler

import (
	"errors"
	"sync"
)

// Queue is a collection that implements First In First Out (FIFO) semanthics.
//
// It is built on top of an array type, so on the average case reading and inserting is O(1), except
// when inserting beyond the capacity of the underlying array which will trigger a full O(n) copy of
// the array to allocate more memory.
//
// It uses a sync.RWMutex to make reads and writes concurrent safe.
type Queue[T any] struct {
	q     []T
	mutex sync.RWMutex
}

// NewQueue returns a pointer to a new queue of type T
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

// Push inserts an element at the back of the queue
func (q *Queue[T]) Push(value T) {
	q.mutex.Lock()
	q.q = append(q.q, value)
	q.mutex.Unlock()
}

// Pop removes the first element of the queue and returns it to the caller
func (q *Queue[T]) Pop() (T, error) {
	if q.IsEmpty() {
		err := errors.New("nothing to pop, queue is empty")
		var zero T
		return zero, err
	}

	q.mutex.Lock()
	value := q.q[0]
	q.q = q.q[1:]
	q.mutex.Unlock()
	return value, nil
}

// IsEmpty returns true if there are no elements in the queue
func (q *Queue[T]) IsEmpty() bool {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	return len(q.q) == 0
}

// Size returns the number of elements in the queue
func (q *Queue[T]) Size() int {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	return len(q.q)
}
