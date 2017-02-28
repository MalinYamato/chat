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
	send chan Message

	Token  string
	Cookie string
}

func (c *Client) token() (string) {

	return c.Token
}

func (c *Client) validSession() bool {
	if c.Cookie == "" {
		log.Println("No Cookie was set")
		return false
	}
	sess := sessionStore.New(sessionName)
	err := securecookie.DecodeMulti(sessionName, c.Cookie, &sess.Values, sessionStore.Codecs...)
	if err != nil {
		log.Println("Cookie was invalid, perhaps expired")
		return false
	}
	log.Println("Client: Cookie is valid")
	return true
}

func (c *Client) readPump() {
	defer func() {

		//person, ok := Persons[c.token()]
		if true {
			log.Println("User: Reading processes terminates because websocet connection terminated! Invalidate token:", c.token())
			//hub.broadcast <- Message{"ExitUser", "",person.UserID, person.FirstName, person.PictureURL, person.Gender, "出室、 またね　" + person.FirstName + " " + person.LastName}
			//delete(Persons, c.token)
		} else {
			log.Println("Client: Reading processs terminates because connection terminated. Token not seet or invalid! Token:", c.token())
		}
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, json_message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("Read Wesocked was closed: %v", err)
			}
			log.Printf("Client: Wesocked was closed: %v ", err)
			break
		}

		var message Message
		json.Unmarshal(json_message, &message)

		log.Println("Client: message from Browser ", message)

		if message.Op == "SendAllMessages" {
			list := c.hub.messages.GetAllAsList()
			for i := 0; i < len(list); i++ {
				c.conn.WriteJSON(list[i])
			}
		} else if ! c.validSession() {
			json_message := Message{Op: "Control", Timestamp: "null", Token: "Invalid Session", Sender: "Server", PictureURL: "null", Gender: "null", Content: "Unauthorized" }
			err := c.conn.WriteJSON(json_message)
			if err != nil {
				log.Println("Client Fail to write JSON on websockets because: ", err)
				return
			}
		} else if value, ok := Persons[message.Token]; ok {

			sender := value.FirstName
			if len(value.Nic) > 5 {
				sender = value.Nic
			}

			log.Println("Client: Token was set", c.token())
			message := Message{Op: "Message", Timestamp: message.Timestamp, Token: value.UserID, Sender: sender, PictureURL: value.PictureURL, Gender: value.Gender, Content: message.Content  }
			c.hub.broadcast <- message
		} else {
			log.Println("Client: Invalid Token: ", message.Token)
			json_message := Message{Op: "Control", Timestamp: "null", Token: "Invalid Token", Sender: "Server", PictureURL: "null", Gender: "null", Content: "Unauthorized" }
			err := c.conn.WriteJSON(json_message)
			if err != nil {
				log.Println("Client: Fail to write JSON on  websockets! Err: ", err)
				return
			}

		}

	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		//person, ok := Persons[c.token()]
		if true {
			log.Println("User: Writing process terminates because websocet connection terminated!  Invalidate token: ", c.token())
			// hub.broadcast <- Message{"ExitUser", "",person.UserID, person.FirstName, person.PictureURL, person.Gender, "出室 またね　" + person.FirstName + " " + person.LastName + " "}
			//delete(Persons, c.token)
		} else {
			log.Println("User: Wrting process terminates because websocet connection terminated. Token not seet or invalid! Token:", c.token())
		}
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			log.Println("Client: message from Hub ", message.Sender, message.Content, message.Timestamp, message.Content)
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				log.Println("Client: The hub closed the connection because:", ok)
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.conn.WriteJSON(message)
			if err != nil {
				log.Println("Client: Fail to write  websockets because: ", err)
				return
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				log.Println("Send one more message")
				err := c.conn.WriteJSON(<-c.send)
				//w.Write(newline)
				if err != nil {
					log.Println("Client: Fail to flush remaining messages and write JSON on  websockets! ")
					return
				}
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("Ping Error ")
				return
			}
		}
	}
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {

	var token string = ""
	var cookie string = ""
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		log.Println("Client: Call to sessionStore.Get returned ", err)
	} else if session != nil {
		atoken, ok := session.Values[sessionToken]
		if ok  {
			log.Println(atoken)
			if atoken != nil {
				token = atoken.(string)
			} else {
				log.Println("Client: token was nil")
			}
			mycookie, err := r.Cookie(sessionName)
			if (err == nil) {
				cookie = mycookie.Value
			} else {
				cookie = ""
			}
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Client: Upgraded to websocket!")

	client := &Client{hub: hub, conn: conn, send: make(chan Message, 256), Token: token, Cookie: cookie}
	client.hub.register <- client
	go client.writePump()
	client.readPump()
}
