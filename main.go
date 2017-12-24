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
//     * Neither the name of Yamato Digital Audio Inc. nor the names of its
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
	"fmt"
	"github.com/dghubble/gologin/facebook"
	facebookOAuth2 "golang.org/x/oauth2/facebook"
)

type Config struct {
	ClientID_FB     string
	ClientSecret_FB string
	ClientID        string
	ClientSecret    string
	ChatHost        string
	ChatPrivateKey  string
}

type HTMLReplace struct {
	Host      string
	LoggedIn  string
	LoggedOut string
	Person    Person
}

type Endpoint struct {
	protocol string
	host     string
	port     string
}

type Date struct {
	Year  string  `json:"year"`
	Month string  `json:"month"`
	Day   string  `json:"day"`
}

const
(
	ERROR   = "ERROR"
	WARNING = "WARNING"
	SUCCESS = "SUCCESS"
)
const (
	GREEN = "GREEN" // sender and target are sending pvt messages to each other
	BLUE  = "BLUE"  // sender sends pvt messages to the target but not the other way around
	BLACK = "BLACK" // The target is blocking, black listening the sender
)

type Status struct {
	Status string `json:"status"`
	Detail string `json:"detail"`
}

type VideoFormat struct {
	Codec   string   		`json:"codec"`
	Width   int16    		`json:"width"`// in pixels
	Height  int16   		`json:"height"`// in pixels
	BitRate int16   		`json:"bitRate"`// bits per second
}

type AudioFormat struct {
	Codec      string   		`json:"codec"`
	Channels   int16   		`json:"channels"`
	BitRate    int16        	`json:"bitRate"`	// bits per second
	BitDepth   int16    		`json:"bitDepth"`	// vertical resolution,  PCM
	SampleRate int32    		`json:"sampleRate"`	// Number of vertical snapshots per second, PCM
}

// publishers[].Targets[]

