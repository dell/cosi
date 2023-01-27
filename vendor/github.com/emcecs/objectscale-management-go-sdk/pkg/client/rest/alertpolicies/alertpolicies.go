package alertpolicies

import (
	"net/http"
	"path"

	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/api"
	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/model"
	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/rest/client"
)

// AlertPolicies is a REST implementation of the AlertPolicies interface
type AlertPolicies struct {
	Client client.RemoteCaller
}

var _ api.AlertPoliciesInterface = &AlertPolicies{}

// Get implements the AlertPolicy interface
func (ap *AlertPolicies) Get(policyName string) (*model.AlertPolicy, error) {
	req := client.Request{
		Method:      http.MethodGet,
		Path:        path.Join("vdc", "alertpolicy", policyName),
		ContentType: client.ContentTypeXML,
	}
	alertpolicy := &model.AlertPolicy{}
	err := ap.Client.MakeRemoteCall(req, alertpolicy)
	if err != nil {
		return nil, err
	}
	return alertpolicy, nil
}

// List implements the AlertPolicy interface
func (ap *AlertPolicies) List(params map[string]string) (*model.AlertPolicies, error) {
	req := client.Request{
		Method:      http.MethodGet,
		Path:        path.Join("vdc", "alertpolicy", "list"),
		ContentType: client.ContentTypeXML,
		Params:      params,
	}
	alertpolicies := &model.AlertPolicies{}
	err := ap.Client.MakeRemoteCall(req, alertpolicies)
	if err != nil {
		return nil, err
	}

	return alertpolicies, nil
}

// Create implements the AlertPolicy interface
func (ap *AlertPolicies) Create(payload model.AlertPolicy) (*model.AlertPolicy, error) {
	req := client.Request{
		Method:      http.MethodPost,
		Path:        path.Join("vdc", "alertpolicy"),
		ContentType: client.ContentTypeXML,
		Body:        &payload,
	}
	alertpolicy := &model.AlertPolicy{}
	err := ap.Client.MakeRemoteCall(req, alertpolicy)
	if err != nil {
		return nil, err
	}
	return alertpolicy, nil
}

// Delete implements the AlertPolicy interface
func (ap *AlertPolicies) Delete(policyName string) error {
	req := client.Request{
		Method:      http.MethodDelete,
		Path:        path.Join("vdc", "alertpolicy", policyName),
		ContentType: client.ContentTypeXML,
	}
	alertpolicy := &model.AlertPolicy{}
	err := ap.Client.MakeRemoteCall(req, alertpolicy)
	if err != nil {
		return err
	}
	return nil
}

// Update implements the AlertPolicy interface
func (ap *AlertPolicies) Update(payload model.AlertPolicy, policyName string) (*model.AlertPolicy, error) {
	req := client.Request{
		Method:      http.MethodPut,
		Path:        path.Join("vdc", "alertpolicy", policyName),
		ContentType: client.ContentTypeXML,
		Body:        &payload,
	}
	alertpolicy := &model.AlertPolicy{}
	err := ap.Client.MakeRemoteCall(req, alertpolicy)
	if err != nil {
		return nil, err
	}
	return alertpolicy, nil
}
