package window

import (
	"fmt"
	"testing"

	"github.com/mad01/k8s-node-terminator/pkg/annotations"
	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestParseTest(t *testing.T) {
	testCases := []struct {
		testName      string
		input         string
		expectedError bool
	}{

		{
			testName:      "time :01",
			input:         "01:01",
			expectedError: true,
		},

		{
			testName:      "time 01:",
			input:         "01:",
			expectedError: true,
		},

		{
			testName:      "time 01:01 AM",
			input:         "01:01 AM",
			expectedError: false,
		},

		{
			testName:      "time 01:01 am",
			input:         "01:01 am",
			expectedError: true,
		},

		{
			testName:      "time 01:01 PM",
			input:         "01:01 PM",
			expectedError: false,
		},

		{
			testName:      "time 1:01",
			input:         "1:01",
			expectedError: true,
		},

		{
			testName:      "time 1:01 AM",
			input:         "1:01 AM",
			expectedError: false,
		},

		{
			testName:      "time 1:01 PM",
			input:         "1:01 PM",
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		if tc.expectedError == true {
			assert.NotNil(t, ParseTest(tc.input), tc.testName, tc.input)
		} else if tc.expectedError == false {
			assert.Nil(t, ParseTest(tc.input), tc.testName, tc.input)
		}
	}

}

func TestMaintainWindow(t *testing.T) {
	testCases := []struct {
		testName string
		node     *v1.Node
		expected bool
	}{
		{
			testName: "in window",
			expected: true,
			node: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: "node0",
					Annotations: map[string]string{
						annotations.NodeAnnotationFromWindow: "05:01 AM",
						annotations.NodeAnnotationToWindow:   "10:01 PM",
						annotations.NodeAnnotationReboot:     "true",
					},
				},
				Spec: v1.NodeSpec{
					ProviderID: "node0",
				},
			},
		},

		{
			testName: "outside window in morning",
			expected: false,
			node: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: "node0",
					Annotations: map[string]string{
						annotations.NodeAnnotationFromWindow: "02:01 AM",
						annotations.NodeAnnotationToWindow:   "03:01 AM",
						annotations.NodeAnnotationReboot:     "true",
					},
				},
				Spec: v1.NodeSpec{
					ProviderID: "node0",
				},
			},
		},

		{
			testName: "outside window in evening",
			expected: false,
			node: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: "node0",
					Annotations: map[string]string{
						annotations.NodeAnnotationFromWindow: "10:01 PM",
						annotations.NodeAnnotationToWindow:   "11:01 PM",
						annotations.NodeAnnotationReboot:     "true",
					},
				},
				Spec: v1.NodeSpec{
					ProviderID: "node0",
				},
			},
		},

		{
			testName: "outside window in evening min 41",
			expected: false,
			node: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: "node0",
					Annotations: map[string]string{
						annotations.NodeAnnotationFromWindow: "10:41 PM",
						annotations.NodeAnnotationToWindow:   "11:41 PM",
						annotations.NodeAnnotationReboot:     "true",
					},
				},
				Spec: v1.NodeSpec{
					ProviderID: "node0",
				},
			},
		},
	}

	for _, tc := range testCases {
		node := fmt.Sprintf("%#v", tc.node)
		window, err := GetMaintenanceWindowFromAnnotations(tc.node)
		assert.Nil(t, err, tc.testName)
		assert.Equal(t, tc.expected, window.InMaintenanceWindow(), tc.testName, node,
			fmt.Sprintf("from time: %v", window.from),
			fmt.Sprintf("to time: %v", window.to),
		)
	}
}
