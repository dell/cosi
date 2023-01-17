package main

import (
	objectscaleRest "github.com/emcecs/objectscale-management-go-sdk/pkg/client/rest"
	"k8s.io/client-go/kubernetes"
	bucketclientset "sigs.k8s.io/container-object-storage-interface-api/client/clientset/versioned"
)

// place for storing global variables like specs
var (
	clientset    *kubernetes.Clientset
	bucketClient *bucketclientset.Clientset
	objectscale  *objectscaleRest.ClientSet
)
