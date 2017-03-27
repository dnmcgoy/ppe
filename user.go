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
	var u userResource
	err := org.PPE.get(fmt.Sprintf("/orgs/%s/users/%s", org.PrimaryDomain, email), &u)
	if err != nil {
		return &User{}, err
	}
	return userFromUserResource(org, u), nil
}

// Users retreives all the users of an organization
func (org *Organization) Users() ([]*User, error) {
	var us []userResource
	err := org.PPE.get(fmt.Sprintf("/orgs/%s/users", org.PrimaryDomain), &us)
	if err != nil {
		return []*User{}, err
	}
	users := make([]*User, len(us))
	for i, u := range us {
		users[i] = userFromUserResource(org, u)
	}
	return users, nil
}

// NewUser creates a new user on the organization
func (org *Organization) NewUser(user NewUser) error {
	return nil
}

func userFromUserResource(org *Organization, res userResource) *User {
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

// NewUser is the type used for user creation
type NewUser struct {
	PrimaryEmail string   `json:"primary_email"`
	Firstname    string   `json:"firstname,omitempty"`
	Lastname     string   `json:"lastname"`
	AliasEmails  []string `json:"alias_emails"`
	Type         string   `json:"type"`
}

type userResource struct {
	Firstname        string   `json:"firstname,omitempty"`
	Surname          string   `json:"surname,omitempty"`
	PrimaryEmail     string   `json:"primary_email"`
	AliasEmails      []string `json:"alias_emails,omitempty"`
	WhiteListSenders []string `json:"white_list_senders,omitempty"`
	BlackListSenders []string `json:"black_list_senders,omitempty"`
	IsActive         int      `json:"isactive,omitempty"`
	Type             string   `json:"type,omitempty"`
}
