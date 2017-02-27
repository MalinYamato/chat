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
	sessionToken  =   "SessionToken"
)

// Config configures the main ServeMux.


func issueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		googleUser, err := google.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		secret := uuid.NewV4()   // used as a secret to verify identity of users who sends websocket messages from the brwoser to the server
		userID := uuid.NewV4()   // used to identify a user to all other users, not a secret.

		// remove possible old cookies
		if isAuthenticated(req) {
			log.Println("There was an old cookie. Removing it")
			sessionStore.Destroy(w, sessionName)
		}
		// issue a new cookie
		session := sessionStore.New(sessionName)
		//session.Values[sessionUserKey] = secret
		session.Values[sessionUserKey] = googleUser.Id
		session.Values[sessionToken] = secret.String()
		err = session.Save(w)

		if err != nil {
			log.Println("could not set session ", err)
		}


		/*
	        type Person struct {
			Nic               string        `json:"nic,omitempty"`
			FirstName         string        `json:"firstName,omitempty"`
			LastName          string        `json:"lastName,omitempty"`
			Email             string        `json:"email,omitempty"`
			Gender            string        `json:"gender,omitempty"`
			Town              string        `json:"country,omitempty"`
			Country           string        `json:"town,omitempty"`
			PictureURL        string        `json:"pictureURL,omitempty"`
			SexualOrientation string        `json:"sexualOrienation,omitempty"`
			Languages         map[string]string `json:"Languages,omitempty"`
			Profession        string        `json:"profession,omitempty"`
			Education         string        `json:"education,omitempty"`
			GoogleID          string        `json:"googleId,omitempty"`
			UserID            string        `json:"userId,omitempty"`
			Token             string        `json:"token,omitempty"`
		}
		*/

		Persons[secret.String()] = Person{
			Nic:"",
			FirstName: googleUser.GivenName,
			LastName:googleUser.FamilyName,
			Email: googleUser.Email,
			Gender: googleUser.Gender,
			Town: "",
			Country:googleUser.Locale,
			PictureURL:googleUser.Picture,
			SexualOrientation:"null",
			Languages: map[string]string{},
			Profession: "",
			Education: "",
			GoogleID: googleUser.Id,
			UserID : userID.String(),
			Token: secret.String()}

		hub.broadcast <- Message{ "NewUser","null", userID.String(), googleUser.GivenName, googleUser.Picture, googleUser.Gender, "NULL"  }
		log.Println("Successful Login ", googleUser.Email)
		http.Redirect(w, req, "/session", http.StatusFound)
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
		session, _ := sessionStore.Get(req,sessionName)
		token := session.Values[sessionToken].(string)
		person, ok := Persons[token]
		if ok == true {
			hub.broadcast <- Message{"ExitUser", "出ました",person.UserID, person.FirstName, person.PictureURL, person.Gender, "出室、またね　" + person.FirstName + " " + person.LastName}
			delete(Persons, token)
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
	_, err := sessionStore.Get(req, sessionName);
	if err == nil {
		return true
	}
	log.Println("authentication failed ", err)
	return false
}
