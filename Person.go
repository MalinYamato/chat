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

import ()
import (
	"log"
)

type Person struct {
	Keep              bool          `json:"keep"`
	Nic               string        `json:"nic"`
	FirstName         string        `json:"firstName"`
	LastName          string        `json:"lastName"`
	Email             string        `json:"email"`
	Gender            string        `json:"gender"`
	Town              string        `json:"country"`
	Country           string        `json:"town"`
	PictureURL        string        `json:"pictureURL,omitempty"`
	SexualOrientation string        `json:"sexualOrienation"`
	BirthDate         Date          `json:"birthDate"`
	Languages         map[string]string `json:"Languages,omitempty"`
	Profession        string        `json:"profession"`
	Education         string        `json:"education"`
	Description       string        `json:"description,omitempty"`
	GoogleID          string        `json:"googleId,omitempty"`
	UserID            UserId        `json:"userId,omitempty"`
	Token             string        `json:"token,omitempty"`
	Room              string        `json:"room"`
	LoggedIn          bool          `json:"loggedIn,omitempty"`
}

func (p *Person) getNic() string {
	if p.Nic == "" {
		return p.FirstName + " " + p.LastName
	} else {
		return p.Nic
	}
}
type Persons struct {
	__pers map[UserId]Person
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
func (pers *Persons) findPersonByGoogleID(GoogleId string) (person Person, ok bool) {
	for _, p := range pers.__pers {
		if p.GoogleID == GoogleId {
			return p, true
		}
	}
	return Person{}, false
}
func (pers *Persons) findPersonByUserId(UserId UserId) (person Person, ok bool) {
	person, ok = pers.__pers[UserId]
	return
}
func (pers *Persons) Save(person Person) bool {
	pers.__pers[ person.UserID ] = person
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
