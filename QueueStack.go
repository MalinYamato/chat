//
// (C) 2017 Yamato Digital Audio
// Author: Malin af Lääkkö
//

package main


type QueueStack struct {
	top *Element
	size int
}

type Element struct {
	value interface{}
	next *Element
}


func (s *QueueStack) Len() int {
	return s.size
}

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


func (s *QueueStack) Push(value interface{}) {
	s.top = &Element{value, s.top}
	s.size++
}

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

