package objmt

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/model"
	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/rest/client"
)

type bucketIdsReqBody struct {
	XMLName xml.Name `xml:"bucket_list"`
	Ids     []string `xml:"id"`
}

type accountIdsReqBody struct {
	XMLName xml.Name `xml:"account_list"`
	Ids     []string `xml:"id"`
}

type storeIdsReqBody struct {
	XMLName xml.Name `xml:"store_list"`
	Ids     []string `xml:"id"`
}

type replicationPairsReqBody struct {
	XMLName      xml.Name         `xml:"replication_list"`
	Replications []replicationIds `xml:"replication"`
}

type replicationIds struct {
	XMLName xml.Name `xml:"replication"`
	Src     string   `xml:"src"`
	Dest    string   `xml:"dest"`
}

func newReplicationIds(ids [][]string) *replicationPairsReqBody {
	ret := &replicationPairsReqBody{}
	ret.Replications = []replicationIds{}
	for _, id := range ids {
		ret.Replications = append(ret.Replications, replicationIds{Src: id[0], Dest: id[1]})
	}
	return ret
}

// Objmt is a REST implementation of the Objmt interface
type Objmt struct {
	Client client.RemoteCaller
}

// GetAccountBillingInfo returns billing info metrics for defined accounts
func (o *Objmt) GetAccountBillingInfo(ids []string, params map[string]string) (*model.AccountBillingInfoList, error) {
	req := client.Request{
		Method:      http.MethodPost,
		Path:        "/object/mt/account/info",
		ContentType: client.ContentTypeXML,
		Body:        accountIdsReqBody{Ids: ids},
		Params:      params,
	}
	ret := &model.AccountBillingInfoList{}
	err := o.Client.MakeRemoteCall(req, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// GetAccountBillingSample returns billing sample (time-window) metrics for defined accounts
func (o *Objmt) GetAccountBillingSample(ids []string, params map[string]string) (*model.AccountBillingSampleList, error) {
	req := client.Request{
		Method:      http.MethodPost,
		Path:        "/object/mt/account/sample",
		ContentType: client.ContentTypeXML,
		Body:        accountIdsReqBody{Ids: ids},
		Params:      params,
	}
	ret := &model.AccountBillingSampleList{}
	err := o.Client.MakeRemoteCall(req, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// GetBucketBillingInfo returns billing info metrics for defined buckets and account
func (o *Objmt) GetBucketBillingInfo(account string, ids []string, params map[string]string) (*model.BucketBillingInfoList, error) {
	// TODO prepare request body with IDs
	req := client.Request{
		Method:      http.MethodPost,
		Path:        fmt.Sprintf("/object/mt/account/%s/bucket/info", account),
		ContentType: client.ContentTypeXML,
		Body:        bucketIdsReqBody{Ids: ids},
		Params:      params,
	}
	ret := &model.BucketBillingInfoList{}
	err := o.Client.MakeRemoteCall(req, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// GetBucketBillingSample returns billing sample (time-window) metrics for defined buckets and account
func (o *Objmt) GetBucketBillingSample(account string, ids []string, params map[string]string) (*model.BucketBillingSampleList, error) {
	req := client.Request{
		Method:      http.MethodPost,
		Path:        fmt.Sprintf("/object/mt/account/%s/bucket/sample", account),
		ContentType: client.ContentTypeXML,
		Body:        bucketIdsReqBody{Ids: ids},
		Params:      params,
	}
	ret := &model.BucketBillingSampleList{}
	err := o.Client.MakeRemoteCall(req, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// GetBucketBillingPerf returns performance metrics for defined buckets and account
func (o *Objmt) GetBucketBillingPerf(account string, ids []string, params map[string]string) (*model.BucketPerfDataList, error) {
	req := client.Request{
		Method:      http.MethodPost,
		Path:        fmt.Sprintf("/object/mt/account/%s/bucket/perf", account),
		ContentType: client.ContentTypeXML,
		Body:        bucketIdsReqBody{Ids: ids},
		Params:      params,
	}
	ret := &model.BucketPerfDataList{}
	err := o.Client.MakeRemoteCall(req, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// GetReplicationInfo returns billing info metrics for defined replication pairs and account
func (o *Objmt) GetReplicationInfo(account string, replicationPairs [][]string, params map[string]string) (*model.BucketReplicationInfoList, error) {
	req := client.Request{
		Method:      http.MethodPost,
		Path:        fmt.Sprintf("/object/mt/account/%s/replication/info", account),
		ContentType: client.ContentTypeXML,
		Body:        newReplicationIds(replicationPairs),
		Params:      params,
	}
	ret := &model.BucketReplicationInfoList{}
	err := o.Client.MakeRemoteCall(req, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// GetReplicationSample returns billing sample (time-window) metrics for defined replication pairs and account
func (o *Objmt) GetReplicationSample(account string, replicationPairs [][]string, params map[string]string) (*model.BucketReplicationSampleList, error) {
	req := client.Request{
		Method:      http.MethodPost,
		Path:        fmt.Sprintf("/object/mt/account/%s/replication/sample", account),
		ContentType: client.ContentTypeXML,
		Body:        newReplicationIds(replicationPairs),
		Params:      params,
	}
	ret := &model.BucketReplicationSampleList{}
	err := o.Client.MakeRemoteCall(req, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// GetStoreBillingInfo returns billing info metrics for object store
func (o *Objmt) GetStoreBillingInfo(params map[string]string) (*model.StoreBillingInfoList, error) {
	req := client.Request{
		Method:      http.MethodGet,
		Path:        "/object/mt/store/info",
		ContentType: client.ContentTypeXML,
		Params:      params,
	}
	ret := &model.StoreBillingInfoList{}
	err := o.Client.MakeRemoteCall(req, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// GetStoreBillingSample returns billing sample (time-window) metrics for object store
func (o *Objmt) GetStoreBillingSample(params map[string]string) (*model.StoreBillingSampleList, error) {
	req := client.Request{
		Method:      http.MethodGet,
		Path:        "/object/mt/store/sample",
		ContentType: client.ContentTypeXML,
		Params:      params,
	}
	ret := &model.StoreBillingSampleList{}
	err := o.Client.MakeRemoteCall(req, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// GetStoreReplicationData returns CRR metrics for defined object stores
func (o *Objmt) GetStoreReplicationData(ids []string, params map[string]string) (*model.StoreReplicationDataList, error) {
	req := client.Request{
		Method:      http.MethodPost,
		Path:        "/object/mt/store/replication",
		ContentType: client.ContentTypeXML,
		Body:        storeIdsReqBody{Ids: ids},
		Params:      params,
	}
	ret := &model.StoreReplicationDataList{}
	err := o.Client.MakeRemoteCall(req, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
