package status

import (
	"net/http"

	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/model"
	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/rest/client"
)

// Status is a REST implementation of the Status interface
type Status struct {
	Client client.RemoteCaller
}

// GetRebuildStatus implements the status interface
func (b *Status) GetRebuildStatus(objStoreName, ssPodName, ssPodNameSpace, level string, params map[string]string) (*model.RebuildInfo, error) {
	// URL Example: https://10.240.117.5:4443/vdc/recovery-status/devices/youmin-test-1-ss-0.youmin-test-1-ss.dellemc-globalmarcu-domain-c10.svc.cluster.local/levels/1
	requestURL := "vdc/recovery-status/devices/" + ssPodName + "." +
		objStoreName + "-ss." + ssPodNameSpace + ".svc.cluster.local/levels/" + level
	req := client.Request{
		Method:      http.MethodGet,
		Path:        requestURL,
		ContentType: client.ContentTypeJSON,
		Params:      params,
	}
	rebuildInfo := &model.RebuildInfo{}
	err := b.Client.MakeRemoteCall(req, rebuildInfo)
	if err != nil {
		return nil, err
	}

	return rebuildInfo, nil
}
