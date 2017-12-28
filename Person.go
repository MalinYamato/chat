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
	"log"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"os"
	"fmt"

)

type Person struct {
	Keep              bool          `json:"keep"`
	Nic               string        `json:"nic"`
	FirstName         string        `json:"firstName"`
	LastName          string        `json:"lastName"`
	Email             string        `json:"email"`
	Gender            string        `json:"gender"`
	Height            string        `json:"height,omitempty"`
	Town              string        `json:"country"`
	Country           string        `json:"town"`
	Long              string        `json:"long,omitempty"`
	Lat               string        `json:"lat,omitempty"`
	PictureURL        string        `json:"pictureURL,omitempty"`
	SexualOrientation string        `json:"sexualOrienation"`
	BirthDate         Date          `json:"birthDate"`
	Languages         map[string]string `json:"Languages,omitempty"`
	Profession        string        `json:"profession"`
	Education         string        `json:"education"`
	Description       string        `json:"description,omitempty"`
	GoogleID          string        `json:"googleId,omitempty"`
	FacebookID        string        `json:"facebookId,omitempty"`
	UserID            UserId        `json:"userId"`
	Token             string        `json:"token,omitempty"`
	Room              string        `json:"room"`
	LoggedIn          bool          `json:"loggedIn"`
	_Persons          *Persons
}

/////////////// Person factory ////////////////////

type Persons struct {
	__pers map[UserId]Person
}

func (pers *Persons) load() {

	if _, err := os.Stat(pers.path()); err != nil {

		if os.IsNotExist(err) {
			log.Println("The directory: " + pers.path() + " does not exist, ignore loading" , err)
			return
		}
	}

	files, err := ioutil.ReadDir( pers.path())
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		content, err := ioutil.ReadFile(pers.path() + "/" + file.Name() + "/profile.json")
		if err != nil {
			log.Fatal(err)
		}
		var person Person
		err = json.Unmarshal(content, &person);
		if err != nil {
			fmt.Println("error:", err)
		}
		pers.__pers[ person.UserID ] = person
	}
}

func (p *Person) getNic() string {
	if p.Nic == "" {
		return p.FirstName + " " + p.LastName
	} else {
		return p.Nic
	}
}


func (pers *Persons) getAll() (persons []Person) {
	var l = []Person{}
	for _, p := range pers.__pers {
		l = append(l, p)
	}
	return l
}
func (pers *Persons) getAllLoggedIn() (persons []Person) {
	var l = []Person{}
	for _, p := range pers.__pers {
		if p.LoggedIn == true {
			l = append(l, p)
		}
	}
	return l
}
func (pers *Persons) getAllInRoom(Room string) (persons []Person) {
	var l = []Person{}
	for _, p := range pers.__pers {
		if p.LoggedIn == true && p.Room == Room {
			l = append(l, p)
		}
	}
	return l
}
func (pers *Persons) findPersonByToken(token string) (person Person, ok bool) {
	for _, p := range pers.__pers {
		if p.Token == token {
			return p, true
		}
	}
	return Person{}, false
}
func (pers *Persons) findPersonByNickName(nic string) (person Person, ok bool) {
	for _, p := range pers.__pers {
		if p.Nic == nic {
			return p, true
		}
	}
	return Person{}, false
}

func (pers *Persons) findPersonByCookie(r *http.Request) (person Person, status Status) {
	var client Person
	var ok bool
	token, _, err := getCookieAndTokenfromRequest(r, true)
	if err != nil {
		return Person{}, Status{ERROR, err.Error()}
	} else {
		client, ok = _persons.findPersonByToken(token)
		if ! ok {
			return Person{}, Status{ERROR, err.Error()}
		}
	}
	return client, Status{SUCCESS, ""}
}
func (pers *Persons) findPersonByGoogleID(GoogleId string) (person Person, ok bool) {
	for _, p := range pers.__pers {
		if p.GoogleID == GoogleId {
			return p, true
		}
	}
	return Person{}, false
}
func (pers *Persons) findPersonByFacebookID(facebookId string) (person Person, ok bool) {
	for _, p := range pers.__pers {
		if p.FacebookID == facebookId {
			return p, true
		}
	}
	return Person{}, false
}


func (pers *Persons) findPersonByUserId(UserId UserId) (person Person, ok bool) {
	person, ok = pers.__pers[UserId]
	return
}

func (pers *Persons) Add(person Person) bool {
	person._Persons = pers
	pers.__pers[ person.UserID ] = person
	return true
}

func (pers *Persons) Save(person Person) bool {
	person.Keep = true
	pers.Add(person)


	if _, err := os.Stat(pers.path()); err != nil {

		if os.IsNotExist(err) {
			log.Println("Creating " + pers.path(), err)
			path := pers.path()
			err := os.Mkdir(path, 0777)
			log.Println("Mkdirerr err ", err)
			if err != nil {
				panic(err)
			}
		}
	}
	path := person.path()
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(path, 0777)
			log.Println("Mkdirerr err ", err)
			if err != nil {
				panic(err)
			}
		}
	}

	path = person.path() + "/img"
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(path, 0777)
			log.Println("Mkdirerr err ", err)
			if err != nil {
				panic(err)
			}
		}
	}

	json_person, _ := json.Marshal(person)
	err := ioutil.WriteFile(person.path()+"/profile.json", json_person, 0777)
	if err != nil {
		panic(err)
	}

	log.Println("Number of persons ", len(pers.__pers))
	return true
}
func (pers *Persons) DeleteById(UserId UserId) bool {
	delete(pers.__pers, UserId)
	return true
}
func (pers *Persons) Delete(user Person) bool {
	delete(pers.__pers, user.UserID)
	return true
}
func (pers *Persons) path() string {
	return "./users"
}
//////////// Person //////////////

func (p *Person) path() string {
	return p._Persons.path() + "/" + string(p.UserID)
}