package utils

import "log"

// Stack LIFO implementation for usage with CHIP-8 Stack
type Stack struct {
	stack []uint16
}

func (s *Stack) Push(elem uint16) {
	s.stack = append(s.stack, elem)
}

func (s *Stack) Pop() uint16 {
	l := len(s.stack)

	if l == 0 {
		log.Println("Tried popping from stack when length is 0")
		return 0
	}

	elem := s.stack[l-1]
	s.stack = s.stack[:l-1]
	return elem
}

func NewStack() Stack {
	return Stack{stack: make([]uint16, 0)}
}
