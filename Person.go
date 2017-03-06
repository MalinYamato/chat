package main

import ()
import "log"

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
