package federatedobjectstores

import (
	"net/http"

	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/api"
	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/model"
	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/rest/client"
)

// FederatedObjectStores is a REST implementation of the FederateObjectStores interface
type FederatedObjectStores struct {
	Client *client.Client
}

var _ api.FederatedObjectStoresInterface = &FederatedObjectStores{}

// List implements the federatedobjectstores interface
func (t *FederatedObjectStores) List(params map[string]string) (*model.FederatedObjectStoreList, error) {
	req := client.Request{
		Method:      http.MethodGet,
		Path:        "/replication/info",
		ContentType: client.ContentTypeXML,
		Params:      params,
	}
	federatedStoreList := &model.FederatedObjectStoreList{}
	err := t.Client.MakeRemoteCall(req, federatedStoreList)
	if err != nil {
		return nil, err
	}
	return federatedStoreList, nil
}
