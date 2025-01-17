package util

import (
	"sync"
)

// LockingList an ummutable list of objects
type LockingList[T comparable] struct {
	mu      sync.Mutex
	objects []T
}

// NewLockingList returns an initialized RPSLObjectList
func NewLockingList[T comparable](size int) LockingList[T] {
	return LockingList[T]{objects: make([]T, 0, size)}
}

// Add adds an object the list
func (l *LockingList[T]) Add(obj T) {
	l.mu.Lock()
	l.objects = append(l.objects, obj)
	l.mu.Unlock()
}

// GetBatch will return a slice of objects only if 'size' are available. They are removed from the list.
func (l *LockingList[T]) GetBatch(size int) []T {
	res := []T{}
	l.mu.Lock()
	if len(l.objects) >= size {
		res = l.objects[:size]
		l.objects = l.objects[size:]
	}
	l.mu.Unlock()
	return res
}

// GetAll returns all RPSL objects and empties the internal list.
func (l *LockingList[T]) GetAll() []T {
	l.mu.Lock()
	res := l.objects
	l.objects = []T{}
	l.mu.Unlock()
	return res
}
