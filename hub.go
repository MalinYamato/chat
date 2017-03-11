// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package main


import (
)
import "log"

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

}




type Room struct {
	Name string
        Messages []Message
}

type Command struct {
		client *Client
		label string
		value string
}

func newHub(stack QueueStack) *Hub {
	return &Hub {
		broadcast:  make(chan Message),
		multicast : make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms:	    make(map[string]Room),
		command:    make(chan Command),
		messages:   RoomManager_getRooms(),
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
	                        person, ok  := _persons.findPersonByUserId(client.UserId);
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

		}
	}
}
