package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/model"
)

// Client is a REST client that communicates with the ObjectScale management API
type ServiceClient struct {
	// Endpoint is the URL of the management API
	Endpoint string `json:"endpoint"`

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

	HTTPClient  *http.Client
	token       string
	authRetries int

	// Should X-EMC-Override header be added into the request
	OverrideHeader bool
}

var _ RemoteCaller = (*ServiceClient)(nil)

func NewServiceClient(endpoint, gateway, namespace, podName, sharedSecret, objectscaleId string, h *http.Client, overrideHdr bool) *ServiceClient {
	return &ServiceClient{
		Endpoint:       endpoint,
		Gateway:        gateway,
		Namespace:      namespace,
		PodName:        podName,
		SharedSecret:   sharedSecret,
		ObjectScaleID:  objectscaleId,
		HTTPClient:     h,
		OverrideHeader: overrideHdr,
	}
}

func (c *ServiceClient) login() error {
	// urn:osc:{ObjectScaleId}:{ObjectStoreId}:service/{ServiceNameId}
	serviceUrn := fmt.Sprintf("urn:osc:%s:%s:service/%s", c.ObjectScaleID, "", c.PodName)
	// B64-{ObjectScaleId},{ObjectStoreId},{ServiceK8SNamespace},{ServiceNameId}
	userNameRaw := fmt.Sprintf("%s,%s,%s,%s", c.ObjectScaleID, "", c.Namespace, c.PodName)
	userNameEncoded := base64.StdEncoding.EncodeToString([]byte(userNameRaw))
	userName := "B64-" + userNameEncoded
	// current time in milliseconds (rounded to nearest 30 seconds)
	timeFactor := time.Now().UTC().Round(30*time.Second).UnixNano() / int64(time.Millisecond)

	data := serviceUrn + strconv.FormatInt(timeFactor, 10)
	h := hmac.New(sha256.New, []byte(c.SharedSecret))
	if _, wrr := h.Write([]byte(data)); wrr != nil {
		return fmt.Errorf("server error: problem writing hmac sha256 %w", wrr)
	}
	password := base64.StdEncoding.EncodeToString(h.Sum(nil))

	u, err := url.Parse(c.Gateway)
	if err != nil {
		return err
	}
	u.Path = "/mgmt/serviceLogin"
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(userName, password)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	if err = HandleResponse(resp); err != nil {
		return err
	}
	c.token = resp.Header.Get("X-SDS-AUTH-TOKEN")
	if c.token == "" {
		return fmt.Errorf("server error: login failed")
	} else {
		c.authRetries = 0
	}
	return nil
}

func (c *ServiceClient) isLoggedIn() bool {
	return c.token != ""
}

// MakeRemoteCall executes an API request against the client endpoint, returning
// the object body of the response into a response object
// NOTE: this is WET not DRY, as the same code is copied for Client
func (c *ServiceClient) MakeRemoteCall(r Request, into interface{}) error {
	var (
		obj []byte
		err error
		q   = url.Values{}
	)
	switch r.ContentType {
	case ContentTypeXML:
		obj, err = xml.Marshal(r.Body)
	case ContentTypeJSON:
		if raw, ok := r.Body.(json.RawMessage); ok {
			obj, err = raw.MarshalJSON()
		} else {
			obj, err = json.Marshal(r.Body)
		}
	default:
		return errors.New("invalid content-type")
	}
	if err != nil {
		return err
	}
	u, err := url.Parse(c.Endpoint)
	if err != nil {
		return err
	}
	u.Path = r.Path
	if r.Params != nil {
		for key, value := range r.Params {
			q.Add(key, value)
		}
	}
	u.RawQuery = q.Encode()
	req, err := http.NewRequest(r.Method, u.String(), bytes.NewBuffer(obj))
	if err != nil {
		return err
	}
	if !c.isLoggedIn() {
		if err = c.login(); err != nil {
			return err
		}
	}
	req.Header.Add("Accept", r.ContentType)
	req.Header.Add("Content-Type", r.ContentType)
	req.Header.Add("Accept", "application/xml")
	req.Header.Add("X-SDS-AUTH-TOKEN", c.token)
	if c.OverrideHeader {
		req.Header.Add("X-EMC-Override", "true")
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return err
	}

	switch {
	case resp.StatusCode == http.StatusUnauthorized:
		if c.authRetries < AuthRetriesMax {
			c.authRetries += 1
			c.token = ""
			return c.MakeRemoteCall(r, into)
		}
		return errors.New(strings.ToLower(resp.Status))
	case resp.StatusCode > 399:
		ecsError := &model.Error{}
		switch r.ContentType {
		case ContentTypeXML:
			if err = xml.Unmarshal(body, ecsError); err != nil {
				return err
			}
		case ContentTypeJSON:
			if err = json.Unmarshal(body, ecsError); err != nil {
				return err
			}
		}
		return fmt.Errorf("%s: %s", ecsError.Description, ecsError.Details)
	case into == nil:
		// No errors found, and no response object defined, so just return
		// without error
		return nil
	default:
		if len(body) == 0 {
			return nil
		}
		switch r.ContentType {
		case ContentTypeXML:
			if err = xml.Unmarshal(body, into); err != nil {
				return err
			}
		case ContentTypeJSON:
			if err = json.Unmarshal(body, into); err != nil {
				return err
			}
		}
	}
	return nil
}
