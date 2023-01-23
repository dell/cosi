package model

import "encoding/xml"

type AlertPolicy struct {
	XMLName              xml.Name             `xml:"alert_policy"`
	PolicyName           string               `json:"policyName" xml:"policyName"`
	MetricType           string               `json:"metricType" xml:"metricType"`
	MetricName           string               `json:"metricName" xml:"metricName"`
	CreatedBy            string               `json:"createdBy" xml:"createdBy"`
	IsEnabled            string               `json:"isEnabled" xml:"isEnabled"`
	IsPerInstanceMetric  string               `json:"isPerInstanceMetric" xml:"isPerInstanceMetric"`
	Period               int                  `json:"period" xml:"period"`
	PeriodUnits          string               `json:"periodUnits" xml:"periodUnits"`
	DatapointsToConsider int                  `json:"datapointsToConsider" xml:"datapointsToConsider"`
	DatapointsToAlert    int                  `json:"datapointsToAlert" xml:"datapointsToAlert"`
	Statistic            string               `json:"statistic" xml:"statistic"`
	Operator             string               `json:"operator" xml:"operator"`
	Condition            AlertPolicyCondition `json:"condition" xml:"condition"`
}

type AlertPolicyCondition struct {
	ThresholdUnits string `json:"thresholdUnits,omitempty" xml:"thresholdUnits,omitempty"`
	ThresholdValue string `json:"thresholdValue,omitempty" xml:"thresholdValue,omitempty"`
	SeverityType   string `json:"severityType,omitempty" xml:"severityType,omitempty"`
}

// AlertPolicies is a list of alert policies
type AlertPolicies struct {
	// XMLName is the name of the xml tag used XML marshalling
	XMLName xml.Name `json:"alert_policies" xml:"alert_policies"`

	// Items is the list of alert policies
	Items []AlertPolicy `json:"alert_policy" xml:"alert_policy"`

	// MaxBuckets is the maximum number of alert policies requested in the listing
	MaxPolicies int `json:"MaxPolicies,omitempty" xml:"MaxPolicies"`

	// NextMarker is a reference object to receive the next set of alert policies
	NextMarker string `json:"next_marker,omitempty" xml:"next_marker,omitempty"`

	// Filter is a string query used to limit the returned alert policies in the
	// listing
	Filter string `json:"Filter,omitempty" xml:"Filter,omitempty"`

	// NextPageLink is a hyperlink to the next page in the alert policy listing
	NextPageLink string `json:"next_page_link,omitempty" xml:"next_page_link,omitempty"`
}
