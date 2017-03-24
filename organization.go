package ppe

import "fmt"

// Organization is the main organization type
type Organization struct {
	PPE           *PPE
	Name          string
	PrimaryDomain string
	domains       []string
}

// Organization retrieves the organization given by the name
func (ppe *PPE) Organization(domain string) (*Organization, error) {
	var o orgResponse
	err := ppe.get(fmt.Sprintf("/orgs/%s", domain), &o)
	if err != nil {
		return &Organization{}, err
	}
	return orgFromOrgResponse(ppe, o), nil
}

// Organizations retrieves the organizations registered under an organization
func (org *Organization) Organizations() ([]*Organization, error) {
	var os []orgResponse
	err := org.PPE.get(fmt.Sprintf("/orgs/%s/orgs", org.PrimaryDomain), &os)
	if err != nil {
		return []*Organization{}, err
	}
	orgs := make([]*Organization, len(os))
	for i, o := range os {
		orgs[i] = orgFromOrgResponse(org.PPE, o)
	}
	return orgs, nil
}

func orgFromOrgResponse(ppe *PPE, res orgResponse) *Organization {
	domains := make([]string, len(res.Domains))
	for i, dr := range res.Domains {
		domains[i] = dr.Name
	}
	return &Organization{
		PPE:           ppe,
		Name:          res.Name,
		PrimaryDomain: res.PrimaryDomain,
		domains:       domains,
	}
}

type orgResponse struct {
	PrimaryDomain        string              `json:"primary_domain"`
	Name                 string              `json:"name"`
	Type                 string              `json:"type"`
	WWW                  string              `json:"www"`
	Address              string              `json:"address"`
	Postcode             string              `json:"postcode"`
	Country              string              `json:"country"`
	Phone                string              `json:"phone"`
	ActiveUsers          int                 `json:"active_users"`
	LicensingPackage     string              `json:"licensing_package"`
	IsBeginnerPlus       bool                `json:"is_beginner_plus"`
	BeginnerPlusEnabled  bool                `json:"beginner_plus_enabled"`
	UserLicenses         int                 `json:"user_licences"`
	OnTrial              int                 `json:"on_trial"`
	WhenRenewal          string              `json:"when_renewal"`
	OutgoingServers      []string            `json:"outgoing_servers"`
	WhiteListSenders     []string            `json:"white_list_senders"`
	BlackListSenders     []string            `json:"black_list_senders"`
	AVBypassSenders      []string            `json:"av_bypass_senders"`
	EXEBypassEnabled     int                 `json:"exe_bypass_enabled"`
	SMTPDiscoveryEnabled int                 `json:"smtp_discovery_enabled"`
	LDAPUrl              string              `json:"ldap_url"`
	LDAPUsername         string              `json:"ldap_username"`
	LDAPBasedn           string              `json:"ldap_basedn"`
	AdminUser            userResponse        `json:"admin_user"`
	IsActive             int                 `json:"isactive"`
	Domains              []orgDomainResponse `json:"domains"`
}

type orgDomainResponse struct {
	Name string `json:"name"`
}
