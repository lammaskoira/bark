package util

import "sync"

type Stack[T any] struct {
	queue []T
	lock  sync.Mutex
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		queue: []T{},
		lock:  sync.Mutex{},
	}
}

func (s *Stack[T]) Push(t T) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.queue = append(s.queue, t)
}

func (s *Stack[T]) Pop() T {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Can't use IsEmpty here as it
	// would deadlock
	if len(s.queue) == 0 {
		var t T
		// should be nil
		return t
	}

	t := s.queue[len(s.queue)-1]
	s.queue = s.queue[:len(s.queue)-1]
	return t
}

func (s *Stack[T]) IsEmpty() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return len(s.queue) == 0
}
