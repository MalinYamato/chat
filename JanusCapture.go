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
"github.com/jmoiron/jsonq"
"net/http"
"encoding/json"
"bytes"
"fmt"
"strconv"
"io/ioutil"
"strings"
	"log"
)

const http_server = "http://media.raku.cloud:7088"
const encrypted_server = "https://media.raku.cloud:7889"
const server = encrypted_server
var _debug = false


type JanusRequest struct {
	Janus      string `json:"janus"`
	Transation string `json:"transaction"`
	Secret     string `json:"admin_secret"`
}
type JanusSessions struct {
	Janus      string `json:"janus"`
	Transation string `json:"transaction"`
	Sessions   []int  `json:"sessions"`
}
type JanusHandles struct {
	Janus      string `json:"janus"`
	Transation string `json:"transaction"`
	Session    int    `json:"session"`
	Handles    []int  `json:"handles"`
}

type handleID int

type Publishment struct {
	RoomID int
}
type Subscription struct {
	RoomID    int
	ID        int    //subscriber
	Display   string //display of subscriber
	HandleID  handleID
	PrivateID int //owner of feed
}
type MediaUser struct {
	ID            int
	PrivateID     int
	Display       string
	SessionID     int
	Publishments  map[handleID]Publishment
	Subscriptions map[handleID]Subscription
}
type MediaUsers struct {
	__mus map[string]MediaUser
}

func (mus *MediaUsers) findByDisplay(display string) (MediaUser, bool) {
	mu, err := mus.__mus[display]
	return mu, err
}
func (mus *MediaUsers) update(mu MediaUser) {
	mus.__mus[mu.Display] = mu
}
func (mus *MediaUsers) listenersOf(display string) ([]MediaUser) {
	result := []MediaUser{}
	for _, mediaUser := range mus.__mus {
		for _, aSubby := range mediaUser.Subscriptions {
			if aSubby.Display == display {
				result = append(result, mediaUser)
			}
		}
	}
	return result
}

func (mus *MediaUsers) getAll() (map[string]MediaUser) {
	return mus.__mus
}
func (mus *MediaUsers) count() (int) {
	return len(mus.__mus)
}

func recover() {
	log.Println("getDocument failed, Janus server problaby donw")
}


func getDocument(mess string, path string) (r *http.Response, e error) {
    defer recover()

	url := server + "/admin" + "/" + path
	message := JanusRequest{Janus: mess, Transation: "123", Secret: "janusoverlord"}

	if _debug == true { fmt.Println(url) }
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(message)
	res, err := http.Post(url, "application/json; charset=utf-8", b)
	if (err != nil) {
		log.Println("http.post returned error " + err.Error())
		return nil,  err
	}
	return res, nil
}
func JanusCapture() (MediaUsers) {

	defer recover()

	publishers := MediaUsers{map[string]MediaUser{}}
	subscriptions := map[handleID]Subscription{}

	var sessions JanusSessions
	res, e := getDocument("list_sessions", "")
	if e != nil {
		return MediaUsers{}
	}
	err := json.NewDecoder(res.Body).Decode(&sessions)
	if err != nil {
		fmt.Println("err")
	}
	for i := 0; i < len(sessions.Sessions); i++ {
		var handles JanusHandles
		res, e := getDocument("list_handles", strconv.Itoa(sessions.Sessions[i]))
		if e != nil {
			return MediaUsers{}
		}
		err := json.NewDecoder(res.Body).Decode(&handles)
		if err != nil {
			fmt.Println("err")
		}
		for h := 0; h < len(handles.Handles); h++ {
			res, e = getDocument("handle_info", strconv.Itoa(sessions.Sessions[i])+"/"+strconv.Itoa(handles.Handles[h]))
			if e != nil {
				return MediaUsers{}
			}
			body, _ := ioutil.ReadAll(res.Body)
			data := map[string]interface{}{}
			dec := json.NewDecoder(strings.NewReader(string(body)))
			dec.Decode(&data)
			jq := jsonq.NewQuery(data)
			pubsub, _ := jq.String("info", "plugin_specific", "type")
			if (pubsub == "publisher") {
				display, _ := jq.String("info", "plugin_specific", "display")
				var aPublisher MediaUser
				_, err := jq.Int("info", "streams", "0", "id")
				if err != nil {
					fmt.Println("no streams")
				} else {
					// the publisher is broadcasting
					_, ok := publishers.findByDisplay(display)
					if ! ok {
						aPublisher := MediaUser{}
						aPublisher.Display = display
						aPublisher.Publishments = map[handleID]Publishment{}
						publishers.update(aPublisher)

					}
					aPublisher, _ = publishers.findByDisplay(display)
					aPublisher.SessionID, _ = jq.Int("session_id")
					aPublisher.ID, _ = jq.Int("info", "plugin_specific", "id")
					aPublisher.PrivateID, _ = jq.Int("info", "plugin_specific", "private_id")
					id, _ := jq.Int("handle_id")
					handle_id := handleID(id)
					room, _ := jq.Int("info", "plugin_specific", "room")
					aPublisher.Publishments[handle_id] = Publishment{room}
					publishers.update(aPublisher)
				}
			} else if (pubsub == "listener") {
				id, _ := jq.Int("handle_id")
				handle_id := handleID(id)
				_, err := jq.Int("info", "streams", "0", "id")
				if err != nil {
					// the listener is not listening
					fmt.Println("no listening streams")
				} else {
					_, ok := subscriptions[handle_id]
					if ! ok {
						subscriptions[handle_id] = Subscription{}
					}
					subby := subscriptions[handle_id]
					subby.RoomID, _ = jq.Int("info", "plugin_specific", "room")
					subby.PrivateID, _ = jq.Int("info", "plugin_specific", "private_id")
					subby.ID, _ = jq.Int("info", "plugin_specific", "feed_id")
					subby.Display, _ = jq.String("info", "plugin_specific", "feed_display")
					subby.HandleID = handle_id
					subscriptions[handle_id] = subby
				}
			}
		}
	}
	for _, user := range publishers.__mus {
		for _, subby := range subscriptions {
			if user.PrivateID == subby.PrivateID {
				if (user.Subscriptions == nil) {
					user.Subscriptions = map[handleID]Subscription{}
				}
				user.Subscriptions[subby.HandleID] = subby
				publishers.update(user)
			}
		}
	}

	if _debug == true { testJanusCapture()}

	return publishers
}


func testJanusCapture() {

	publishers := JanusCapture()
	fmt.Printf("count %d\n",  publishers.count() )
	for _, user := range publishers.__mus {
		fmt.Print("User: ")
		fmt.Printf("Display %s ID %d PvtID %d  Session %d\n", user.Display, user.ID, user.PrivateID, user.SessionID)
		fmt.Println("publishes: ")
		for h, pub := range user.Publishments {
			fmt.Printf("Using handle %d in Room %d \n", h, pub.RoomID)
		}
		fmt.Println("subscribes to: ")
		for s, sub := range user.Subscriptions {
			fmt.Printf("Using handle %d in  Room %d to %s with ID %d PvtID %d\n", s, sub.RoomID, sub.Display, sub.ID, sub.PrivateID)
		}
		fmt.Println("Listeners: ")
		listeners := publishers.listenersOf(user.Display)
		for l := 0; l < len(listeners); l++ {
			fmt.Println(listeners[l].Display + " listens on " + user.Display)
		}
		fmt.Println()
	}
}

