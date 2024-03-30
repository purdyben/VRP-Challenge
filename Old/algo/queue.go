package algo

import (
	"errors"
	"sync"
)

var EmptyQueue error = errors.New("Queue Is Empty")

type Queue[T any] struct {
	mu    sync.RWMutex
	Items []T
}

func NewQueue[T any]() *Queue[T] {
	return new(Queue[T])
}

func (q *Queue[T]) lock() {
	q.mu.Lock()
}

func (q *Queue[T]) unlock() {
	q.mu.Unlock()
}

func (q *Queue[T]) rLock() {
	q.mu.RLock()
}

func (q *Queue[T]) rUnlock() {
	q.mu.RUnlock()
}

// Size returns the number of Nodes in the queue
func (q *Queue[T]) Size() int {
	q.rLock()
	defer q.rUnlock()
	return len(q.Items)
}

// Enqueue adds an Node to the end of the queue
func (q *Queue[T]) Enqueue(t T) {
	q.lock()
	defer q.unlock()
	q.Items = append(q.Items, t)
}

// Dequeue removes an Node from the start of the queue
func (q *Queue[T]) Dequeue() (*T, error) {
	q.lock()
	defer q.unlock()

	if len(q.Items) == 0 {
		return nil, EmptyQueue
	}

	item := q.Items[0]
	q.Items = q.Items[1:]
	return &item, nil
}

func (q *Queue[T]) Peek() (*T, error) {
	q.rLock()
	defer q.rUnlock()
	if len(q.Items) == 0 {
		return nil, EmptyQueue
	}
	return &q.Items[0]
}

// IsEmpty returns true if the queue is empty
func (q *Queue[T]) IsEmpty() bool {
	return q.Size() == 0
}
