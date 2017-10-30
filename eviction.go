package main

import (
	"k8s.io/client-go/kubernetes"
)

func drainNode(nodename string, client *kubernetes.Clientset) error {
	// look at eviction policy
	// client.Core().Pods().Evict
	return nil
}
