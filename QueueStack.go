//
// (C) 2017 Yamato Digital Audio
// Author: Malin af Lääkkö
//
//
// Copyright 2017 Malin Lääkkö -- Yamato Digital Audio.  All rights reserved.
// https://github.com/MalinYamato
//
// Yamato Digital Audio https://yamato.xyz
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Yamato Digital Audio. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

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

