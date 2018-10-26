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
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var (
	LANGUAGES   = []string{"English", "Finnish", "Same", "Swedish", "German", "French", "Spannish", "Italian", "Portogese", "Russian", "Chinese", "Japanese", "Korean", "Thai"}
	ORIENTATION = []string{"Straight", "Gay", "Lesbian", "BiSexual", "ASexual"}
	GENDER      = []string{"Female", "Male", "TransF", "TransM", "Other"}
	RELATION    = []string{"Single", "Married", "Partner", "Divorced"}
)

type PersonResponse struct {
	Status Status `json:"status"`
	Person Person `json:"person"`
}

type PersonRequest struct {
	Op     string `json:"op"`
	Token  string `json:"token"`
	UserID UserId `json:"userID"`
	Nic    string `json:"nic"`
}

func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

func getSessionUser(r *http.Request) (p Person, s Status) {
	var status Status
	var person Person
	token, _, err := getCookieAndTokenfromRequest(r, true)
	status = Status{SUCCESS, ""}
	var ok bool = false
	if err != nil {
		status = Status{ERROR, err.Error()}
	} else {
		person, ok = _persons.findPersonByToken(token)
		if !ok {
			status = Status{ERROR, err.Error()}
		}
	}
	return person, status
}

func checkNewNicname(nic string) Status {
	var status Status
	var _, found = _persons.findPersonByNickName(nic)
	if found == true {
		status.Status = WARNING
		status.Detail = "Nicname is taken"
	} else {
		status.Status = SUCCESS
		status.Detail = "Nicname is unique and valid"
	}
	return status
}

//////////////////////////////////////////////////////////////////////////////////////

func LaunchRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	session, err := _sessionStore.Get(r, sessionName)
	if err != nil {
		log.Println("Main: mainProfileHandler() Call to sessionStore.Get returned ", err)
		return
	}
	if session == nil {
		log.Println("Main: mainProfileHander() returned session was nil")
		return
	}
	token := session.Values[sessionToken].(string)
	p, _ := _persons.findPersonByToken(token)
	t := template.New("fieldname example")
	t = template.Must(template.ParseFiles("registration.html"))
	var prot = "https"
	if r.Proto == "HTTP/1.1" {
		prot = "http"
	}
	t.Execute(w, struct {
		P        Person
		Host     string
		Protocol string
	}{
		P:        p,
		Protocol: prot,
		Host:     r.Host,
	})
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	var request PersonRequest
	var response PersonResponse
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&request)
		if err != nil {
			log.Println("ERR> ", err)
		}
		defer r.Body.Close()
		log.Printf("Main: Profile request for user UserID: %s \n", request.UserID)
		var person Person
		var ok bool
		response.Status = Status{SUCCESS, ""}
		if request.Op == "getUserByToken" {
			person, ok = _persons.findPersonByToken(request.Token)
		} else if request.Op == "getUserByNic" {
			person, ok = _persons.findPersonByNickName(request.Nic)
		} else if request.Op == "getUserByID" {
			person, ok = _persons.findPersonByUserId(request.UserID)
		}
		if !ok {
			response.Status = Status{WARNING, "fail to find person"}
		}
		if request.Op == "getMyself" {
			person, response.Status = getSessionUser(r)
			if response.Status.Status == SUCCESS {
				ok = true
			} else {
				ok = false
			}
		}
		if !ok {
			log.Printf("Main: User not found for UserID %s \n", request.UserID)
		} else {
			log.Printf("Main: Profile request for user %s UserID %s token %s \n", person.Email, person.UserID, person.Token)
			response.Person = person
		}
		data, err := json.Marshal(response)
		if err != nil {
			panic(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)

	} else {
		log.Println("Main Unknown HTTP method ", r.Method)
	}
}

// case a   CLIENT ---> TARGET

