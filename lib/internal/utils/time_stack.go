package utils

import "time"

// RequestStack is a stack of requests
type TimeStack struct {
	stack []time.Time
}

// create a new stack
func NewTimeStack() *TimeStack {
	return &TimeStack{stack: make([]time.Time, 0)}
}

// push adds a new element to the top of the stack
func (s *TimeStack) Push(t time.Time) {
	s.stack = append(s.stack, t)
}

// pop removes the top element of the stack
func (s *TimeStack) Pop() time.Time {
	if len(s.stack) == 0 {
		return time.Time{}
	}

	t := s.stack[0]
	s.stack = s.stack[1:]
	return t
}

// peek returns the top element of the stack
func (s *TimeStack) Peek() time.Time {
	if len(s.stack) == 0 {
		return time.Time{}
	}

	return s.stack[0]
}

func (s *TimeStack) Size() int {
	return len(s.stack)
}
