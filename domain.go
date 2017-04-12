package ppe

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Domain is the main domain type
type Domain struct {
	Organization *Organization
	Name         string
	Destination  string
	Failover     string
	Relay        bool
	Active       bool
}

// NewDomain is the creation type passed to CreateDomain
type NewDomain struct {
	DomainName  string `json:"domain_name"`
	Destination string `json:"destination"`
	Failover    string `json:"failover,omitempty"`
	IsRelay     int    `json:"is_relay,omitempty"`
}

// Domain retrieves a domain from an organization
func (org *Organization) Domain(domain string) (*Domain, error) {
	var d domainsResponse
	err := org.PPE.get(fmt.Sprintf("/domains/%s/%s", org.PrimaryDomain, domain), &d)
	if err != nil {
		return &Domain{}, err
	}
	return domainFromDomainResponse(org, d.Domains[0]), nil
}

// Domain retrieves a domain without knowing it's organization
func (ppe *PPE) Domain(domain string) (*Domain, error) {
	var d domainsResponse
	err := ppe.get(fmt.Sprintf("/domains/%s", domain), &d)
	if err != nil {
		return &Domain{}, err
	}
	drs := d.Domains
	var dom *Domain
	for _, dr := range drs {
		if dr.DomainName == domain {
			org, err := ppe.Organization(domain)
			if err != nil {
				return &Domain{}, err
			}
			dom = domainFromDomainResponse(org, dr)
			break
		}
	}
	return dom, nil
}

// Domains retrieves the domains associated with an organization
func (org *Organization) Domains() ([]*Domain, error) {
	var d domainsResponse
	err := org.PPE.get(fmt.Sprintf("/domains/%s", org.PrimaryDomain), &d)
	if err != nil {
		return []*Domain{}, err
	}
	drs := d.Domains
	doms := make([]*Domain, len(org.domains))
	for _, dr := range drs {
		for i, dom := range org.domains {
			if dr.DomainName == dom {
				doms[i] = domainFromDomainResponse(org, dr)
			}
		}
	}
	return doms, nil
}

type domCreationResponse struct {
	TotalCreated int                     `json:"total_created"`
	FailResults  []domCreationFailResult `json:"fail_results"`
}

type domCreationFailResult struct {
	Result domCreationResult `json:"result"`
}

type domCreationResult struct {
	ResultCode int    `json:"result_code"`
	Message    string `json:"message"`
}

// CreateDomain creates a new domain for the organization
func (org *Organization) CreateDomain(newDom NewDomain) error {
	var (
		r domCreationResponse
		b bytes.Buffer
	)
	newDomL := []NewDomain{newDom}
	err := json.NewEncoder(&b).Encode(newDomL)
	if err != nil {
		return err
	}
	err = org.PPE.post(fmt.Sprintf("/domains/%s", org.PrimaryDomain), &b, &r)
	if err != nil {
		return err
	}
	if len(r.FailResults) > 0 {
		errs := make([]string, len(r.FailResults))
		for i, fr := range r.FailResults {
			errs[i] = fr.Result.Message
		}
		return errors.New(strings.Join(errs, ", "))
	}
	return nil
}

func domainFromDomainResponse(org *Organization, res domainResponse) *Domain {
	return &Domain{
		Organization: org,
		Name:         res.DomainName,
		Destination:  res.Destination,
		Failover:     res.Failover,
		Relay:        res.IsRelay != 0,
		Active:       res.IsRelay != 0,
	}
}

type domainsResponse struct {
	Domains []domainResponse `json:"message"`
}

type domainResponse struct {
	DomainName  string `json:"domain_name"`
	Destination string `json:"destination"`
	Failover    string `json:"failover"`
	IsRelay     int    `json:"is_relay"`
	IsActive    int    `json:"is_active"`
}
