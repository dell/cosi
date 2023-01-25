/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"net/http"

	rest "k8s.io/client-go/rest"
	v1alpha1 "sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"
	"sigs.k8s.io/container-object-storage-interface-api/client/clientset/versioned/scheme"
)

type ObjectstorageV1alpha1Interface interface {
	RESTClient() rest.Interface
	BucketsGetter
	BucketAccessesGetter
	BucketAccessClassesGetter
	BucketClaimsGetter
	BucketClassesGetter
}

// ObjectstorageV1alpha1Client is used to interact with features provided by the objectstorage.k8s.io group.
type ObjectstorageV1alpha1Client struct {
	restClient rest.Interface
}

func (c *ObjectstorageV1alpha1Client) Buckets() BucketInterface {
	return newBuckets(c)
}

func (c *ObjectstorageV1alpha1Client) BucketAccesses(namespace string) BucketAccessInterface {
	return newBucketAccesses(c, namespace)
}

func (c *ObjectstorageV1alpha1Client) BucketAccessClasses() BucketAccessClassInterface {
	return newBucketAccessClasses(c)
}

func (c *ObjectstorageV1alpha1Client) BucketClaims(namespace string) BucketClaimInterface {
	return newBucketClaims(c, namespace)
}

func (c *ObjectstorageV1alpha1Client) BucketClasses() BucketClassInterface {
	return newBucketClasses(c)
}

// NewForConfig creates a new ObjectstorageV1alpha1Client for the given config.
// NewForConfig is equivalent to NewForConfigAndClient(c, httpClient),
// where httpClient was generated with rest.HTTPClientFor(c).
func NewForConfig(c *rest.Config) (*ObjectstorageV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	httpClient, err := rest.HTTPClientFor(&config)
	if err != nil {
		return nil, err
	}
	return NewForConfigAndClient(&config, httpClient)
}

// NewForConfigAndClient creates a new ObjectstorageV1alpha1Client for the given config and http client.
// Note the http client provided takes precedence over the configured transport values.
func NewForConfigAndClient(c *rest.Config, h *http.Client) (*ObjectstorageV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientForConfigAndClient(&config, h)
	if err != nil {
		return nil, err
	}
	return &ObjectstorageV1alpha1Client{client}, nil
}

// NewForConfigOrDie creates a new ObjectstorageV1alpha1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *ObjectstorageV1alpha1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new ObjectstorageV1alpha1Client for the given RESTClient.
func New(c rest.Interface) *ObjectstorageV1alpha1Client {
	return &ObjectstorageV1alpha1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1alpha1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *ObjectstorageV1alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
