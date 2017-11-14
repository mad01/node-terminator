package annotations

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCheckAnnotationsExists(t *testing.T) {
	testCases := []struct {
		testName string
		node     *v1.Node
		expected error
	}{
		{
			testName: "node with correct annotations",
			node: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: "node0",
					Annotations: map[string]string{
						NodeAnnotationReboot: "true",
					},
				},
				Spec: v1.NodeSpec{
					ProviderID: "node0",
				},
			},
			expected: nil,
		},

		{
			testName: "node with correct annotations but not set to true",
			node: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: "node0",
					Annotations: map[string]string{
						NodeAnnotationReboot: "false",
					},
				},
				Spec: v1.NodeSpec{
					ProviderID: "node0",
				},
			},
			expected: errNodeAnnotationNotTrue,
		},

		{
			testName: "node with correct annotations",
			node: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: "node0",
				},
				Spec: v1.NodeSpec{
					ProviderID: "node0",
				},
			},
			expected: errMissingNodeAnnotation,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, CheckAnnotationsExists(tc.node))
	}
}
