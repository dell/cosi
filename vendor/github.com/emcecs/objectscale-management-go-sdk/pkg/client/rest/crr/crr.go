package crr

import (
	"net/http"
	"path"

	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/model"
	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/rest/client"
)

// CRR is a REST implementation of the CRR interface
type CRR struct {
	Client client.RemoteCaller
}

// PauseReplication implements the CRR interface
func (c *CRR) PauseReplication(destObjectScale string, destObjectStore string, params map[string]string) error {
	req := client.Request{
		Method:      http.MethodPost,
		Path:        path.Join("replication", "control", destObjectScale, destObjectStore, "pause"),
		ContentType: client.ContentTypeXML,
		Params:      params,
	}
	return c.Client.MakeRemoteCall(req, nil)
}

// SuspendReplication implements the CRR interface
func (c *CRR) SuspendReplication(destObjectScale string, destObjectStore string, params map[string]string) error {
	req := client.Request{
		Method:      http.MethodPost,
		Path:        path.Join("replication", "control", destObjectScale, destObjectStore, "suspend"),
		ContentType: client.ContentTypeXML,
		Params:      params,
	}
	return c.Client.MakeRemoteCall(req, nil)
}

// ResumeReplication implements the CRR interface
func (c *CRR) ResumeReplication(destObjectScale string, destObjectStore string, params map[string]string) error {
	req := client.Request{
		Method:      http.MethodPost,
		Path:        path.Join("replication", "control", destObjectScale, destObjectStore, "resume"),
		ContentType: client.ContentTypeXML,
		Params:      params,
	}
	return c.Client.MakeRemoteCall(req, nil)
}

// ResumeReplication implements the CRR interface
func (c *CRR) UnthrottleReplication(destObjectScale string, destObjectStore string, params map[string]string) error {
	req := client.Request{
		Method:      http.MethodPost,
		Path:        path.Join("replication", "control", destObjectScale, destObjectStore, "unthrottle"),
		ContentType: client.ContentTypeXML,
		Params:      params,
	}
	return c.Client.MakeRemoteCall(req, nil)
}

// ThrottleReplication implements the CRR interface
func (c *CRR) ThrottleReplication(destObjectScale string, destObjectStore string, params map[string]string) error {
	req := client.Request{
		Method:      http.MethodPost,
		Path:        path.Join("replication", "control", destObjectScale, destObjectStore, "throttle"),
		ContentType: client.ContentTypeXML,
		Params:      params,
	}
	return c.Client.MakeRemoteCall(req, nil)
}

// Get implements the CRR interface
func (c *CRR) Get(destObjectScale string, destObjectStore string, params map[string]string) (*model.CRR, error) {
	req := client.Request{
		Method:      http.MethodGet,
		Path:        path.Join("replication", "control", destObjectScale, destObjectStore),
		ContentType: client.ContentTypeXML,
		Params:      params,
	}
	config := &model.CRR{}
	err := c.Client.MakeRemoteCall(req, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
