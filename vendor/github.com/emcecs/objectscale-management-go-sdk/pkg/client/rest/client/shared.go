package client

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/model"
)

const (
	// AuthRetriesMax is the maximum number of times the client will attempt to
	// login before returning an error
	AuthRetriesMax = 3

	// ContentTypeXML
	ContentTypeXML = "application/xml"

	// ContentTypeJSON
	ContentTypeJSON = "application/json"
)

// Request is an ObjectScale API request wrapper
type Request struct {
	// Method is the method of REST API request
	Method string

	// Path is the path of REST API request
	Path string

	// Body is the body of REST API request
	Body interface{}

	// ContentType is the body of REST API request
	ContentType string

	// Params are the parameters of the REST API request
	Params map[string]string
}

// RemoteCaller interface is used to create backend calls.
// into represents type, _into_ which data will be unmarshalled.
// Naming follows Effective Go naming convention https://go.dev/doc/effective_go#interface-names
type RemoteCaller interface {
	MakeRemoteCall(r Request, into interface{}) error
}

func HandleResponse(resp *http.Response) error {
	if resp.StatusCode > 399 {
		switch resp.Body {
		case nil:
			switch {
			case resp.Status != "":
				return fmt.Errorf("server error: %s", strings.ToLower(resp.Status))
			case resp.StatusCode != 0:
				return fmt.Errorf("server error: status code %d", resp.StatusCode)
			}
		default:
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("server errror: %s", strings.ToLower(resp.Status))
				return err
			}
			apiError := &model.Error{}
			err = xml.Unmarshal(body, apiError)
			if err != nil {
				return errors.New(string(body))
			}
			switch {
			case apiError.Code == 1004:
				return errors.New("server error: not found")
			default:
				return fmt.Errorf("server error: %s", strings.ToLower(apiError.Description))
			}
		}
	}
	return nil
}
