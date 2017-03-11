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
//     * Neither the name of Google Inc. nor the names of its
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
	"net/http"
	"encoding/json"
	"log"
)

type RoomRequest struct {
	Op   string
	Room string
}

func flushMessagesInRoom(person Person, targets Targets) {
	theRoom := hub.messages[person.Room]
	list := theRoom.GetAllAsList()
	//log.Println("Room List length ",len(list))
	for i := 0; i < len(list); i++ {
		var msg Message
		msg = list[i].(Message)
		msg.Token = "flash"
		msg.Targets = targets
		hub.multicast <- msg
	}
}

func RoomManager_getRooms() map[string]QueueStack {
	queueStack := map[string]QueueStack{
	"Main":QueueStack{},
	"ReimersHotel" : QueueStack{},
	"MalinFriends":QueueStack{},
	"Japanese": QueueStack{},
	"Lesbian": QueueStack{},
	"Gay": QueueStack{},
	"Trans":  QueueStack{}}
        return queueStack
}

func RoomManagerHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var status Status
	var request RoomRequest
	if r.Method == "POST" {
		var person Person
		var ok bool
		token, _, err := getCookieAndTokenfromRequest(r, true)
		if err != nil {
			status = Status{ERROR, err.Error()}
		} else {
			person, ok = _persons.findPersonByToken(token)
			if ! ok {
				status = Status{ERROR, err.Error()}
			} else {
				decoder := json.NewDecoder(r.Body)
				err = decoder.Decode(&request)
				if err != nil {
					log.Println("Json decoder error> ", err.Error())
					panic(err)
				}
				if request.Op == "ChangeRoom" {
					leavingRoom := person.Room
					person.Room = request.Room
					_persons.Save(person)
					targets := make(Targets)
					targets[person.UserID] = true
					flushMessagesInRoom(person, targets)
					RoomUsers := _persons.getAllInRoom(person.Room)
					hub.broadcast <- Message{Op: "ExitUser", Token: "", Room: leavingRoom, Timestamp: timestamp(), Sender: person.UserID, Nic: person.getNic(), PictureURL: person.PictureURL, Content: "Leaving " + person.getNic() }
					hub.broadcast <- Message{Op: "NewUser", Token: "", Room: person.Room, Timestamp: timestamp(), Sender: person.UserID, Nic: person.getNic(), PictureURL: person.PictureURL, Content: "Entering " + person.getNic() }
					hub.multicast <- Message{Op: "RefreshRoomUsers", Token: "", Room: person.Room, Timestamp: timestamp(), Targets: targets, Sender: person.UserID, Nic: person.getNic(), PictureURL: person.PictureURL, Content: "RoomUsers" + person.getNic(), RoomUsers: RoomUsers }
					status = Status{Status: SUCCESS}
				} else if request.Op == "RefressAllMessages" {
					targets := make(Targets)
					targets[person.UserID] = true
					flushMessagesInRoom(person, targets)
					status = Status{Status: SUCCESS}
				} else {
					status = Status{Status: ERROR, Detail:"Unknown operation > " + request.Op}
				}

			}

		}

	} else {
		status = Status{Status: ERROR}
		log.Println("Main Unknown HTTP method ", r.Method)
	}
	json_response, err := json.Marshal(status)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_response)
}