//
// This work is mostly done by the Gorilla team with some modificaiton by
// Malin Lääkkö

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Copyright 2017 Malin Yamato Lääkkö --  All rights reserved.
// https://github.com/MalinYamato
//
// MIT License
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
//     * Neither the name of Rakuen. nor the names of its
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

import ()
import (
	"log"
)

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan Message

	multicast chan Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	messages map[string]QueueStack

	command chan Command

	rooms map[string]Room

	updateMediaUsers chan MediaUsers
}

type Room struct {
	Name     string
	Messages []Message
}

type Command struct {
	client *Client
	label  string
	value  string
}

func newHub(stack QueueStack) *Hub {
	return &Hub{

		broadcast:  make(chan Message),
		multicast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]Room),
		command:    make(chan Command),
		messages:   RoomManager_getRooms(),
		updateMediaUsers: make(chan MediaUsers),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			if client.Token != "" {
				person, _ := _persons.findPersonByToken(client.Token)
				client.UserId = person.UserID
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				//_persons.RemoveByToken(client.Token)
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.multicast:

			log.Printf("Hub: multicast from %s to %s in room %s ", message.Sender, message.Room, message.Targets)
			for client := range h.clients {
				if _, ok := message.Targets[client.UserId]; ok == true {
					select {
					case client.send <- message:
					default:
						log.Println("Hub: Close Client")
						close(client.send)
						delete(h.clients, client)
					}
				}
			}



		case message := <-h.broadcast:
			room := h.messages[message.Room]
			room.Push(message)
			h.messages[message.Room] = room
			log.Printf("Hub: broadcast to all from %s in room %s", message.Sender, message.Room)

			if room.Len() > 50 {
				room.TailPop()
			}
			for client := range h.clients {
				person, ok := _persons.findPersonByUserId(client.UserId)
				if ok && person.Room == message.Room {
					select {
					case client.send <- message:
					default:
						log.Println("Hub: Close Client")
						close(client.send)
						delete(h.clients, client)
					}
				}
			}

		case mediaUsers := <-h.updateMediaUsers:
			 setMediaUsers(mediaUsers)

		}
	}
}
