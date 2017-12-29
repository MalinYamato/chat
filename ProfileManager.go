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
	"net/http"
	"log"
	"fmt"
	"encoding/json"
	"html/template"
	//"google.golang.org/api/adexchangeseller/v1"
)

var (
	LANGUAGES   = []string{"English", "Finnish", "Same", "Swedish", "German", "French", "Spannish", "Italian", "Portogese", "Russian", "Chinese", "Japanese", "Korean", "Thai" }
	ORIENTATION = []string{"Straight", "Gay", "Lesbian", "BiSexual", "ASexual"}
	GENDER      = []string{"Female", "Male", "TranssexualF", "TranssexualM", "CrossDresser", "None"}
)

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, sessionName)
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
	t = template.Must(template.ParseFiles( homepath + "registration.html"))
	t.Execute(w, struct {
		P                  Person
		Host               string
	}{
		P:                  p,
		Host:               r.Host,
	})
}



func profileHandler(w http.ResponseWriter, r *http.Request) {
	var request Person;
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&request)
		if err != nil {
			log.Println("ERR> ", err)
		}
		defer r.Body.Close()
		log.Printf("Main: Profile request for user UserID: %s \n", request.UserID)
		var person Person
		person, ok := _persons.findPersonByUserId(request.UserID)
		person.Token = ""
		if ok {
			log.Printf("Main: User not found for UserID %s \n", request.UserID)
		} else {
			log.Printf("Main: Profile request for user %s UserID %s token %s \n", person.Email, person.UserID, person.Token)
		}
		data, err := json.Marshal(person)
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

func mainProfileHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, sessionName)
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
	t = template.Must(template.ParseFiles( homepath + "profile.html"))
	t.Execute(w, struct {
		Languages          []string
		Genders            []string
		SexualOrientations []string
		P                  Person
		Host               string
	}{
		Languages:          LANGUAGES,
		Genders:            GENDER,
		SexualOrientations: ORIENTATION,
		P:                  p,
		Host:               r.Host,
	})
}
func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

func checkNicname(nic string) Status {
	var status Status
	var _, found = _persons.findPersonByNickName(nic)
	if found == true {
		status.Status = "WARNING"
		status.Detail = "Nicna me is taken"
	} else {
		status.Status = "SUCCESS"
		status.Detail = "Nicname is valid"
	}
	return status
}

func updateProfileHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	var status Status
	if r.Method == "POST" {
		session, err := sessionStore.Get(r, sessionName)
		if err != nil {
			log.Println("Main: UpdateProfileHandler() Call to sessionStore.Get returned ", err)
			status.Status = ERROR
			status.Detail = "Failed to get a valid cookie!"
		} else if session == nil {
			log.Println("Main: UpdateProfileHandler() returned session was nil")
			status.Status = ERROR
			status.Detail = "The session is not valid!"
		} else {
			token := session.Values[sessionToken].(string)
			var p Person
			p, _ = _persons.findPersonByToken(token)
			var op = r.FormValue("OP")
			var nic = r.FormValue("NicName")

			if op == "cancel" {
				_persons.Delete(p);
				sessionStore.Destroy(w, sessionName)
				log.Println("User Deleted, session destroyed")
			}

			if op == "checkNicname" {
				status = checkNicname(nic)

			} else if op == "register" {
				status = checkNicname(nic)
				if status.Status == "SUCCESS" {
					p.Nic = nic
					_persons.Save(p)
					status.Detail = "Registration successful"
					hub.broadcast <- Message{Op: "NewUser", Token: "", Room: p.Room, Timestamp: timestamp(), Sender: p.UserID, Nic: p.getNic(), PictureURL: p.PictureURL, Content: "新入社員　" + p.getNic() }
				}

			} else if op == "update" {

				p.FirstName = r.Form.Get("FirstName")
				p.PictureURL = r.Form.Get("PictureURL")
				p.LastName = r.Form.Get("LastName")
				p.Gender = r.Form.Get("Gender")
				p.Country = r.Form.Get("Country")
				p.Town = r.Form.Get("Town")
				p.Lat = r.Form.Get("Lat")
				p.Long = r.Form.Get("Long")
				//p.Nic = r.Form.Get("Nic")
				p.Profession = r.Form.Get("Profession")
				p.Education = r.Form.Get("Education")
				p.SexualOrientation = r.Form.Get("SexualOrientation")
				p.Description = r.Form.Get("Description")
				p.BirthDate.Year = r.Form.Get("BirthYear")
				p.BirthDate.Month = r.Form.Get("BirthMonth")
				p.BirthDate.Day = r.Form.Get("BirthDay")
				p.LoggedIn = true
				fmt.Printf("%+v\n", r.Form)
				productsSelected := r.Form["Language"]
				log.Println(Contains(productsSelected, "English"))
				for i := 0; i < len(LANGUAGES); i++ {
					if Contains(r.Form["Language"], LANGUAGES[i]) {
						p.Languages[LANGUAGES[i]] = "checked"
					}
				}
				_persons.Save(p)
				status.Status = "Updated"
				status.Detail = "Success! Your profile was updated!";
			}
		}
	} else {
		status.Status = "ERROR"
		status.Detail = "Wrong HTTPS Method. Reqire POST but you sent: " + r.Method;
		log.Println( status.Detail + " method " + r.Method)
	}
	data, err := json.Marshal(status)
	if err != nil {
		panic(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
