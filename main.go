// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
	"html/template"
	//"github.com/kabukky/httpscerts"

	"github.com/dghubble/gologin"
	googleOAuth2 "golang.org/x/oauth2/google"

	"github.com/dghubble/gologin/google"
	"golang.org/x/oauth2"
	"strings"
	"github.com/kabukky/httpscerts"
)

type Config struct {
	ClientID     string
	ClientSecret string
}

type HTMLReplace struct {
	Host string
	HostImg string
	LoginVisibility string
	LogoutVisibility string
	Person Person
}

type PersonsMAP map[string]Person


var addr = flag.String("addr", ":8080", "http service address")
var homeTemplate = template.Must(template.ParseFiles("home.html"))

func serveHome(w http.ResponseWriter, r *http.Request, stuff HTMLReplace) {

	log.Println(">> ",r.URL)

	stuff.Host = r.Host;
	stuff.HostImg = "https://secure.krypin.org/"

	if ( strings.Contains( r.URL.Path,"/session") ) {
		log.Println("Redirect to / from", r.URL.Path)

		r.URL.Path = "/"
	}
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTemplate.Execute(w, stuff)
}

func serveHomeLogin(w http.ResponseWriter, r *http.Request) {
        var p = (Person{})
	log.Println("here fail")
	serveHome(w,r, HTMLReplace{ "null", "null", "visible", "hidden", p } )
}
func serveHomeLogout(w http.ResponseWriter, r *http.Request) {

	id := r.FormValue("id")
	log.Println("here session")
	//u := strings.Split(user,"@")
	//secret := uuid.NewV4()

	serveHome(w,r, HTMLReplace { "null", "null", "hidden", "visible", Persons[id] } )

}

func imageHandler(w http.ResponseWriter, r *http.Request) {

	http.FileServer(http.Dir("path/to/file"))

}


func NewMux(config *Config, hub *Hub) *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/session/", serveHomeLogout)
	mux.HandleFunc("/", serveHomeLogin)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	mux.Handle("/profile", requireLogin(http.HandlerFunc(profileHandler)))
	mux.HandleFunc("/logout", logoutHandler)
	// 1. Register Login and Callback handlers
	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  "https://secure.krypin.org:8080/google/callback",
		Endpoint:     googleOAuth2.Endpoint,
		Scopes:       []string{"profile", "email"},
	}
	// state param cookies require HTTPS by default; disable for localhost development
	stateConfig := gologin.DebugOnlyCookieConfig
	mux.Handle("/google/login", google.StateHandler(stateConfig, google.LoginHandler(oauth2Config, nil)))
	mux.Handle("/google/callback", google.StateHandler(stateConfig, google.CallbackHandler(oauth2Config, issueSession(), nil)))
	return mux
}

var Persons PersonsMAP

func main() {

	// Check if the cert files are available.
	err := httpscerts.Check("cert.pem", "key.pem")
	//f they are not available, generate new ones.
	if err != nil {
		err = httpscerts.Generate("cert.pem", "key.pem", "localhost")
		if err != nil {
			log.Fatal("Error: Couldn't create https certs.")
		}
	}

	// read credentials from environment variables if available

	//config := &Config{
	//	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	//	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	//}

	Persons = make(PersonsMAP)

	config := &Config{
		ClientID:     "585900153728-tu7nr57i15m1d8sq8ljiv1e00nol2djr.apps.googleusercontent.com",
		ClientSecret: "Qll7wns7E-5uePpE7nqsm56o",
	}

	// allow consumer credential flags to override config fields
	clientID := flag.String("client-id", "", "Google Client ID")
	clientSecret := flag.String("client-secret", "", "Google Client Secret")
	flag.Parse()
	if *clientID != "" {
		config.ClientID = *clientID
	}
	if *clientSecret != "" {
		config.ClientSecret = *clientSecret
	}
	if config.ClientID == "" {
		log.Fatal("Missing Google Client ID")
	}
	if config.ClientSecret == "" {
		log.Fatal("Missing Google Client Secret")
	}

	queue := new(QueueStack)

	flag.Parse()
	hub := newHub(*queue)
	go hub.run()

	err := http.ListenAndServeTLS(*addr,"cert.pem", "key.pem", NewMux(config, hub) )
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
