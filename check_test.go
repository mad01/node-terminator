package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCheckIfMaster(t *testing.T) {
	testCases := []struct {
		testName string
		node     *v1.Node
		expected bool
	}{
		{
			testName: "node as type master",
			node: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: "node0",
					Labels: map[string]string{
						"kubernetes.io/role": "master",
					},
				},
				Spec: v1.NodeSpec{
					ProviderID: "node0",
				},
			},
			expected: true,
		},

		{
			testName: "node as type node",
			node: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: "node0",
					Labels: map[string]string{
						"kubernetes.io/role": "node",
					},
				},
				Spec: v1.NodeSpec{
					ProviderID: "node0",
				},
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, checkIfMaster(tc.node))
	}
}
