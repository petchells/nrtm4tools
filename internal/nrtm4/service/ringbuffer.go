package service

import (
	"sync"
)

type RingBuffer[T any] struct {
	buffer []T
	size   int
	mu     sync.Mutex
	write  int
	count  int
}

// NewRingBuffer creates a new ring buffer with a fixed size.
func NewRingBuffer[T any](size int) *RingBuffer[T] {
	return &RingBuffer[T]{
		buffer: make([]T, size),
		size:   size,
	}
}

// Add inserts a new element into the buffer, overwriting the oldest if full.
func (rb *RingBuffer[T]) Add(value T) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.buffer[rb.write] = value
	rb.write = (rb.write + 1) % rb.size

	if rb.count < rb.size {
		rb.count++
	}
}

// GetAll returns the contents of the buffer in FIFO order.
func (rb *RingBuffer[T]) GetAll() []T {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	result := make([]T, rb.count)

	for i := range rb.count {
		index := (rb.write + rb.size - rb.count + i) % rb.size
		result[i] = rb.buffer[index]
	}
	return result
}

// Len returns the current number of elements in the buffer.
func (rb *RingBuffer[T]) Len() int {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return rb.count
}
