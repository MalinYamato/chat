package main

import (

"fmt"
"io/ioutil"

	"net/http"

	"log"
	"github.com/satori/go.uuid"
	"github.com/dghubble/gologin/google"
)

const (
	sessionName    = "secure.krypin.xyz"
	sessionSecret  = "secure,krypin.xyz secret key developer"
	sessionUserKey = "googleID"
)

// Config configures the main ServeMux.
type Person struct {
	Nic string		`json:"nic,omitempty"`
	FirstName string	`json:"firstName,omitempty"`
	LastName string		`json:"lastName,omitempty"`
	Email string		`json:"email,omitempty"`
	Gender string		`json:"gender,omitempty"`
	Location string         `json:"location,omitempty"`
	PictureURL string	`json:"pictureURL,omitempty"`
	GoogleID string		`json:"googleID,omitempty"`
	UserID string           `json:"googleID,omitempty"`
	Token string		`json:"token,omitempty"`
}


func issueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		googleUser, err := google.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// remove possible old cookies
		if isAuthenticated(req) {
			log.Println("There was an old cookie. Removing it")
			sessionStore.Destroy(w, sessionName)
		}
		// issue a new cookie
		session := sessionStore.New(sessionName)
		session.Values[sessionUserKey] = googleUser.Id
		session.Save(w)

		secret := uuid.NewV4()   // used as a secret to verify identity of users who sends websocket messages from the brwoser to the server
		userID := uuid.NewV4()   // used to identify a user to all other users, not a secret.

		Persons[secret.String()] = Person{"null",googleUser.GivenName,googleUser.FamilyName,googleUser.Email,googleUser.Gender,googleUser.Locale,googleUser.Picture,googleUser.Id,userID.String(), secret.String()}
		hub.broadcast <- Message{ "NewUser","null", userID.String(), googleUser.GivenName, googleUser.Picture, googleUser.Gender, "NULL"  }
		log.Println("Successful Login ", googleUser.Email)
		http.Redirect(w, req, "/session/?id=" + secret.String(), http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

// welcomeHandler shows a welcome message and login button.
func welcomeHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	if isAuthenticated(req) {
		http.Redirect(w, req, "/profile", http.StatusFound)
		return
	}
	page, _ := ioutil.ReadFile("home.html")
	fmt.Fprintf(w, string(page))
}


// logoutHandler destroys the session on POSTs and redirects to home.
func logoutHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("User Loggeed out, Delete",req.Method)
	if req.Method == "POST" {
		req.ParseForm()
		token := req.Form["token"]
		person, ok := Persons[token[0]]
		if ok == true {
			hub.broadcast <- Message{"ExitUser", "出ました",person.UserID, person.FirstName, person.PictureURL, person.Gender, "出室、またね　" + person.FirstName + " " + person.LastName}
			delete(Persons, token[0])
		}
		sessionStore.Destroy(w, sessionName)
	}
	http.Redirect(w, req, "/", http.StatusFound)
}

// requireLogin redirects unauthenticated users to the login route.
func requireLogin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if !isAuthenticated(req) {
			http.Redirect(w, req, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

// isAuthenticated returns true if the user has a signed session cookie.
func isAuthenticated(req *http.Request) bool {
	if _, err := sessionStore.Get(req, sessionName); err == nil {
		return true
	}
	return false
}

