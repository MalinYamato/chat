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
	"encoding/json"
	"log"
	"net/http"
)

type RoomRequest struct {
	Op   string
	Room string
}

func flushMessagesInRoom(person Person, targets Targets) {
	theRoom := _hub.messages[person.Room]
	list := theRoom.GetAllAsList()
	//log.Println("Room List length ",len(list))
	for i := 0; i < len(list); i++ {
		var msg Message
		msg = list[i].(Message)
		msg.Token = "flash"
		msg.Targets = targets
		_hub.multicast <- msg
	}
}

func RoomManager_getRooms() map[string]QueueStack {
	queueStack := map[string]QueueStack{
		"Main":     QueueStack{},
		"Japanese": QueueStack{},
		"Lesbian":  QueueStack{},
		"Gay":      QueueStack{},
		"Trans":    QueueStack{}}
	return queueStack
}

type RoomElem struct {
	Name     string `json:"name"`
	Japanese string `json:"japanese"`
	Owner    UserId `json:"owner"`
}

type Response struct {
	Status Status     `json:"status"`
	Rooms  []RoomElem `json:"rooms"`
}

var rooms = []RoomElem{
	{"Main", "本館", ""},
	{"Private", "秘密屋", ""},
	{"Japanese", "日本語屋", ""},
	{"Gay", "ゲイ屋", ""},
	{"Lesbian", "レス屋", ""},
	{"Trans", "性転換屋", ""},
}

func roomUsersToTargets(rooms string) {

}

func RoomManagerHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var response Response
	response.Rooms = rooms
	response.Status = Status{"SUCCESS", ""}
	var request RoomRequest
	if r.Method == "POST" {
		var person Person
		var ok bool
		token, _, err := getCookieAndTokenfromRequest(r, true)
		if err != nil {
			response.Status = Status{ERROR, err.Error()}
		} else {
			person, ok = _persons.findPersonByToken(token)
			if !ok {
				response.Status = Status{ERROR, err.Error()}
			} else {
				decoder := json.NewDecoder(r.Body)
				err = decoder.Decode(&request)
				if err != nil {
					log.Println("Json decoder error> ", err.Error())
					panic(err)
				}
				if request.Op == "GetAllRooms" {
					response.Status = Status{Status: SUCCESS} // not implemented
				} else if request.Op == "ChangeRoom" {
					leavingRoom := person.Room
					enterRoom := request.Room

					person.Room = enterRoom
					_persons.Save(person)

					targets := make(Targets)
					targets[person.UserID] = true
					flushMessagesInRoom(person, targets)

					EnterRoomUsers := _persons.getAllInRoom(enterRoom)
					LeavingRoomUsers := _persons.getAllInRoom(leavingRoom)
					EnterRoomTargets := make(Targets)
					LeavingRoomTargets := make(Targets)
					for i := 0; i < len(EnterRoomUsers); i++ {
						EnterRoomTargets[EnterRoomUsers[i].UserID] = true
					}
					for i := 0; i < len(LeavingRoomUsers); i++ {
						LeavingRoomTargets[LeavingRoomUsers[i].UserID] = true
					}

					_hub.multicast <- Message{Op: "UserEnteredRoom", Token: "UserEnteredRoom", Room: person.Room, Timestamp: timestamp(),
						Targets: EnterRoomTargets, Sender: person.UserID, Nic: person.getNic(), PictureURL: person.PictureURL,
						Content: "RoomUsers" + person.getNic(), RoomUsers: EnterRoomUsers}
					response.Status = Status{Status: SUCCESS}

					_hub.multicast <- Message{Op: "UserLeftRoom", Token: "UserLeftRoom", Room: person.Room, Timestamp: timestamp(),
						Targets: LeavingRoomTargets, Sender: person.UserID, Nic: person.getNic(), PictureURL: person.PictureURL,
						Content: "RoomUsers" + person.getNic(), RoomUsers: LeavingRoomUsers}
					response.Status = Status{Status: SUCCESS}

				} else if request.Op == "RefressAllMessages" {
					targets := make(Targets)
					targets[person.UserID] = true
					flushMessagesInRoom(person, targets)
					response.Status = Status{Status: SUCCESS}
				} else {
					response.Status = Status{Status: ERROR, Detail: "Unknown operation > " + request.Op}
				}
			}
		}
	} else {
		response.Status = Status{Status: ERROR}
		log.Println("Main Unknown HTTP method ", r.Method)
	}
	json_response, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_response)
}
