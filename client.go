// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"
	"time"
	"github.com/gorilla/websocket"
	//"encoding/json"
	"encoding/json"
	"github.com/gorilla/securecookie"

)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	SessionID = "example-egoogle-app"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	id string

	token string

	UserName string

}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.

func validSession(value string) bool {
	if value == "" {
		log.Println("empty cookie")
		return false
	}
	sess := sessionStore.New("example-google-app")
	err := securecookie.DecodeMulti("example-google-app", value, &sess.Values, sessionStore.Codecs...)
	if err != nil {
		log.Println("Invalid")
		return false
	}
	log.Println("Valid")
	return true
}

func (c *Client) readPump(cookieValue string) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		var m Message
		json.Unmarshal(message, &m)
		log.Println(">>>> ", m.Op)

		if ( ! validSession(cookieValue)) {
			json_message, _ := json.Marshal(Message{Op: "Control", Timestamp: "null", Token: "Invalid Session", Sender: "Server", PictureURL: "null", Gender: "null", Content: "Unauthorized" })
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			log.Println("write Error")
			w.Write(json_message)
		} else {

			if value, ok := Persons[m.Token];  ok {
				log.Println("message", value.PictureURL)
				message, _ := json.Marshal(Message{Op: "Message", Timestamp: m.Timestamp, Token: "null", Sender: value.FirstName, PictureURL: value.PictureURL, Gender: value.Gender, Content: m.Content  })
				c.hub.broadcast <- message
				c.token = m.Token
			} else {
				json_message, _ := json.Marshal(Message{Op: "Control", Timestamp: "null", Token: "Wrong Token", Sender: "Server", PictureURL: "null", Gender: "null",Content: "Unauthorized" })
				w, err := c.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}
				log.Println("Write Wrong Token")
				w.Write(json_message)
			}

		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	list := c.hub.messages.GetAllAsList()
	for i := 0; i < len(list); i++ {
		var m Message
		json.Unmarshal(list[i].([]byte), &m)
		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(list[i].([]byte))
	}
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}
	json_message, _ := json.Marshal(Message{Op: "Control", Timestamp: "null", Sender: "System", PictureURL: "null", Gender: "null", Content: "Ignore" })
	w.Write(json_message)

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {

	var cookieValue = ""
	mycookie, err := r.Cookie("example-google-app")
	if (err == nil) {
		cookieValue = mycookie.Value
	} else {
		cookieValue = ""
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	go client.writePump()
	client.readPump(cookieValue)
}