func LaunchProfileHandler(w http.ResponseWriter, r *http.Request) {
	session, err := _sessionStore.Get(r, sessionName)
	if err != nil {
		log.Println("Main: mainProfileHandler() Call to sessionStore.Get returned ", err)
		return
	}
	if session == nil {
		log.Println("Main: mainProfileHander() returned session was nil")
		return
	}
	token := session.Values[sessionToken].(string)
	p, _ := _persons.findPersonByToken(token)
	var prot = "https"
	if r.Proto == "HTTP/1.1" {
		prot = "http"
	}
	t := template.New("fieldname example")
	t = template.Must(template.ParseFiles("profile.html"))
	t.Execute(w, struct {
		Languages          []string
		Genders            []string
		SexualOrientations []string
		Relationship       []string
		P                  Person
		Protocol           string
		Host               string
	}{
		Languages:          LANGUAGES,
		Genders:            GENDER,
		SexualOrientations: ORIENTATION,
		Relationship:       RELATION,
		P:                  p,
		Protocol:           prot,
		Host:               r.Host,
	})
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	var status Status
	var p Person
	if r.Method == "POST" {
		var op = r.FormValue("OP")
		var nic = r.FormValue("NicName")
		p, status = getSessionUser(r)
		if op == "checkNewNicname" {
			status = checkNewNicname(nic)
		} else if op == "register" {
			status = checkNewNicname(nic)
			if status.Status == SUCCESS {
				p.Nic = nic
				_persons.Save(p)
				status.Detail = "Registration successful"
				_hub.broadcast <- Message{Op: "NewUser", Token: "", Room: p.Room, Timestamp: timestamp(), Sender: p.UserID, Nic: p.getNic(), PictureURL: p.PictureURL, Content: "新入社員　" + p.getNic()}
			}
		} else if op == "cancel" {
			log.Println("User Deleted, session destroyed")
			_persons.Delete(p)
			_sessionStore.Destroy(w, sessionName)
		}
	} else {
		status.Status = ERROR
		status.Detail = "Wrong HTTPS Method. Reqire POST but you sent: " + r.Method
		log.Println(status.Detail + " method " + r.Method)
	}
	data, err := json.Marshal(status)
	if err != nil {
		panic(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	var status = Status{SUCCESS, ""}
	var p Person
	r.ParseForm()
	if r.Method == "POST" {
		p, status = getSessionUser(r)
		if status.Status == SUCCESS {
			var op = r.FormValue("OP")
			if op == "update" {
				p.FirstName = r.Form.Get("FirstName")
				p.LastName = r.Form.Get("LastName")
				p.FirstNamePublic, _ = strconv.ParseBool(r.Form.Get("FirstNamePublic"))
				p.LastNamePublic, _ = strconv.ParseBool(r.Form.Get("LastNamePublic"))
				p.PictureURL = r.Form.Get("PictureURL")
				p.Gender = r.Form.Get("Gender")
				p.Country = r.Form.Get("Country")
				p.Town = r.Form.Get("Town")
				p.Lat = r.Form.Get("Lat")
				p.Long = r.Form.Get("Long")
				//p.Nic = r.Form.Get("Nic")
				p.Relationship = r.Form.Get("Relationship")
				p.Children, _ = strconv.Atoi(r.Form.Get("Children"))
				p.Profession = r.Form.Get("Profession")
				p.Education = r.Form.Get("Education")
				p.SexualOrientation = r.Form.Get("SexualOrientation")
				p.Description = r.Form.Get("Description")
				p.BirthDate.Year = r.Form.Get("BirthYear")
				p.BirthDate.Month = r.Form.Get("BirthMonth")
				p.BirthDate.Day = r.Form.Get("BirthDay")
				p.LoggedIn = true
				fmt.Printf("%+v\n", r.Form)
				p.Languages = make(map[string]string)
				//
				//Patch. Could not get map to work on the javascript side forcing me to use list for languages
				//instead of map.
				p.LanguagesList = []string{}
				for i := 0; i < len(LANGUAGES); i++ {
					if Contains(r.Form["Language"], LANGUAGES[i]) {
						p.Languages[LANGUAGES[i]] = "checked"
						p.LanguagesList = append(p.LanguagesList, LANGUAGES[i])
					}
				}
				_persons.Save(p)
				status.Status = "Updated"
				status.Detail = "Success! Your profile was updated!"
			}
		}
	} else {
		status.Status = "ERROR"
		status.Detail = "Wrong HTTPS Method. Reqire POST but you sent: " + r.Method
		log.Println(status.Detail + " method " + r.Method)
	}
	data, err := json.Marshal(status)
	if err != nil {
		panic(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