// Media Session Protocol
//
type MediaSession struct {
	MediaServerURL string  			`json:"idMediaServerURL"`
	IdMediaSession string  			`json:"idHandle"`
	IdHandle string        			`json:"id"`
	Id       string    	    		`json:"id"`
	IdRoom   string         	        `json:"room"`
	Audio    bool          			`json:"audio"`
	Video    bool          			`json:"video"`
	PubOrSub string        			`json:"pubOrSub"`
	OnOrOff  string        			`json:"onOrOff"`
	VideoFormat VideoFormat        	 	`json:"VideoFormat,omitempty"`
	AudioFormat AudioFormat        		`json:"AudioFormast,omitempty"`
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Simplifed protocol and structure when handshake and interpreation of SDP
// (Session Description Protocol) is supported by another layer such as JANUS. The purpose of the protocol
// is to let users know without maintaning state on the server-side who are currently publishing video, audio
// or both as well as control of whom is allwoed to se certain other users broadcasts. To dissalov certain users
// from subscribing a certain stream, the normal procedure is to look up the UserId of the AnyPublishers package
// received and decide to reply with a MediaStatus response or not based on that.
// Most of media information is included in SDP and are therefore omitted except for video hight and width.

type MediaStatus struct {
        MedaiServerURL string                          `json:"mediaServerURL"`  // The url of SFU and MediaGateway
	OnOff          string                          `json:"onOff"`
	JanusId        string                          `json:"janusId"`         // the Id used by JAnus to identify streams
	PubOrSub       string                          `json:"pubOrSub"`        // Janus room
	Room           string                          `json:"room"`
	Audio          bool                            `json:"audio"`
	Video          bool                            `json:"video"`
	VideoHeight    int16                           `json:"videoHeight"`     // Pixels. hint how to arrange the GUI to present video
	VideoWidth     int16                           `json:"videoWidth"`      // Pixels. hint how to arrange the GUI to present video
}

//   AnyPuiblishers, broadcasted when interresed in knowing who is/are publishing.
//   MediaStatus     sent as a response upon reception of AnyPublishers if
//                          publishing, not blocking prospective subscribers or when start or stop publishing.

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Message struct {
	Op         string                           `json:"op"`
	Token      string                           `json:"token"`
	Room       string                           `json:"room"`
	Sender     UserId                           `json:"sender"`
	Targets    Targets                          `json:"targets,omitempty"`
	Nic        string                           `json:"nic,omitempty"`
	Timestamp  string                           `json:"timestamp,omitempty"`
	PictureURL string                           `json:"pictureURL,omitemtpy"`

	//payload
	Content   string                           `json:"content"`
	Graph     Graph                            `json:"graph,omitempty"`
	RoomUsers []Person                         `json:"roomUsers,omitempty"`
	MediaSession MediaSession                  `json:"mediaSession,omitempty"`
}

func (endpoint *Endpoint) url() (string) {
	return endpoint.protocol + "://" + endpoint.host + ":" + endpoint.port
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if ( strings.Contains(r.URL.Path, "/session") ) {
		log.Println("Main: Set path ", r.URL.Path)
		r.URL.Path = "/"
	} else if ( strings.Contains(r.URL.Path, "/user") ) {
		//log.Println("Serve ", DocumentRoot+r.URL.Path)
		fp := path.Join(DocumentRoot + r.URL.Path)
		http.ServeFile(w, r, fp)
		return
	} else if ( strings.Contains(r.URL.Path, "/test") ) {
		//log.Println("Serve ", DocumentRoot+r.URL.Path)
		fp := path.Join(DocumentRoot + r.URL.Path)
		http.ServeFile(w, r, fp)
		return
	} else if ( strings.Contains(r.URL.Path, "/css") ) {
		//log.Println("Serve ", DocumentRoot+r.URL.Path)
		fp := path.Join(DocumentRoot + r.URL.Path)
		http.ServeFile(w, r, fp)
		return
	} else if ( strings.Contains(r.URL.Path, "/js") ) {
		//log.Println("Serve ", DocumentRoot+r.URL.Path)
		fp := path.Join(DocumentRoot + r.URL.Path)
		http.ServeFile(w, r, fp)
		return
	} else if ( strings.Contains(r.URL.Path, "/images") ) {
		//log.Println("Serve ", DocumentRoot+r.URL.Path)
		fp := path.Join(DocumentRoot + r.URL.Path)
		http.ServeFile(w, r, fp)
		return
	}

	if r.URL.Path != "/" {
		http.Error(w, "Main: Illegal path "+r.URL.Path, 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Main: Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	room := hub.messages["Main"]
	ifs := room.GetAllAsList()
	var msgs []Message
	msgs = make([]Message, len(ifs), len(ifs))
	for i := 0; i < len(ifs); i++ {
		msgs[i] = ifs[i].(Message)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTemplate.Execute(w, struct {
		Host      string
		LoggedIn  string
		LoggedOut string
		Person    Person
		Messages  []Message
		Persons   []Person
		Targets   []GreenBlue
	}{
		Host:      r.Host,
		LoggedIn:  "none",
		LoggedOut: "flex",
		Person:    Person{},
		Messages:  msgs,
		Persons:   _persons.getAllInRoom("Main"),
		Targets:   nil,
	})
}

type GreenBlue struct {
	Color  string
	Target Person
}

func sessionHandler(w http.ResponseWriter, r *http.Request) {

	sess, err := sessionStore.Get(r, sessionName)
	if err != nil {
		log.Println("Main: sessionHandler: Error in getting and verifying coookie ", err)
	}
	token := sess.Values[sessionToken].(string)
	log.Println("session token from cookie ", token)
	var person Person
	var ok bool
	person, ok = _persons.findPersonByToken(token)
	if ! ok {
		log.Println("Main: sessionHandler: User does not exist for token ", person.Token)
		w.Write([]byte("Authorization Failure! User does not exist, The following token is invalid: " + token ))
	}
	room := hub.messages[person.Room]
	ifs := room.GetAllAsList()
	var msgs []Message
	msgs = make([]Message, len(ifs), len(ifs))
	for i := 0; i < len(ifs); i++ {
		msgs[i] = ifs[i].(Message)
	}
	var targets []GreenBlue
	for k, _ := range _publishers[person.UserID] {
		target, ok := _persons.findPersonByUserId(k)
		if ok {
			//	color := updateMPRStatus(person.UserID, target.UserID)
			targets = append(targets, GreenBlue{BLUE, target})
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTemplate.Execute(w, struct {
		Host      string
		LoggedIn  string
		LoggedOut string
		Person    Person
		Messages  []Message
		Persons   []Person
		Targets   []GreenBlue
	}{
		Host:      r.Host,
		LoggedIn:  "flex",
		LoggedOut: "none",
		Person:    person,
		Messages:  msgs,
		Persons:   _persons.getAllInRoom(person.Room),
		Targets:   targets,
	})
}

func NewMux(config *Config, hub *Hub) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", serveHome)
	mux.Handle("/session/", requireLogin(http.HandlerFunc(sessionHandler)))
	mux.Handle("/profile", requireLogin(http.HandlerFunc(profileHandler)))
	mux.Handle("/ProfileUpdate", requireLogin(http.HandlerFunc(updateProfileHandler)))
	mux.Handle("/MainProfile", requireLogin(http.HandlerFunc(mainProfileHandler)))
	mux.Handle("/TargetManager", requireLogin(http.HandlerFunc(TargetManagerHandler)))
	mux.Handle("/RoomManager", requireLogin(http.HandlerFunc(RoomManagerHandler)))
	mux.Handle("/ImageManager", requireLogin(http.HandlerFunc(ImageManager_UploadHandler)))
	mux.Handle("/ImageManagerGetAll", requireLogin(http.HandlerFunc(ImageManger_GetHandler)))
	mux.Handle("/ImageManagerDelete", requireLogin(http.HandlerFunc(ImageManager_DeleteHandler)))

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	mux.HandleFunc("/logout", logoutHandler)

	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  endpoint.url() + "/google/callback",
		Endpoint:     googleOAuth2.Endpoint,
		Scopes:       []string{"profile", "email"},
	}
	// state param cookies require HTTPS by default; disable for localhost development
	stateConfig :=  gologin.DebugOnlyCookieConfig
	mux.Handle("/google/login", google.StateHandler(stateConfig, google.LoginHandler(oauth2Config, nil)))
	mux.Handle("/google/callback", google.StateHandler(stateConfig, google.CallbackHandler(oauth2Config, issueSession(), nil)))


	oauth2ConfigFB := &oauth2.Config{
		ClientID:     config.ClientID_FB,
		ClientSecret: config.ClientSecret_FB,
		RedirectURL:  endpoint.url() + "/facebook/callback",
		Endpoint:     facebookOAuth2.Endpoint,
		//Scopes:       []string{"profile", "email"},
	}
	stateConfigFB := gologin.DefaultCookieConfig
	mux.Handle("/facebook/login", facebook.StateHandler(stateConfigFB, facebook.LoginHandler(oauth2ConfigFB, nil)))
	mux.Handle("/facebook/callback", facebook.StateHandler(stateConfigFB, facebook.CallbackHandler(oauth2ConfigFB, issueSessionFB(), nil)))

	return mux
}

func getCookieAndTokenfromRequest(r *http.Request, onlyTooken bool) (token string, cookie string, err error) {

	if (!onlyTooken) {
		//retrieve encrypted cookie
		cookieInfo, err := r.Cookie(sessionName)
		if (err != nil) {
			return "", "", fmt.Errorf("No cookie found for give cookie name %s detail %s", sessionName, err)
		}
		cookie = cookieInfo.Value
	}

	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		return "", "", fmt.Errorf("Fail to retrieve cookie to create session %s detail %s", sessionName, err)
	}
	atoken, ok := session.Values[sessionToken]
	if !ok {
		return "", "", fmt.Errorf("The sesstion did not contain %s ", sessionToken)
	}
	if atoken != nil {
		token = atoken.(string)
	} else {
		token = ""
	}
	return token, cookie, nil
}

var homepath = ""
var _persons Persons
var hub *Hub
var DocumentRoot string
var endpoint Endpoint
var homeTemplate = template.Must(template.ParseFiles("/var/www/rakuen/home.html"))
var sessionStore *sessions.CookieStore
var _publishers PublishersTargets

func main() {
	//testA()
	//return

	_publishers = make(PublishersTargets)
	_persons = Persons{__pers: make(map[UserId]Person)}
	config := &Config{
		ClientID_FB:      os.Getenv("FACEBOOK_CLIENT_ID"),
		ClientSecret_FB:  os.Getenv("FACEBOOK_CLIENT_SECRET"),
		ClientID:         os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret:     os.Getenv("GOOGLE_CLIENT_SECRET"),
		ChatHost:         os.Getenv("CHAT_HOST"),
		ChatPrivateKey:   os.Getenv("CHAT_PRIVATE_KEY"),
	}
	sessionStore = sessions.NewCookieStore([]byte(config.ChatPrivateKey), nil)
	endpoint = Endpoint{"https", config.ChatHost, "443"}
	dir, _ := os.Getwd()
	DocumentRoot = strings.Replace(dir, " ", "\\ ", -1)
	queue := new(QueueStack)
	var addr = flag.String("addr", ":"+endpoint.port, "http service address")

	// Check if the cert files are available.
	err := httpscerts.Check("fullchain.pem", "privkey.pem")
	//f they are not available, generate new ones.
	if err != nil {
		log.Println("Issuing autosigned Certs..")
		err = httpscerts.Generate("fullchain.pem", "privkey.pem", endpoint.host)
		if err != nil {
			log.Fatal("Error: Couldn't create https certs.")
		}
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

	log.Println("Starting service at ", endpoint.url())
	err = http.ListenAndServeTLS(*addr, "fullchain.pem", "privkey.pem", NewMux(config, hub))
	//err = http.ListenAndServe(*addr, NewMux(config, hub) )
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
