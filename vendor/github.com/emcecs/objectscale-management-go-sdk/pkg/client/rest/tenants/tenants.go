package tenants

import (
	"fmt"
	"net/http"

	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/api"
	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/model"
	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/rest/client"
)

// Tenants is a REST implementation of the Tenants interface
type Tenants struct {
	Client client.RemoteCaller
}

var _ api.TenantsInterface = &Tenants{}

// List implements the tenants interface
func (t *Tenants) List(params map[string]string) (*model.TenantList, error) {
	req := client.Request{
		Method:      http.MethodGet,
		Path:        "/object/tenants",
		ContentType: client.ContentTypeXML,
		Params:      params,
	}
	tenantList := &model.TenantList{}
	err := t.Client.MakeRemoteCall(req, tenantList)
	if err != nil {
		return nil, err
	}
	return tenantList, nil
}

// Get implements the tenants interface
func (t *Tenants) Get(tenantID string, params map[string]string) (*model.Tenant, error) {
	req := client.Request{
		Method:      http.MethodGet,
		Path:        fmt.Sprintf("object/tenants/tenant/%s", tenantID),
		ContentType: client.ContentTypeXML,
		Params:      params,
	}
	tenant := &model.Tenant{}
	err := t.Client.MakeRemoteCall(req, tenant)
	if err != nil {
		return nil, err
	}
	return tenant, nil
}

// Create implements the tenants interface
func (t *Tenants) Create(payload model.TenantCreate) (*model.Tenant, error) {
	req := client.Request{
		Method:      http.MethodPost,
		Path:        "object/tenants/tenant/",
		ContentType: client.ContentTypeXML,
		Body:        payload,
	}
	tenant := &model.Tenant{}
	err := t.Client.MakeRemoteCall(req, tenant)
	if err != nil {
		return nil, err
	}
	return tenant, nil
}

// Delete implements the tenants interface
func (t *Tenants) Delete(tenantID string) error {
	req := client.Request{
		Method:      http.MethodPost,
		Path:        fmt.Sprintf("object/tenants/tenant/%s/delete/", tenantID),
		ContentType: client.ContentTypeXML,
	}
	tenant := &model.Tenant{}
	err := t.Client.MakeRemoteCall(req, tenant)
	if err != nil {
		return err
	}
	return nil
}

// Update implements the tenants interface
func (t *Tenants) Update(payload model.TenantUpdate, tenantID string) error {
	req := client.Request{
		Method:      http.MethodPut,
		Path:        fmt.Sprintf("object/tenants/tenant/%s/", tenantID),
		ContentType: client.ContentTypeXML,
		Body:        payload,
	}
	tenant := &model.Tenant{}
	err := t.Client.MakeRemoteCall(req, tenant)
	if err != nil {
		return err
	}
	return nil
}

// GetQuota implements the tenants interface
func (t *Tenants) GetQuota(tenantID string, params map[string]string) (*model.TenantQuota, error) {
	req := client.Request{
		Method:      http.MethodGet,
		Path:        fmt.Sprintf("object/tenants/tenant/%s/quota", tenantID),
		ContentType: client.ContentTypeXML,
		Params:      params,
	}
	quota := &model.TenantQuota{}
	err := t.Client.MakeRemoteCall(req, quota)
	return quota, err
}

// DeleteQuota implements the tenants interface
func (t *Tenants) DeleteQuota(tenantID string) error {
	req := client.Request{
		Method:      http.MethodDelete,
		Path:        fmt.Sprintf("object/tenants/tenant/%s/quota", tenantID),
		ContentType: client.ContentTypeXML,
	}
	quota := &model.TenantQuota{}
	err := t.Client.MakeRemoteCall(req, quota)
	return err
}

// SetQuota implements the tenants interface
func (t *Tenants) SetQuota(tenantID string, payload model.TenantQuotaSet) error {
	req := client.Request{
		Method:      http.MethodPut,
		Path:        fmt.Sprintf("object/tenants/tenant/%s/quota", tenantID),
		ContentType: client.ContentTypeXML,
		Body:        payload,
	}
	quota := &model.TenantQuota{}
	err := t.Client.MakeRemoteCall(req, quota)
	return err
}
