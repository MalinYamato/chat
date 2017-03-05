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

type Endpoint struct {
	protocol string
	host     string
	port     string
}

type Date struct {
	Year  string `json:"year"`
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

type PublishRequest struct {
	Op  string   `json:"op"`
	Ids []string `json:"ids"`
}

type PublishRequestResponse struct {
	Op     string       `json:"op"`
	Status Status      `json:"status"`
	Person Person      `json:"person"`
}

var (
	LANGUAGES   = []string{"English", "Finnish", "Same", "Swedish", "German", "French", "Spannish", "Italian", "Portogese", "Russian", "Chinese", "Japanese", "Korean", "Thai" }
	ORIENTATION = []string{"Straight", "Gay", "Lesbian", "BiSexual", "ASexual"}
	GENDER      = []string{"Female", "Male", "TranssexualF", "TranssexualM", "CrossDresser", "None"}
)

func (endpoint *Endpoint) url() (string) {
	return endpoint.protocol + "://" + endpoint.host + ":" + endpoint.port
}

func serveHome(w http.ResponseWriter, r *http.Request) {

	if ( strings.Contains(r.URL.Path, "/session") ) {
		log.Println("Main: Set path ", r.URL.Path)
		r.URL.Path = "/"
	} else if ( strings.Contains(r.URL.Path, "/images") ) {
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
		Persons:   _persons.getAllLoggedIn(),
		Targets:   nil,
	})
}

type GreenBlue struct {
	Color string
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
	for k, _ := range _publishers[person.UserID].Targets {
		target, ok := _persons.findPersonByUserId(k)
		if ok {
			color := updateMPRStatus(person.UserID, target.UserID)
			targets = append(targets, GreenBlue{color,target})
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
		Persons:   _persons.getAllLoggedIn(),
		Targets:   targets,
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

func updateMPRStatus(clientID string, targetID string) string {
	MPRStatus := BLUE
	var two int = 0
	client, ok := _publishers[clientID]
	if ok {
		if _, ok := client.Targets[targetID]; ok == true {
			two++
		}
	}
	target, ok := _publishers[targetID]
	if ok {
		if _, ok := target.Targets[clientID]; ok == true {
			two++
		}
	}
	if two == 2 {
		MPRStatus = GREEN // target and client are sending messages to each other, they have formed a Multicast Private Room
	}

	targets := make(Targets)
	targets[targetID] = true
	targets[clientID] = true
	hub.multicast <- Message{"UpdateTarget", "", "", clientID, "", targets, timestamp(), "", MPRStatus }
	return MPRStatus
}

func TargetManagerHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var request PublishRequest
	var MPRStatus string
	response := PublishRequestResponse{"RequestResponse", Status{}, Person{}}
	if r.Method == "POST" {
		var client Person
		var ok bool
		token, _, err := getCookieAndTokenfromRequest(r, true)
		if err != nil {
			response.Status = Status{ERROR, err.Error()}
		} else {
			client, ok = _persons.findPersonByToken(token)
			if ! ok {
				response.Status = Status{ERROR, err.Error()}
			} else {
				decoder := json.NewDecoder(r.Body)
				err = decoder.Decode(&request)
				if err != nil {
					log.Println("Json decoder error> ", err.Error())
					panic(err)
				}
				log.Println(request)
				targetID := request.Ids[0]
				log.Println("target", targetID)
				target, ok := _persons.findPersonByUserId(targetID)
				if ! ok {
					log.Printf("Main: Target  not found for UserID %s \n", targetID)
					response.Status = Status{Status: WARNING, Detail: fmt.Sprintf("Receiver not found for UserID %s \n", targetID) }
				} else {
					log.Printf("Main: Profile request for Target %s UserID %s token %s \n", target.Email, target.UserID, target.Token)
					publisher, ok := _publishers[client.UserID]
					if request.Op == "RemoveTarget" {
						if ok && len(publisher.Targets) >= 1 {
							log.Println("Remove a Target")
							delete(publisher.Targets, request.Ids[0])
						}
						if ok && len(publisher.Targets) < 1 {
							delete(_publishers, client.UserID)
						}
					} else if request.Op == "AddTarget" {
						if ! ok {
							publisher = Publisher{client.UserID, make(Targets)}
						}
						publisher.Targets[target.UserID] = true
						log.Printf("Receiver %s added \n", request.Ids[0])
						_publishers[client.UserID] = publisher
					}
					MPRStatus = updateMPRStatus(client.UserID, target.UserID)
					response.Status = Status{SUCCESS, MPRStatus}
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
			p.BirthDate.Year = r.Form.Get("BirthYear")
			p.BirthDate.Month = r.Form.Get("BirthMonth")
			p.BirthDate.Month = r.Form.Get("BirthDay")
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
				status.Detail = " A new profile was successfully created! <br> Public key: " + p.UserID + " <br>Private Key: " + p.Token + " <br>(used for secure broadcasts)"
			} else {
				status.Status = "Updated"
				status.Detail = "The profile was successfully updated! <br> Public key: " + p.UserID + " <br>Private Key: " + p.Token + " <br>(used for secure broadcasts)"
			}
			p.Keep = true
			_persons.Save(p)
		}
	}
	data, err := json.Marshal(status)
	if err != nil {
		panic(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func NewMux(config *Config, hub *Hub) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", serveHome)
	mux.Handle("/session/", requireLogin(http.HandlerFunc(sessionHandler)))
	mux.Handle("/profile", requireLogin(http.HandlerFunc(profileHandler)))
	mux.Handle("/ProfileUpdate", requireLogin(http.HandlerFunc(updateProfileHandler)))
	mux.Handle("/MainProfile", requireLogin(http.HandlerFunc(mainProfileHandler)))
	mux.Handle("/TargetManager", requireLogin(http.HandlerFunc(TargetManagerHandler)))
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

type Targets map[string]bool

type Publisher struct {
	UserId  string
	Targets Targets
}
type Publishers map[string]Publisher

var _persons Persons
var hub *Hub
var DocumentRoot string
var endpoint Endpoint
var homeTemplate = template.Must(template.ParseFiles("home.html"))
var sessionStore *sessions.CookieStore
var _publishers Publishers

func main() {

	_publishers = make(Publishers)
	_persons = Persons{__pers: make(map[string]Person)}

	config := &Config{
		ClientID:       os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret:   os.Getenv("GOOGLE_CLIENT_SECRET"),
		ChatHost:       os.Getenv("CHAT_HOST"),
		ChatPrivateKey: os.Getenv("CHAT_PRIVATE_KEY"),
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
