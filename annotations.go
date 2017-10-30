package main

import (
	"fmt"

	"k8s.io/api/core/v1"
)

const (
	nodeAnnotation = "node.updater.reboot" // true as string
)

var (
	errMissingNodeAnnotation = fmt.Errorf("missing annotation %v", nodeAnnotation)
)

func checkAnnotationsExists(node *v1.Node) error {
	annotations := node.GetAnnotations()
	if _, ok := annotations[nodeAnnotation]; ok {
		return nil
	}
	return errMissingNodeAnnotation
}
