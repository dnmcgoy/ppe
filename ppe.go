package ppe

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// PPE is the Proofpoint Essentials api client
type PPE struct {
	apiroot string
	user    string
	pass    string
}

// New creates a new proofpoint API client
func New(apiroot, user, pass string) *PPE {
	return &PPE{
		apiroot: fmt.Sprintf("https://%s/api", apiroot),
		user:    user,
		pass:    pass,
	}
}

// UnauthorizedError is returned when the Proofpoint API responds with 401
type UnauthorizedError struct{}

func (err UnauthorizedError) Error() string {
	return http.StatusText(http.StatusUnauthorized)
}

func (pp PPE) request(method, req string, reqBody io.Reader, res interface{}) error {
	// Build request
	httpReq, err := http.NewRequest(method, fmt.Sprintf("%s%s", pp.apiroot, req), reqBody)
	if err != nil {
		return err
	}
	httpReq.Header.Set("X-User", pp.user)
	httpReq.Header.Set("X-Password", pp.pass)
	if method == "POST" || method == "PUT" {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	// Send request
	httpClient := http.Client{}
	httpRes, err := httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	// Check for HTTP errors
	switch httpRes.StatusCode {
	case http.StatusUnauthorized:
		return UnauthorizedError{}
	}

	// Unmarshal response
	err = json.NewDecoder(httpRes.Body).Decode(res)
	if err != nil {
		return err
	}
	return nil
}

func (pp PPE) get(req string, res interface{}) error {
	return pp.request("GET", req, nil, res)
}

func (pp PPE) post(req string, reqBody io.Reader, res interface{}) error {
	return pp.request("POST", req, reqBody, res)
}

func (pp PPE) put(req string, reqBody io.Reader, res interface{}) error {
	return pp.request("PUT", req, reqBody, res)
}

func (pp PPE) delete(req string, res interface{}) error {
	return pp.request("DELETE", req, nil, res)
}
