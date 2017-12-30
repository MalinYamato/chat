
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
	"os"
	"strings"
	"io"
	"log"
	"encoding/json"
	"net/http"
	"strconv"
	"image"
	"image/gif"
	"image/png"
	"image/jpeg"
	"github.com/robfig/graphics-go/graphics"
	"io/ioutil"

	"github.com/satori/go.uuid"
	"unicode"
	"golang.org/x/net/html/atom"
)



type VideoResponse struct {
	Status          Status       `json:"status"`
	CamID           string       `json:"camID"`
}

type VideoRequest struct {
	Op              string       `json:"op"`
	CamID           string       `json:"camID"`
}


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
			status = Status{ERROR, err.Error()}
		} else {
			p, ok = _persons.findPersonByToken(token)
			if ! ok {
				status = Status{ERROR, err.Error()}
			} else {
				decoder := json.NewDecoder(r.Body)
				err = decoder.Decode(&request)
				if err != nil {
					log.Println("Json decoder error> ", err.Error())
					panic(err)
				}
				if request.Op == "publish" {
					p.CamID = request.CamID
					p.CamState = "ON"
					hub.broadcast <- Message{Op: "VideoStarted", Token: "", Timestamp: timestamp(), Sender: p.UserID, Nic: p.getNic(), PictureURL: p.PictureURL, Content: "映像放送開始" + p.getNic() }
					status.Status = SUCCESS
				} else if request.Op == "unpublish" {
					p.CamState = "OFF"
					hub.broadcast <- Message{Op: "VideoStopped", Token: "", Timestamp: timestamp(), Sender: p.UserID, Nic: p.getNic(), PictureURL: p.PictureURL, Content: "映像放送停止" + p.getNic() }
					status.Status = SUCCESS
				} else if request.Op == "getCamId" {
					response.CamID = p.CamID
					status.Status = SUCCESS
					}

			}
		}
	} else {
		status = Status{Status: ERROR, Detail:"Bad HTTPS method"}
		log.Println("ImageManager: Unknown HTTP method ", r.Method)
	}
	json_response, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_response)
}


