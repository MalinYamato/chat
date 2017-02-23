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

	"golang.org/x/oauth2"
	"strings"
	"github.com/kabukky/httpscerts"
	"github.com/dghubble/gologin/google"
	//"path"
	"path"
	"os"
	"github.com/dghubble/sessions"
)

type Config struct {
	ClientID     string
	ClientSecret string
}

type HTMLReplace struct {
	Host      string
	LoggedIn  string
	LoggedOut string
	Person    Person
}

type PersonsMAP map[string]Person

type Endpoint struct {
	protocol string
	host     string
	port     string
}

func (endpoint *Endpoint) url() (string) {
	return endpoint.protocol + "://" + endpoint.host + ":" + endpoint.port
}



func serveHome(w http.ResponseWriter, r *http.Request, stuff HTMLReplace) {

	stuff.Host = r.Host;

	if ( strings.Contains(r.URL.Path, "/session") ) {
		log.Println(" Set path to / ", r.URL.Path)
		r.URL.Path = "/"
	} else if ( strings.Contains(r.URL.Path, "/images") ) {
		log.Println("Serve ", DocumentRoot+r.URL.Path)
		fp := path.Join(DocumentRoot + r.URL.Path)
		http.ServeFile(w, r, fp)
		return
	} else if ( strings.Contains(r.URL.Path, "/css") ) {
		log.Println("Serve ", DocumentRoot+r.URL.Path)
		fp := path.Join(DocumentRoot + r.URL.Path)
		http.ServeFile(w, r, fp)
		return
	} else if ( strings.Contains(r.URL.Path, "/js") ) {
		log.Println("Serve ", DocumentRoot+r.URL.Path)
		fp := path.Join(DocumentRoot + r.URL.Path)
		http.ServeFile(w, r, fp)
		return
	}

	if r.URL.Path != "/" {
		http.Error(w, "Illegal path "+r.URL.Path, 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	log.Println("Serve ", r.URL)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTemplate.Execute(w, stuff)
}

func serveHomeLogin(w http.ResponseWriter, r *http.Request) {
	var p = (Person{})
	serveHome(w, r, HTMLReplace{"null", "none", "flex", p })
}
func sessionHandler(w http.ResponseWriter, r *http.Request) {

	id := r.FormValue("id")
	serveHome(w, r, HTMLReplace{"null", "flex", "none", Persons[id] })
}

func NewMux(config *Config, hub *Hub) *http.ServeMux {

	mux := http.NewServeMux()
	mux.Handle("/session/", requireLogin(http.HandlerFunc(sessionHandler)))
	mux.HandleFunc("/", serveHomeLogin)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	mux.HandleFunc("/logout", logoutHandler)
	// 1. Register Login and Callback handlers
	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  endpoint.url() + "/google/callback",
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
var hub *Hub
var DocumentRoot string
var endpoint Endpoint
var homeTemplate = template.Must(template.ParseFiles("home.html"))
var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

func main() {

	endpoint = Endpoint{"https", "secure.krypin.xyz", "443"}
	dir, _ := os.Getwd()
	DocumentRoot = strings.Replace(dir, " ", "\\ ", -1)
	queue := new(QueueStack)
	var addr = flag.String("addr", ":"+endpoint.port, "http service address")

	// Check if the cert files are available.
	err := httpscerts.Check("cert.pem", "key.pem")
	//f they are not available, generate new ones.
	if err != nil {
		log.Println("Issuing autosigned Certs..")
		err = httpscerts.Generate("cert.pem", "key.pem", endpoint.host)
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

	flag.Parse()
	hub = newHub(*queue)
	go hub.run()

	log.Println("Starting servoce at " + *addr)
	err = http.ListenAndServeTLS(*addr, "cert.pem", "key.pem", NewMux(config, hub))
	//err = http.ListenAndServe(*addr, NewMux(config, hub) )
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
