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
	"strconv"
	"time"
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

type VideoResponse struct {
	Status   Status `json:"status"`
	CamID    int    `json:"camID"`
	CamState string `json:"camState"`
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
func setMediaUsers(mediaUsers MediaUsers) {
	__mediaUsers = mediaUsers
}
func (manager *RTCManager) start() {

	for {
		time.Sleep(5 * time.Second)
		publishers := JanusCapture()
		testJanusCapture()
		log.Printf("%s %d", "CaptureJanus -- available publishers ",publishers.count())
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
	var status Status
	var request VideoRequest
	var response VideoResponse
	if r.Method == "POST" {
		var p Person
		var ok bool
		token, _, err := getCookieAndTokenfromRequest(r, true)
		if err != nil {
			response.Status = Status{ERROR, err.Error()}
		} else {
			p, ok = _persons.findPersonByToken(token)
			if !ok {
				response.Status = Status{ERROR, err.Error()}
			} else {
				decoder := json.NewDecoder(r.Body)
				err = decoder.Decode(&request)
				if err != nil {
					log.Println("Json decoder error> ", err.Error())
					panic(err)
				}

				///// My own webcam //////

				if request.Op == "setMyCamID" {
					log.Println("setMyCamID")
					p.CamID = request.CamID
					_persons.Save(p)
					log.Println("setMyCamID" + strconv.Itoa(request.CamID))

					response.Status = Status{SUCCESS, ""}
				}

				if request.Op == "publish" {
					p.CamState = "ON"
					_persons.Save(p)
					hub.broadcast <- Message{Op: "VideoStarted", Token: "", Timestamp: timestamp(), Room: p.Room, Sender: p.UserID, Nic: p.getNic(), PictureURL: p.PictureURL, Content: "映像放送開始 Vide started!"}
					response.Status = Status{SUCCESS, ""}
				} else if request.Op == "unpublish" {
					p.CamState = "OFF"
					_persons.Save(p)
					hub.broadcast <- Message{Op: "VideoStopped", Token: "", Timestamp: timestamp(), Room: p.Room, Sender: p.UserID, Nic: p.getNic(), PictureURL: p.PictureURL, Content: "映像放送停止 Video stopped!"}
					response.Status = Status{SUCCESS, ""}

				}

				///// Others webcam //////

				if request.Op == "watchRequest" {

					// get the token
				}

				if request.Op == "getCamID" {
					publisher, ok := _persons.findPersonByToken(request.Publisher)
					if !ok {
						response.Status = Status{ERROR, "Could not find person"}
					} else {
						response.CamID = publisher.CamID
						response.Status = Status{SUCCESS, ""}
					}

				}
				if request.Op == "getCamState" {

					publisher, ok := _persons.findPersonByToken(request.Publisher)
					if !ok {
						status = Status{ERROR, "Could not find person"}
					} else {
						response.CamID = publisher.CamID
						response.Status = Status{SUCCESS, ""}
						if publisher.CamState == "ON" {
							response.CamState = "ON"
							response.Status = Status{SUCCESS, ""}
						} else if publisher.CamState == "OFF" {
							response.CamState = "OFF"
							response.Status = Status{SUCCESS, ""}
							status.Status = SUCCESS
						} else {
							response.CamState = "UNKNOWN"
							response.Status = Status{WARNING, "Camstate is unknonw!"}

						}
					}
				}
			}
		}
	} else {
		status = Status{Status: ERROR, Detail: "Bad HTTPS method"}
		log.Println("ImageManager: Unknown HTTP method ", r.Method)
	}
	json_response, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_response)
}
