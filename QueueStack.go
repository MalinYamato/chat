package main


type QueueStack struct {
	top *Element
	size int
}

type Element struct {
	value interface{} // All types satisfy the empty interface, so we can store anything here.
	next *Element
}

// Return the stack's length
func (s *QueueStack) Len() int {
	return s.size
}

// GEt all elements with the oldest element at the top, wihout altering the stack
func (s *QueueStack) GetAllAsList() []interface{}  {
	list := make([]interface{}, s.size)
	var cur *Element
	cur = s.top
	for i := s.size-1; i >= 0; i-- {
		list[i] = cur.value
		cur = cur.next
	}
	return list
}

// Push a new element onto the stack
func (s *QueueStack) Push(value interface{}) {
	s.top = &Element{value, s.top}
	s.size++
}

// Remove the top element from the stack and return it's value
// If the stack is empty, return nil
func (s *QueueStack) HeadPop() (value interface{}) {
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
		return
	}
	return nil
}
func (s *QueueStack) TailPop() (value interface{}) {
	var cur *Element = s.top
	if cur == nil {
		s.size = 0
		return nil
	} else {
		for {
			if cur.next == nil {
				value = cur.value
				cur = nil
				s.size--
				return value
			}
			cur = cur.next
		}
	}
}

