// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package main


import (
	"encoding/json"
)

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	messages QueueStack

	command chan Command

	rooms map[string]Room


}

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
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
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms:	    make(map[string]Room),
		command:    make(chan Command),
		messages:   stack,
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
			h.messages.Push(message)
			if h.messages.Len() > 25 {
				h.messages.TailPop()
			}
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		case cmd := <-h.command:
			if (cmd.label == "chroom") {
				cmd.client.CurrentRoom = cmd.value
				messages := h.rooms[cmd.client.CurrentRoom].Messages
				for i := 0; i < len(messages); i++ {
					jsonMessage, _ := json.Marshal(&messages[i])
					cmd.client.send <- jsonMessage

				}

			}
		}
	}
}
