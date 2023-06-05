// Copyright © 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// AuthRetriesMax is the maximum number of times the client will attempt to
// login before returning an error
const AuthRetriesMax = 3

// Authenticator can perform a Login to the gateway.
type Authenticator interface {
	// IsAuthenticated returns true if the authenticated has been established.  This
	// does not mean the next request is guaranteed to succeed as authentication can
	// become expired.
	IsAuthenticated() bool

	// Login obtains fresh authentication token(s) from the server.
	Login(*http.Client) error

	// Token returns the current authentication token.
	Token() string
}

var _ Authenticator = (*AuthUser)(nil)    // interface guard
var _ Authenticator = (*AuthService)(nil) // interface guard

// AuthService is an in-cluster Authenticator.
type AuthService struct {
	// Gateway is the auth endpoint
	Gateway string `json:"gateway"`

	// SharedSecret is the fedsvc shared secret
	SharedSecret string `json:"sharedSecret"`

	// PodName is the GraphQL Pod name
	PodName string `json:"podName"`

	// Namespace is the GraphQL Namespace name
	Namespace string `json:"namespace"`

	// ObjectScaleID is just that
	ObjectScaleID string `json:"objectScaleID"`

	token string
}

// IsAuthenticated returns true if the authenticated has been established.  This
// does not mean the next request is guaranteed to succeed as authentication can
// become expired.
func (auth *AuthService) IsAuthenticated() bool {
	return auth.token != ""
}

// Login obtains fresh authentication token(s) from the server.
func (auth *AuthService) Login(ht *http.Client) error {
	// urn:osc:{ObjectScaleID}:{ObjectStoreID}:service/{ServiceNameID}
	serviceUrn := fmt.Sprintf("urn:osc:%s:%s:service/%s", auth.ObjectScaleID, "", auth.PodName)
	// B64-{ObjectScaleID},{ObjectStoreID},{ServiceK8SNamespace},{ServiceNameID}
	userNameRaw := fmt.Sprintf("%s,%s,%s,%s", auth.ObjectScaleID, "", auth.Namespace, auth.PodName)
	userNameEncoded := base64.StdEncoding.EncodeToString([]byte(userNameRaw))
	userName := "B64-" + userNameEncoded
	// current time in milliseconds (rounded to nearest 30 seconds)
	timeFactor := time.Now().UTC().Round(30*time.Second).UnixNano() / int64(time.Millisecond)

	data := serviceUrn + strconv.FormatInt(timeFactor, 10)
	h := hmac.New(sha256.New, []byte(auth.SharedSecret))
	if _, wrr := h.Write([]byte(data)); wrr != nil {
		return fmt.Errorf("server error: problem writing hmac sha256 %w", wrr)
	}
	password := base64.StdEncoding.EncodeToString(h.Sum(nil))

	u, err := url.Parse(auth.Gateway)
	if err != nil {
		return err
	}
	u.Path = "/mgmt/serviceLogin"
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(userName, password)
	resp, err := ht.Do(req)
	if err != nil {
		return err
	}
	if err = HandleResponse(resp); err != nil {
		return err
	}
	auth.token = resp.Header.Get("X-SDS-AUTH-TOKEN")
	if auth.token == "" {
		return fmt.Errorf("server error: login failed")
	}
	return nil
}

// Token returns the current authentication token.
func (auth *AuthService) Token() string {
	return auth.token
}

// AuthUser is an out-of-cluster or username+password based Authenticator.
type AuthUser struct {
	// Gateway is the auth endpoint
	Gateway string `json:"gateway"`

	// Username used to authenticate management user
	Username string `json:"username"`

	// Password used to authenticate management user
	Password string `json:"password"`

	token string
}

// IsAuthenticated returns true if the authenticated has been established.  This
// does not mean the next request is guaranteed to succeed as authentication can
// become expired.
func (auth *AuthUser) IsAuthenticated() bool {
	return auth.token != ""
}

// Login obtains fresh authentication token(s) from the server.
func (auth *AuthUser) Login(ht *http.Client) error {
	u, err := url.Parse(auth.Gateway)
	if err != nil {
		return err
	}
	u.Path = "/mgmt/login"
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(auth.Username, auth.Password)
	resp, err := ht.Do(req)
	if err != nil {
		return err
	}
	if err = HandleResponse(resp); err != nil {
		return err
	}
	auth.token = resp.Header.Get("X-SDS-AUTH-TOKEN")
	if auth.token == "" {
		return fmt.Errorf("server error: login failed")
	}
	return nil
}

// Token returns the current authentication token.
func (auth *AuthUser) Token() string {
	return auth.token
}
