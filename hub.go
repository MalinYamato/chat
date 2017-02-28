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

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	messages map[string]QueueStack

	command chan Command

	rooms map[string]Room


}

type Message struct {
	Op        string `json:"op,omitempty"`
	Room      string `json:"Room"`
	Timestamp string `json:"timestamp,omitempty"`
	Token     string `json:"token,omitempty"`
	Sender    string `json:"sender,omitempty"`
	PictureURL string `json:"pictureURL,omitempty"`
	Gender    string  `json:"gender,omitempty"`
	Content   string `json:"content,omitempty"`
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
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms:	    make(map[string]Room),
		command:    make(chan Command),
		messages:   map[string]QueueStack{
				"Main":QueueStack{},
			        "ReimersHotel" : QueueStack{},
			        "EvaMonikaMalin":QueueStack{},
				"Japanese": QueueStack{},
			        "Lesbian": QueueStack{},
			        "Gay": QueueStack{},
			        "Trans":  QueueStack{}},
	}
}

func (h *Hub) run() {


	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.broadcast:


			room := h.messages[message.Room]
			room.Push(message)
		        h.messages[message.Room] = room

			log.Println("Hub: broadcast", message,message.Room, room.Len() )

			if room.Len() > 50 {
				room.TailPop()
			}
			for client := range h.clients {
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
