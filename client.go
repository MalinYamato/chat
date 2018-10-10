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

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
	//"encoding/json"
	"encoding/json"
	"github.com/gorilla/securecookie"
	"strconv"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 5000
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
	hub    *Hub
	conn   *websocket.Conn
	send   chan Message
	UserId UserId
	Token  string
	Cookie string
}

func (c *Client) token() string {

	return c.Token
}

func (c *Client) validSession() bool {
	if c.Cookie == "" {
		log.Println("Client: Empty! Client does not have a Cookie yet.")
		return false
	}
	sess := _sessionStore.New(sessionName)
	err := securecookie.DecodeMulti(sessionName, c.Cookie, &sess.Values, _sessionStore.Codecs...)
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

		if !c.validSession() {
			json_message := Message{Op: "Status", Sender: "Server", Room: "none", Token: "null", Timestamp: "null", Content: "Unauthorized"}
			c.send <- json_message
			if err != nil {
				log.Println("Client Fail to write JSON on websockets because: ", err)
				return
			}
		} else {
			var person Person
			person, ok := _persons.findPersonByToken(message.Token)
			if ok {
				targets, yes := _publishers[c.UserId]
				if yes {
					theMessage := "[" + strconv.Itoa(len(targets)) + "] " + message.Content
					targets[person.UserID] = ok // the message should be sent to the sender herself.
					message := Message{Op: "PrivateMessage", Token: "", Room: person.Room, Sender: person.UserID, Nic: person.getNic(), Targets: targets, Timestamp: message.Timestamp, PictureURL: person.PictureURL, Content: theMessage}
					c.hub.multicast <- message
				} else {
					if person.Room != "Private" { // should not broadcast to this room
						message := Message{Op: "Message", Token: "", Room: person.Room, Sender: person.UserID, Nic: person.getNic(), Timestamp: message.Timestamp, PictureURL: person.PictureURL, Content: message.Content}
						c.hub.broadcast <- message
					}
				}
			} else {
				log.Println("Client: Invalid Token: ", message.Token)
				json_message := Message{Op: "Status", Token: "Invalid Token", Room: person.Room, Sender: "Server", Timestamp: "null", Content: "Unauthorized"}
				c.send <- json_message
				if err != nil {
					log.Println("Client: Fail to write JSON on  websockets! Err: ", err)
					return
				}

			}

		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Println("User: Writing process terminates because websocet connection terminated!  Invalidate token: ", c.token())
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:

			log.Println("Client: Try to send message to browser", message.Sender, message.Content, message.Timestamp, message.Content)
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				log.Println("Client: fail to use websocket connection. It was probaly closed by client.  Reason", ok)
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
				log.Println("Client: Write process : Ping Error ")
				return
			}
		}
	}
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {

	token, cookie, err := getCookieAndTokenfromRequest(r, false)
	if err != nil {
		log.Println(err)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Client: Upgraded to websocket!")

	person, _ := _persons.findPersonByToken(token)

	client := &Client{hub: hub, conn: conn, send: make(chan Message, 256), UserId: person.UserID, Token: token, Cookie: cookie}
	client.hub.register <- client
	go client.writePump()
	client.readPump()
}
