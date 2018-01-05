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

//      expected return value
//
//     [publisherA] --> [targetA],[targetB],[targetC], n
//     [targetA]  --> [target = PublisherA],[target], n
//     [targetB]  --> [target = publiherA],[target],[target], n
//     [targetC]  --> [target], n
//

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Graph struct {
	Basenode          UserId            `json:"basenode"`
	PublishersTargets PublishersTargets `json:"publishersTargets"`
	Pictures          map[UserId]string `json:"pictures"`
}
type UserId string
type Targets map[UserId]bool
type PublishersTargets map[UserId]map[UserId]bool

func (t PublishersTargets) collectAllTargets(pub UserId) (p PublishersTargets, targets Targets) {
	p = make(PublishersTargets)
	targets = make(Targets)
	for k, _ := range t[pub] {
		targets[k] = true
		//	log.Println("target>>>>",k)
		p[k] = make(Targets)
		//log.Println(k)
		for k2, _ := range t[k] {
			//	log.Println(k2)
			p[k][k2] = t[k][k2]
		}
	}
	return p, targets
}

func (t PublishersTargets) Status(pub UserId, target UserId) (publish PublishersTargets, status Status) {

	var pt = make(PublishersTargets)
	var a_to_b bool = false
	var b_to_a bool = false

	if _, ok := t[pub][target]; ok == true {
		if pt[pub] == nil {
			pt[pub] = make(Targets)
		}
		pt[pub][target] = true
		a_to_b = true
	}
	if _, ok := t[target][pub]; ok == true {
		if pt[target] == nil {
			pt[target] = make(Targets)
		}
		pt[target][pub] = true
		b_to_a = true
	}
	if a_to_b && b_to_a {
		status.Status = GREEN
	} else {
		status.Status = BLUE
	}

	return pt, status

}

type PublishRequest struct {
	Op  string   `json:"op"`
	Ids []string `json:"ids"`
}

type PublishRequestResponse struct {
	Op     string `json:"op"`
	Status Status `json:"status"`
	Person Person `json:"person"`
}

func TargetManagerHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var request PublishRequest
	response := PublishRequestResponse{"RequestResponse", Status{}, Person{}}
	if r.Method == "POST" {
		var client Person
		var ok bool
		token, _, err := getCookieAndTokenfromRequest(r, true)
		if err != nil {
			response.Status = Status{ERROR, err.Error()}
		} else {
			client, ok = _persons.findPersonByToken(token)
			if !ok {
				response.Status = Status{ERROR, err.Error()}
			} else {
				decoder := json.NewDecoder(r.Body)
				err = decoder.Decode(&request)
				if err != nil {
					log.Println("Json decoder error> ", err.Error())
					panic(err)
				}
				log.Println(request)
				targetID := UserId(request.Ids[0])
				target, ok := _persons.findPersonByUserId(targetID)
				if !ok {
					log.Printf("Main: Target  not found for UserID %s \n", targetID)
					response.Status = Status{Status: WARNING, Detail: fmt.Sprintf("Receiver not found for UserID %s \n", targetID)}
				} else {
					log.Printf("Main: Profile request for Target %s UserID %s token %s \n", target.Email, target.UserID, target.Token)
					targets, ok := _publishers[client.UserID]
					if request.Op == "RemoveTarget" {
						if ok && len(targets) >= 1 {
							delete(_publishers[client.UserID], UserId(request.Ids[0]))
						}
						if ok && len(targets) < 1 {
							delete(_publishers, client.UserID)
						}
					} else if request.Op == "AddTarget" {
						if !ok {
							targets = make(Targets)
						}
						targets[targetID] = true
						_publishers[client.UserID] = targets
					}

					// all the pictures of all possible targets
					var pictures = make(map[UserId]string)
					for k, _ := range _publishers {
						for k2, _ := range _publishers[k] {
							p, _ := _persons.findPersonByUserId(k2)
							pictures[k2] = p.PictureURL
						}
					}

					targetgraph := Graph{Basenode: target.UserID, PublishersTargets: _publishers, Pictures: pictures}
					hub.multicast <- Message{Op: "UpdateTargetGraph", Token: "", Room: target.Room, Sender: client.UserID, Targets: Targets{target.UserID: true}, Nic: "", Timestamp: timestamp(), Content: "UpdateGraph", Graph: targetgraph}
					clientgraph := Graph{Basenode: client.UserID, PublishersTargets: _publishers, Pictures: pictures}
					hub.multicast <- Message{Op: "UpdateTargetGraph", Token: "", Room: client.Room, Sender: client.UserID, Targets: Targets{client.UserID: true}, Nic: "", Timestamp: timestamp(), Content: "UpdateGraph", Graph: clientgraph}

					response.Status = Status{SUCCESS, "DONT"}
					response.Person = target
				}
			}
		}
		json_response, err := json.Marshal(response)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json_response)

	} else {
		log.Println("Main Unknown HTTP method ", r.Method)
	}
}
