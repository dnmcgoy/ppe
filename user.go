package ppe

import "fmt"

// User is the main user type
type User struct {
	Organization     *Organization
	Firstname        string
	Surname          string
	Email            string
	Aliases          []string
	WhiteListSenders []string
	BlackListSenders []string
	Active           bool
	Type             string
}

// User retreives a single user from an organization
func (org *Organization) User(email string) (*User, error) {
	var u userResponse
	err := org.PPE.get(fmt.Sprintf("/orgs/%s/users/%s", org.PrimaryDomain, email), &u)
	if err != nil {
		return &User{}, err
	}
	return userFromUserResponse(org, u), nil
}

// Users retreives all the users of an organization
func (org *Organization) Users() ([]*User, error) {
	var us []userResponse
	err := org.PPE.get(fmt.Sprintf("/orgs/%s/users", org.PrimaryDomain), &us)
	if err != nil {
		return []*User{}, err
	}
	users := make([]*User, len(us))
	for i, u := range us {
		users[i] = userFromUserResponse(org, u)
	}
	return users, nil
}

func userFromUserResponse(org *Organization, res userResponse) *User {
	return &User{
		Organization:     org,
		Firstname:        res.Firstname,
		Surname:          res.Surname,
		Email:            res.PrimaryEmail,
		Aliases:          res.AliasEmails,
		WhiteListSenders: res.WhiteListSenders,
		BlackListSenders: res.BlackListSenders,
		Active:           res.IsActive != 0,
		Type:             res.Type,
	}
}

type userResponse struct {
	Firstname        string   `json:"firstname"`
	Surname          string   `json:"surname"`
	PrimaryEmail     string   `json:"primary_email"`
	AliasEmails      []string `json:"alias_emails"`
	WhiteListSenders []string `json:"white_list_senders"`
	BlackListSenders []string `json:"black_list_senders"`
	IsActive         int      `json:"isactive"`
	Type             string   `json:"type"`
}
