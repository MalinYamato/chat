package main

import (

"fmt"
"io/ioutil"
"net/http"
"github.com/dghubble/sessions"
	"log"
	"github.com/satori/go.uuid"
	"github.com/dghubble/gologin/google"
)

const (
	sessionName    = "example-google-app"
	sessionSecret  = "example cookie signing secret"
	sessionUserKey = "googleID"
)



// sessionStore encodes and decodes session data stored in signed cookies
var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

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


// issueSession issues a cookie session after successful Google login
func issueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		googleUser, err := google.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 2. Implement a success handler to issue some form of session
		session := sessionStore.New(sessionName)
		session.Values[sessionUserKey] = googleUser.Id
		session.Save(w)
		secret := uuid.NewV4()
		userID := uuid.NewV4()

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


// profileHandler shows protected user content.
func profileHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, `<p>You are logged in!</p><form action="/logout" method="post"><input type="submit" value="Logout"></form>`)
}

// logoutHandler destroys the session on POSTs and redirects to home.
func logoutHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Delete",req.Method)
	if req.Method == "POST" {
		req.ParseForm()
		token := req.Form["token"]
		person, ok := Persons[token[0]]
		if ok == true {
			hub.broadcast <- Message{"ExitUser", "",person.UserID, person.FirstName, person.PictureURL, person.Gender, "User logged out!"}
			delete(Persons, "token[0]")
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

