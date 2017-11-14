package annotations

import (
	"fmt"

	"k8s.io/api/core/v1"
)

const (
	NodeAnnotationReboot     = "k8s.node.terminator.reboot" // true as string
	NodeAnnotationFromWindow = "k8s.node.terminator.fromTimeWindow"
	NodeAnnotationToWindow   = "k8s.node.terminator.toTimeWindow"
)

var (
	errMissingNodeAnnotation = fmt.Errorf("missing annotation %v", NodeAnnotationReboot)
	errNodeAnnotationNotTrue = fmt.Errorf("annotation %v not set true", NodeAnnotationReboot)
)

// CheckAnnotationsExists docs
func CheckAnnotationsExists(node *v1.Node) error {
	a := node.GetAnnotations()
	if _, ok := a[NodeAnnotationReboot]; ok {
		if a[NodeAnnotationReboot] == "true" {
			return nil
		}
		return errNodeAnnotationNotTrue
	}
	return errMissingNodeAnnotation
}
