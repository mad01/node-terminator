package main

import (
	"fmt"

	"k8s.io/api/core/v1"
)

const (
	nodeAnnotation = "k8s.node.terminator.reboot" // true as string
)

var (
	errMissingNodeAnnotation = fmt.Errorf("missing annotation %v", nodeAnnotation)
	errNodeAnnotationNotTrue = fmt.Errorf("annotation %v not set true", nodeAnnotation)
)

func checkAnnotationsExists(node *v1.Node) error {
	annotations := node.GetAnnotations()
	if _, ok := annotations[nodeAnnotation]; ok {
		if annotations[nodeAnnotation] == "true" {
			return nil
		}
		return errNodeAnnotationNotTrue
	}
	return errMissingNodeAnnotation
}
