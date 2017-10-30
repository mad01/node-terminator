package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
)

func TestCheckAnnotationsExists(t *testing.T) {
	annotations := map[string]string{
		nodeAnnotation: "true",
	}
	n := &v1.Node{}
	n.SetName("foo")

	assert.NotNil(t, checkAnnotationsExists(n))

	n.SetAnnotations(annotations)
	assert.Nil(t, checkAnnotationsExists(n))
}
