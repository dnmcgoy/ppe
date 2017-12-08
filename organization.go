package ppe

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

// Organization is the main organization type
type Organization struct {
	PPE           *PPE
	Name          string
	PrimaryDomain string
	UserLicenses  int
	domains       []string
}

// Organization retrieves the organization given by the name
func (ppe *PPE) Organization(domain string) (*Organization, error) {
	var o orgResource
	err := ppe.get(fmt.Sprintf("/orgs/%s", domain), &o)
	if err != nil {
		return &Organization{}, err
	}
	return orgFromOrgResource(ppe, o), nil
}

// Organizations retrieves the organizations registered under an organization
func (org *Organization) Organizations() ([]*Organization, error) {
	var os []orgResource
	err := org.PPE.get(fmt.Sprintf("/orgs/%s/orgs", org.PrimaryDomain), &os)
	if err != nil {
		return []*Organization{}, err
	}
	orgs := make([]*Organization, len(os))
	for i, o := range os {
		orgs[i] = orgFromOrgResource(org.PPE, o)
	}
	return orgs, nil
}

// NewOrganization is the creation type passed to CreateOrganization
type NewOrganization struct {
	// Required fields
	Name         string         `json:"name"`
	AdminUser    NewUser        `json:"admin_user"`
	UserLicenses int            `json:"user_licences"`
	Domains      []NewOrgDomain `json:"domains"`

	// Optional fields
	PrimaryDomain     string `json:"primary_domain,omitempty"`
	WWW               string `json:"www,omitempty"`
	Address           string `json:"address,omitempty"`
	Postcode          string `json:"postcode,omitempty"`
	Country           string `json:"country,omitempty"`
	Phone             string `json:"phone,omitempty"`
	LicencingPackage  string `json:"licencing_package,omitempty"` // Defaults to beginner
	AccountTemplateID string `json:"account_template_id,omitempty"`
}

// NewOrgDomain is used when creating domains together with an organization. We
// need this type because the proofpoint API is inconsistent.
type NewOrgDomain struct {
	Name       string   `json:"name"`
	Transports []string `json:"transports"`
}

type orgCreationResponse struct {
	TotalCreated int                     `json:"total_created"`
	FailResults  []orgCreationFailResult `json:"fail_results"`
}

type orgCreationFailResult struct {
	Result orgCreationResult `json:"result"`
}

type orgCreationResult struct {
	ResultCode int    `json:"result_code"`
	Message    string `json:"message"`
}

// CreateOrganization creates a new organization associated with this
// organization.
func (org *Organization) CreateOrganization(newOrg NewOrganization) error {
	var (
		r orgCreationResponse
		b bytes.Buffer
	)
	err := json.NewEncoder(&b).Encode(newOrg)
	if err != nil {
		return err
	}
	err = org.PPE.post(fmt.Sprintf("/orgs/%s/orgs", org.PrimaryDomain), &b, &r)
	if err != nil {
		return err
	}
	if len(r.FailResults) >= 1 {
		return errors.New(r.FailResults[0].Result.Message)
	}
	return nil
}

func orgFromOrgResource(ppe *PPE, res orgResource) *Organization {
	domains := make([]string, len(res.Domains))
	for i, dr := range res.Domains {
		domains[i] = dr.Name
	}
	return &Organization{
		PPE:           ppe,
		Name:          res.Name,
		PrimaryDomain: res.PrimaryDomain,
		UserLicenses:  res.UserLicenses,
		domains:       domains,
	}
}

type orgResource struct {
	PrimaryDomain        string              `json:"primary_domain"` // Required
	Name                 string              `json:"name"`           // Required
	Type                 string              `json:"type,omitempty"`
	WWW                  string              `json:"www,omitempty"`
	Address              string              `json:"address,omitempty"`
	Postcode             string              `json:"postcode,omitempty"`
	Country              string              `json:"country,omitempty"`
	Phone                string              `json:"phone,omitempty"`
	ActiveUsers          int                 `json:"active_users,omitempty"`
	LicensingPackage     string              `json:"licensing_package,omitempty"` // Defaults to beginner
	IsBeginnerPlus       bool                `json:"is_beginner_plus,omitempty"`
	BeginnerPlusEnabled  bool                `json:"beginner_plus_enabled,omitempty"`
	UserLicenses         int                 `json:"user_licences"` // Required
	OnTrial              int                 `json:"on_trial,omitempty"`
	WhenRenewal          string              `json:"when_renewal,omitempty"`
	OutgoingServers      []string            `json:"outgoing_servers,omitempty"`
	WhiteListSenders     []string            `json:"white_list_senders,omitempty"`
	BlackListSenders     []string            `json:"black_list_senders,omitempty"`
	AVBypassSenders      []string            `json:"av_bypass_senders,omitempty"`
	EXEBypassEnabled     int                 `json:"exe_bypass_enabled,omitempty"`
	SMTPDiscoveryEnabled int                 `json:"smtp_discovery_enabled"`
	LDAPUrl              string              `json:"ldap_url,omitempty"`
	LDAPUsername         string              `json:"ldap_username,omitempty"`
	LDAPBasedn           string              `json:"ldap_basedn,omitempty"`
	AdminUser            userResource        `json:"admin_user"` // Required
	IsActive             int                 `json:"isactive,omitempty"`
	Domains              []orgDomainResource `json:"domains"` // At least one required
	AccountTemplateID    int                 `json:"account_template_id,omitempty"`
}

type orgDomainResource struct {
	Name string `json:"name"`
}
