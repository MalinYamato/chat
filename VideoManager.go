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
	"time"
	"sync"

)

type WebRTCSubscribe struct {
	Display string
	ID      int
}

type WebRTCUser struct {
	Display       string
	ID            int
	Handle        int
	Session       int
	Id            int
	Publishing    bool
	Subscriptions map[string]WebRTCSubscribe
}

type WebRTC map[string]WebRTCUser


type SinglePublisherResponse struct {
	Status   Status   `json:"status"`
	CamID    int      `json:"camID"`
	CamState string   `json:"camState"`
	Persons  []Person `json:"persons"`
}

type PublishersResponse struct {
	Status   Status   `json:"status"`
	Persons  []Person `json:"persons"`
}

type StatusResponse struct {
Status   Status   `json:"status"`
}

type VideoRequest struct {
	Op        string `json:"op"`
	CamID     int    `json:"camID"`
	UserID    string `json:"userID"`
	Publisher string `json:"publisher"`
}

var netClient = &http.Client{
	Timeout: time.Second * 10,
}

//////////// RTC ////////////////

type RTCManager struct {
	hub *Hub
	publishers MediaUsers
}

var __mediaUsers MediaUsers
func getMediaUsers() (*MediaUsers) {
	return  &__mediaUsers
}
func unlockMediaUsers()  {
	_mutex.Unlock()
}
func lockMediaUsers() {
	_mutex.Lock()
}

var _mutex sync.Mutex
func setMediaUsers(mediaUsers MediaUsers) {

	// find new publishers
	mus := __mediaUsers.getAll()
	new_mus := mediaUsers.getAll();
	for k, _ := range new_mus {
		_, ok := mus[k];
		if ! ok {
			p, _ := _persons.findPersonByNickName(k)
			hub.broadcast <- Message{Op: "VideoStarted", Token: "", Timestamp: timestamp(), Room: p.Room, Sender: p.UserID, Nic: p.getNic(), PictureURL: p.PictureURL, Content: "映像放送開始 Video ON!"}
		}
	}
	_mutex.Lock()
	__mediaUsers = mediaUsers
	_mutex.Unlock()
}
func (manager *RTCManager) start() {
	_mutex = sync.Mutex{}
	for {
		time.Sleep(5 * time.Second)
		publishers := JanusCapture()
		// log.Printf("%s %d", "CaptureJanus -- available publishers ",publishers.count())
		manager.hub.updateMediaUsers <- publishers
	}
}
var __rtcManager RTCManager
func startRTCManager() {
	__rtcManager := &RTCManager{hub: hub}
	go __rtcManager.start()
}

//////////// End RTC ////////////////

func VideoManager_handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var json_response []byte
	var request VideoRequest
	if r.Method == "POST" {
		var ok bool
		token, _, err := getCookieAndTokenfromRequest(r, true)
		if err != nil {

		} else {
			_, ok = _persons.findPersonByToken(token)
			if ! ok {
				json_response, err = json.Marshal(StatusResponse{Status{ERROR, "fail to find person by token" }})
			} else {
				decoder := json.NewDecoder(r.Body)
				err = decoder.Decode(&request)
				if err != nil {
					log.Println("Json decoder error> ", err.Error())
					panic(err)
				}

				if request.Op == "getCamID" {
					var response SinglePublisherResponse
					publisher, ok := _persons.findPersonByToken( request.UserID)
					if !ok {
						response.Status = Status{ERROR, "Could not find person"}
					} else {
						lockMediaUsers()
						pubs := getMediaUsers().getAll()
						unlockMediaUsers()
						mu, ok := pubs[publisher.Nic]
						if ! ok {
							response.Status = Status{WARNING, "Publisher not found!"}
							log.Println("user " + request.UserID + " not found")
						} else {
							response.CamID = mu.ID
							response.Status = Status{SUCCESS, ""}
							log.Printf("camid %d\n",response.CamID)
						}
					}
					json_response, err = json.Marshal(response)
					if err != nil {
						panic(err)
					}
				}

				if request.Op == "getAllPublishers" {
					log.Println("getAllPublishers")
					response := PublishersResponse{}
					lockMediaUsers()
					pubs := getMediaUsers().getAll()
					unlockMediaUsers()
					response.Persons = nil
					response.Status = Status{ SUCCESS, "No publishers"}
					for _, v := range pubs {
						person, ok := _persons.findPersonByNickName(v.Display);
						if ok {
							response.Persons = append(response.Persons, person)
							response.Status = Status{ SUCCESS, ""}
						}
					}
					json_response, err = json.Marshal(response)
					if err != nil {
						panic(err)
					}
				}

			}
		}
	} else {
		var err error
		json_response, err = json.Marshal(StatusResponse{Status{ERROR, "Wrong HTTP method" }})
		if err != nil {
			panic(err)
		}
		log.Println("ImageManager: Unknown HTTP method ", r.Method)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json_response)
}
