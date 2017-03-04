package main

import (
	"log"
	"net/http"
	"time"
	"github.com/gorilla/websocket"
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
	hub    *Hub
	conn   *websocket.Conn
	send   chan Message
	UserId string
	Token  string
	Cookie string
}

func (c *Client) token() (string) {

	return c.Token
}

func (c *Client) validSession() bool {
	if c.Cookie == "" {
		log.Println("Client: Empty! Client does not have a Cookie yet.")
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

func (c *Client) flushRoom(room string) {
	theRoom := c.hub.messages[room]
	list := theRoom.GetAllAsList()
	//log.Println("Room List length ",len(list))
	for i := 0; i < len(list); i++ {
		var msg Message
		msg = list[i].(Message)
		msg.Token = "flash"
		c.send <- msg
	}
}

//type Message struct {
//	Op        string   `json:"op"`
//	Token     string   `json:"token"`
//	Room      string   `json:"room"`
//	Sender    string   `json:"sender"`
//	Receivers map[string]bool `json:"receivers,omitempty"`
//	Timestamp string   `json:"timestamp,omitempty"`
//	Content   string   `json:"content,omitempty"`
//}

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
			c.flushRoom("Main")

		} else if ! c.validSession() {
			json_message := Message{Op: "Status", Sender: "Server", Room: "none", Token: "null", Timestamp: "null", Content: "Unauthorized" }
			c.send <- json_message
			if err != nil {
				log.Println("Client Fail to write JSON on websockets because: ", err)
				return
			}
		} else {
			var person Person
			log.Println("messsssssss "+ message.Sender)
			person, ok := _persons.findPersonByToken(message.Token)
			if ok {
				if message.Op == "ChangeRoom" {
					log.Println("Request to change room to ", message.Content)
					person.Room = message.Content;
					_persons.Save(person)
					c.flushRoom(person.Room)
				} else {
					targets, yes := _publishers[c.UserId]
					length := len(targets.Targets)
					theMessage := "[" + strconv.Itoa(length) + "] " + message.Content
					if yes {
						message := Message{Op: "PrivateMessage", Token: "", Room: person.Room, Sender: person.UserID, Nic: person.getNic(), Targets: targets.Targets, Timestamp: message.Timestamp, PictureURL: person.PictureURL, Content: theMessage  }
						c.hub.multicast <- message
					} else {
						if person.Room != "MPR" { // should not broadcast to this room
							message := Message{Op: "Message", Token: "", Room: person.Room, Sender: person.UserID, Nic: person.getNic(), Timestamp: message.Timestamp, PictureURL: person.PictureURL, Content: message.Content  }
							c.hub.broadcast <- message
						}
					}
				}
			} else {
				log.Println("Client: Invalid Token: ", message.Token)
				json_message := Message{Op: "Status", Token: "Invalid Token", Room: person.Room, Sender: "Server", Timestamp: "null", Content: "Unauthorized" }
				c.send <- json_message
				if err != nil {
					log.Println("Client: Fail to write JSON on  websockets! Err: ", err)
					return
				}

			}

		}
	}
}

/*
type Message struct {
	Op        string    `json:"op"`
	Token     string    `json:"token"`
	Room      string    `json:"room"`
	Sender    string    `json:"sender"`
	Nic       string    `json:"nic,omitempty"`
	Receivers Receivers `json:"receivers,omitempty"`
	Timestamp string    `json:"timestamp,omitempty"`
	PictureURL string   `json:"pictureURL,omitemtpy"`
	Content   string    `json:"content"`
}
*/

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
