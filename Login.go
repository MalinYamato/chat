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
	"github.com/dghubble/gologin/facebook"
	"github.com/dghubble/gologin/google"
	"github.com/rs/xid"
	"log"
	"net/http"

	"strconv"
	"strings"
	"time"
)

const (
	sessionName  = "secure.raku.cloud"
	sessionToken = "SessionToken"
)

// Config configures the main ServeMux.

func checkSet(a string, b string) string {
	if a == "" {
		return b
	}
	return a
}

func timestamp() string {
	date := time.Now()
	return strconv.Itoa(date.Day()) + ":" + strconv.Itoa(date.Hour()) + ":" + strconv.Itoa(date.Minute()) + ":" + strconv.Itoa(date.Second())
}

func issueSessionFB() http.Handler {

	var sessionUserKey = "FacebookUser"

	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		facebookUser, err := facebook.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		secret := xid.New() // used as a secret to verify identity of users who sends websocket messages from the brwoser to the server
		userID := xid.New() // used to identify a user to all other users, not a secret.
		// remove possible old cookies
		if isAuthenticated(req) {
			log.Println("Login: There was an old cookie. Removing it")
			sessionStore.Destroy(w, sessionName)
		}
		session := sessionStore.New(sessionName)
		session.Values[sessionUserKey] = facebookUser.ID
		session.Values[sessionToken] = secret.String()
		err = session.Save(w)
		if err != nil {
			log.Println("Login: could not set session ", err)
		}
		var user string
		var v Person
		var ok bool
		s := strings.Split(facebookUser.Name, " ")
		v, ok = _persons.findPersonByFacebookID(facebookUser.ID)
		if ok && v.Keep {
			person := Person{
				Nic:               checkSet(v.Nic, "unregistered"),
				Keep:              v.Keep,
				FirstName:         checkSet(v.FirstName, s[0]),
				LastName:          checkSet(v.LastName, s[1]),
				Email:             checkSet(v.Email, ""),
				Gender:            checkSet(v.Gender, ""),
				BirthDate:         v.BirthDate,
				Country:           checkSet(v.Country, ""),
				Town:              v.Town,
				Long:              checkSet(v.Long, ""),
				Lat:               checkSet(v.Lat, ""),
				PictureURL:        checkSet(v.PictureURL, ""),
				SexualOrientation: v.SexualOrientation,
				Languages:         v.Languages,
				Profession:        v.Profession,
				Education:         v.Education,
				FacebookID:        v.FacebookID,
				UserID:            v.UserID,
				Token:             secret.String(),
				Description:       v.Description,
				Room:              v.Room,
				RakuMail:          v.RakuMail,
				Relationship:      v.Relationship,
				Children:          v.Children,
				FirstNamePublic:   v.FirstNamePublic,
				LastNamePublic:    v.LastNamePublic,
				Password:          v.Password}

			person.LoggedIn = true
			_persons.Save(person)
			user = "registred user"
			hub.broadcast <- Message{Op: "UserLoggedIn", Token: "", Room: person.Room, Timestamp: timestamp(), Sender: person.UserID, Nic: person.getNic(), PictureURL: person.PictureURL, Content: "入室 Enter" + person.getNic()}
		} else if !ok {
			person := Person{
				//Nic:               s[0],
				Keep:              false,
				FirstName:         s[0],
				LastName:          s[1],
				Email:             "",
				Gender:            "",
				BirthDate:         Date{"1900", "1", "1"},
				Town:              "",
				Country:           "",
				PictureURL:        "",
				SexualOrientation: "",
				Languages:         map[string]string{},
				Profession:        "",
				Education:         "",
				FacebookID:        facebookUser.ID,
				UserID:            UserId(userID.String()),
				Token:             secret.String(),
				Description:       "",
				FirstNamePublic:   true,
				LastNamePublic:    true,
				Room:              "Main"}

			if person.PictureURL == "" {
				person.PictureURL = _config.url() + "/images/default.png"
			}
			person.LoggedIn = true
			_persons.Add(person)
			http.Redirect(w, req, "/registration", http.StatusFound)
			return
		}
		person, _ := _persons.findPersonByFacebookID(facebookUser.ID)
		log.Printf("Login: Successful Login of %s Email: %s  FacebookID: %s Token: %s UserID %s ", user, person.Email, person.FacebookID, person.Token, person.UserID)
		http.Redirect(w, req, "/session", http.StatusFound)
	}

	return http.HandlerFunc(fn)
}

func issueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {

		var sessionUserKey = "GoogleUser"

		ctx := req.Context()
		googleUser, err := google.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		secret := xid.New() // used as a secret to verify identity of users who sends websocket messages from the brwoser to the server
		userID := xid.New() // used to identify a user to all other users, not a secret.
		// remove possible old cookies
		if isAuthenticated(req) {
			log.Println("Login: There was an old cookie. Removing it")
			sessionStore.Destroy(w, sessionName)
		}
		session := sessionStore.New(sessionName)
		session.Values[sessionUserKey] = googleUser.Id
		session.Values[sessionToken] = secret.String()
		err = session.Save(w)
		if err != nil {
			log.Println("Login: could not set session ", err)
		}
		var user string
		var v Person
		var ok bool
		v, ok = _persons.findPersonByGoogleID(googleUser.Id)
		if ok && v.Keep {
			person := Person{
				Nic:               checkSet(v.Nic, "uregistered"),
				Keep:              v.Keep,
				FirstName:         checkSet(v.FirstName, googleUser.GivenName),
				LastName:          checkSet(v.LastName, googleUser.FamilyName),
				Email:             checkSet(v.Email, googleUser.Email),
				Gender:            checkSet(v.Gender, googleUser.Gender),
				BirthDate:         v.BirthDate,
				Country:           checkSet(v.Country, googleUser.Locale),
				Town:              v.Town,
				Long:              checkSet(v.Long, ""),
				Lat:               checkSet(v.Lat, ""),
				PictureURL:        checkSet(v.PictureURL, googleUser.Picture),
				SexualOrientation: v.SexualOrientation,
				Languages:         v.Languages,
				Profession:        v.Profession,
				Education:         v.Education,
				GoogleID:          googleUser.Id,
				UserID:            v.UserID,
				Token:             secret.String(),
				Description:       v.Description,
				Room:              v.Room,
				RakuMail:          v.RakuMail,
				Relationship:      v.Relationship,
				Children:          v.Children,
				FirstNamePublic:   v.FirstNamePublic,
				LastNamePublic:    v.LastNamePublic,
				Password:          v.Password}

			if person.PictureURL == "" {
				person.PictureURL = _config.url() + "/images/default.png"
			}

			person.LoggedIn = true
			_persons.Save(person)
			user = "registred user"
			hub.broadcast <- Message{Op: "UserLoggedIn", Token: "", Room: person.Room, Timestamp: timestamp(), Sender: person.UserID, Nic: person.getNic(), PictureURL: person.PictureURL, Content: "入室 Enter" + person.getNic()}
		} else if !ok {
			person := Person{
				//Nic:               googleUser.GivenName,
				Keep:              false,
				FirstName:         googleUser.GivenName,
				LastName:          googleUser.FamilyName,
				Email:             googleUser.Email,
				Gender:            googleUser.Gender,
				BirthDate:         Date{"1900", "1", "1"},
				Town:              "",
				Country:           googleUser.Locale,
				PictureURL:        googleUser.Picture,
				SexualOrientation: "",
				Languages:         map[string]string{},
				Profession:        "",
				Education:         "",
				GoogleID:          googleUser.Id,
				UserID:            UserId(userID.String()),
				Token:             secret.String(),
				FirstNamePublic:   true,
				LastNamePublic:    true,
				Description:       "",
				Room:              "Main"}

			person.LoggedIn = true
			_persons.Add(person)
			http.Redirect(w, req, "/registration", http.StatusFound)
			return
		}
		person, _ := _persons.findPersonByGoogleID(googleUser.Id)
		log.Printf("Login: Successful Login of %s Email: %s  FacebookID: %s Token: %s UserID %s ", user, person.Email, person.FacebookID, person.Token, person.UserID)
		http.Redirect(w, req, "/session", http.StatusFound)

	}
	return http.HandlerFunc(fn)
}

func logoutHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		req.ParseForm()
		session, _ := sessionStore.Get(req, sessionName)
		token := session.Values[sessionToken].(string)
		var person Person
		person, ok := _persons.findPersonByToken(token)
		if ok {
			person.LoggedIn = false
			_persons.Save(person)
			hub.broadcast <- Message{Op: "UserLoggedOut", Token: "", Room: person.Room, Timestamp: timestamp(), Sender: person.UserID, Nic: person.getNic(), PictureURL: person.PictureURL, Content: "出室、またね " + person.getNic()}
			if person.Keep == false {
				log.Printf("Login: Logout user and remove Remove her profile Email %s  UserId %s Token %s", person.Email, person.UserID, person.Token)
				_persons.Delete(person)

			} else {
				log.Printf("Login: Logout user but keep her Profile Email %s  UserId %s Token %s", person.Email, person.UserID, person.Token)
			}
		}
		sessionStore.Destroy(w, sessionName)
	}
	// redirect does not work for AJAX calls. Redirects have to be implemtend by client
	w.Write([]byte(SUCCESS))
	http.Redirect(w, req, "/", http.StatusFound)
}

func requireLoginNonMember(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if !isAuthenticated(req) {
			http.Redirect(w, req, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}
func requireLogin(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, req *http.Request) {
		session, _ := sessionStore.Get(req, sessionName)
		token := session.Values[sessionToken].(string)
		var person Person
		person.Keep = false
		person, ok := _persons.findPersonByToken(token)
		log.Println(" ok %s  %s", person.FirstName, person.Token, token, strconv.FormatBool(ok), strconv.FormatBool(person.Keep))
		if person.Keep == false {
			_persons.Delete(person)
			sessionStore.Destroy(w, sessionName)
			http.Redirect(w, req, "/", http.StatusFound)
			return
		}
		if !isAuthenticated(req) {
			http.Redirect(w, req, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

func isAuthenticated(req *http.Request) bool {
	_, err := sessionStore.Get(req, sessionName)
	if err == nil {
		return true
	}
	log.Println("Login: Authentication failed, reason: ", err)
	return false
}
