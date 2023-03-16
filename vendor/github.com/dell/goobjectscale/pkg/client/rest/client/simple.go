//
//
//  Copyright © 2021 - 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//       http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//
//

package client

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dell/goobjectscale/pkg/client/model"
)

var _ RemoteCaller = (*Simple)(nil) // interface guard

// Allowed content types.
const (
	ContentTypeXML  = "application/xml"
	ContentTypeJSON = "application/json"
)

// Simple is the simple client definition.
type Simple struct {
	// Endpoint is the URL of the management API
	Endpoint string `json:"endpoint"`

	// Authenticator!=nil means Authenticator.Login will be called to
	// obtain login credentials.
	Authenticator Authenticator

	HTTPClient *http.Client

	// Should X-EMC-Override header be added into the request
	OverrideHeader bool
}

// MakeRemoteCall executes an API request against the client endpoint, returning
// the object body of the response into a response object
// NOTE: this is WET not DRY, as the same code is copied for Client
func (c *Simple) MakeRemoteCall(r Request, into interface{}) error {
	err := r.Validate(c.Endpoint)
	if err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}
	//
	// Unmarshal unmarshals resp.Body into v with respect to returned Content-Type or
	// original requested Content-Type.
	Unmarshal := func(resp *http.Response, v interface{}) error {
		// Use content type sent by server first but fall back to original request content type.
		contentType := resp.Header.Get("Content-Type")
		if contentType == "" {
			contentType = r.ContentType
		}
		//
		// Reading all of resp.Body into memory is inefficient and it's better to send
		// it directly into decoders.  However if body is empty the decoders will
		// receive an EOF.  They can also receive EOF for malformed responses.
		// We need to differentiate between EOF from empty responses and malformed
		// responsed.
		var cw CountWriter
		body := io.TeeReader(resp.Body, &cw)
		//
		// HandleError handles a decoding error.
		// If incoming error is EOF and cw.N == 0 then it's EOF due to empty response
		// and error is discarded if-and-only-if v is non-nil; i.e. it's an error
		// if we got empty body but expected a response.
		HandleError := func(err error) error {
			// v == nil check disabled because some test fixtures are missing proper response and/or some API
			// calls are written strangely -- Tenants.SetQuota for example.
			if (errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF)) && cw.N == 0 /* && v == nil TODO Causes problems -- see note */ {
				return nil
			}
			return err
		}
		switch contentType {
		case ContentTypeJSON:
			decoder := json.NewDecoder(body)
			if err := HandleError(decoder.Decode(v)); err != nil {
				return fmt.Errorf("response: json: %w", err)
			}
		case ContentTypeXML:
			decoder := xml.NewDecoder(body)
			if err := HandleError(decoder.Decode(v)); err != nil {
				return fmt.Errorf("response: xml: %w", err)
			}
		default:
			return fmt.Errorf("response: %s: %w", r.ContentType, ErrContentType)
		}
		return nil
	}
	//
	// Do performs a single http request.
	Do := func() error {
		req, err := r.HTTP(c.Endpoint)
		if err != nil {
			return fmt.Errorf("simple client: %w", err)
		}
		//
		req.Header.Add("Accept", r.ContentType)
		req.Header.Add("Content-Type", r.ContentType)
		req.Header.Add("Accept", "application/xml")
		if c.Authenticator != nil {
			req.Header.Add("X-SDS-AUTH-TOKEN", c.Authenticator.Token())
		}
		if c.OverrideHeader {
			req.Header.Add("X-EMC-Override", "true")
		}
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close() // #nosec
		//
		switch {
		case resp.StatusCode == http.StatusUnauthorized:
			return ErrAuthorization
		case resp.StatusCode > 399:
			var ecsError model.Error
			if err = Unmarshal(resp, &ecsError); err != nil {
				return err
			}
			return fmt.Errorf("%s: %s", ecsError.Description, ecsError.Details)
		case into == nil:
			// No errors found, and no response object defined, so just return
			// without error
			return nil
		default:
			if err = Unmarshal(resp, into); err != nil {
				return err
			}
		}
		return nil
	}
	//
	// If Authenticator is nil then just perform a single request; otherwise
	// perform AuthRetriesMax requests but only if returned error is an authorization
	// error.
	if c.Authenticator == nil {
		return Do()
	}
	if !c.Authenticator.IsAuthenticated() {
		if err := c.Authenticator.Login(c.HTTPClient); err != nil {
			return fmt.Errorf("%w: login: %s", ErrAuthorization, err.Error())
		}
	}
	for tries := 0; tries < AuthRetriesMax; tries++ {
		err := Do()
		switch {
		case errors.Is(err, ErrAuthorization):
			if err = c.Authenticator.Login(c.HTTPClient); err != nil {
				// TODO Depending on how the error is constructed we could potentially
				//      leak credentials here.  Must be careful.
				return fmt.Errorf("%w: retry login", err)
			}
			continue
		default:
			return err
		}
	}
	return fmt.Errorf("%w: exhausted authentication tries", ErrAuthorization)
}
