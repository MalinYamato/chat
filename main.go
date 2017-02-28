

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
	"encoding/json"

	"fmt"
	"strconv"

)

type Config struct {
	ClientID       string
	ClientSecret   string
	ChatHost       string
	ChatPrivateKey string
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

func serveHome(w http.ResponseWriter, r *http.Request) {

	if ( strings.Contains(r.URL.Path, "/session") ) {
		log.Println(" Set path ", r.URL.Path)
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
	}{
		Host:      r.Host,
		LoggedIn:  "none",
		LoggedOut: "flex",
		Person:    Person{},
		Messages:  msgs,

	})

}

func sessionHandler(w http.ResponseWriter, r *http.Request) {

	sess, err := sessionStore.Get(r, sessionName)
	if err != nil {
		log.Println("Error in getting and verifying coookie ", err)
	}

	token := sess.Values[sessionToken].(string)

	log.Println("session token from cookie ", token)

	person, ok := Persons[token]
	if !ok {
		log.Println("sessionHandler: User does not exist for token ", person.Token)
		w.Write([]byte("Authorization Failure! User does not exist, The following token is invalid: " + token ))
	}

	room := hub.messages[person.Room]
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
	}{
		Host:      r.Host,
		LoggedIn:  "flex",
		LoggedOut: "none",
		Person:    person,
		Messages:  msgs,
	})

}

type ProfileRequest struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Date struct {
	Year int
	Month int
	Day int
}

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
	UserID            string        `json:"userId,omitempty"`
	Token             string        `json:"token,omitempty"`
	Room              string        `json:"room"`
}

type Status struct {
	Status string `json:"status"`
	Detail string `json:"detail"`
}

func profileHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("User requested a profile")

	var request Person;
	if r.Method == "POST" {

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&request)
		if err != nil {
			log.Println("ERR> ", err)
		}
		defer r.Body.Close()
		log.Printf("%s\n", request.UserID)

		var person Person
		var ok = false
		for k, v := range Persons {
			fmt.Printf("key[%s] value[%s]\n", k, v)
			if v.UserID == request.UserID {
				person = v
				person.Token = "secret"
				ok = true
				break;
			}
			ok = false
			fmt.Printf("key[%s] value[%s]\n", k, v)
		}

		if ok == false {
			log.Println("Person not foond for ID: ")
		}

		data, err := json.Marshal(person)
		if err != nil {
			panic(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)

	} else {
		log.Println("Unknown method ", r.Method)

	}
}

var LANGUAGES = []string{"English", "Finnish", "Same", "Swedish", "German", "French", "Spannish", "Italian", "Portogese", "Russian", "Chinese", "Japanese", "Korean", "Thai" }
var ORIENTATION = []string{"Straight", "Gay", "Lesbian", "BiSexual", "ASexual"}
var GENDER = []string{"Female", "Male", "TranssexualF", "TranssexualM", "CrossDresser", "None"}

func mainProfileHandler(w http.ResponseWriter, r *http.Request) {

	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		log.Println("Client: Call to sessionStore.Get returned ", err)
		return
	}

	if session == nil {
		log.Println("Client: returned session was nil")
		return
	}

	token := session.Values[sessionToken].(string)

	p, _ := Persons[token]
	t := template.New("fieldname example")
	t = template.Must(template.ParseFiles("profile.html"))


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

func updateProfileHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("UpdateProfile called ", r.Method)
	r.ParseForm()

	var status Status

	if r.Method == "POST" {

		session, err := sessionStore.Get(r, sessionName)
		if err != nil {
			log.Println("Client: Call to sessionStore.Get returned ", err)
			status.Status = "Error"
			status.Detail = "Failed to get a valid cookie!"
		} else if session == nil {
			log.Println("Client: returned session was nil")
			status.Status = "Error"
			status.Detail = "The session is not valid!"

		} else {

			token := session.Values[sessionToken].(string)

			p := Persons[token]

			p.FirstName = r.Form.Get("FirstName")
			p.LastName = r.Form.Get("LastName")
			p.Gender = r.Form.Get("Gender")
			p.Country = r.Form.Get("Country")
			p.Town = r.Form.Get("Town")
			p.Nic = r.Form.Get("Nic")
			p.Profession = r.Form.Get("Profession")
			p.Education = r.Form.Get("Education")
			p.SexualOrientation = r.Form.Get("SexualOrientation")
			p.Description = r.Form.Get("Description")
			p.BirthDate.Year, _ =  strconv.Atoi(r.Form.Get("BirthYear"))
			p.BirthDate.Month, _ =  strconv.Atoi(r.Form.Get("BirthMonth"))
			p.BirthDate.Day, _ =  strconv.Atoi(r.Form.Get("BirthDay"))

			fmt.Printf("%+v\n", r.Form)
			productsSelected := r.Form["Language"]
			log.Println(Contains(productsSelected, "English"))

			for i := 0; i < len(LANGUAGES); i++ {
				if Contains(r.Form["Language"], LANGUAGES[i]) {
					p.Languages[LANGUAGES[i]] = "checked"
				}
			}

			if p.Keep == false {
				status.Status = "New"
				status.Detail = "The profile was successfully created!"
			} else {
				status.Status = "Updated"
				status.Detail = "The profile was successfully updated!"
			}

			p.Keep = true
			Persons[token] = p

			log.Println("Name ", r.Form["Gender"])
		}
	}

	data, err := json.Marshal(status)
	if err != nil {
		panic(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

	//http.Redirect(w, r, "/session", http.StatusFound)
}

func NewMux(config *Config, hub *Hub) *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/", serveHome)
	mux.Handle("/session/", requireLogin(http.HandlerFunc(sessionHandler)))
	mux.Handle("/profile", requireLogin(http.HandlerFunc(profileHandler)))
	mux.Handle("/ProfileUpdate", requireLogin(http.HandlerFunc(updateProfileHandler)))
	mux.Handle("/MainProfile", requireLogin(http.HandlerFunc(mainProfileHandler)))

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
var sessionStore *sessions.CookieStore
func main() {

	config := &Config{
		ClientID:      os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret:  os.Getenv("GOOGLE_CLIENT_SECRET"),
		ChatHost:      os.Getenv("CHAT_HOST"),
		ChatPrivateKey : os.Getenv("CHAT_PRIVATE_KEY"),
	}

	sessionStore = sessions.NewCookieStore([]byte(config.ChatPrivateKey), nil)

	endpoint = Endpoint{"https", config.ChatHost, "443"}
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

	Persons = make(PersonsMAP)

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
	err = http.ListenAndServeTLS(*addr, "cert.pem", "key.pem", NewMux(config, hub))
	//err = http.ListenAndServe(*addr, NewMux(config, hub) )
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
