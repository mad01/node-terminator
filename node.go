package main

import (
	"fmt"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

func setNodeUnschedulable(nodename string, client *kubernetes.Clientset) error {
	patch := `{"spec":{"unschedulable":true}}`
	_, err := client.Core().Nodes().Patch(nodename, types.StrategicMergePatchType, []byte(patch))
	if err != nil {
		return fmt.Errorf("failed to set node as unschedulable: %v", err.Error())
	}
	return nil
}

// func setNodeSchedulable(nodename string, client *kubernetes.Clientset) error {
// 	patch := `{"spec":{"$patch":"delete", "unschedulable":true}}`
// 	_, err := client.Core().Nodes().Patch(nodename, types.StrategicMergePatchType, []byte(patch))
// 	if err != nil {
// 		return fmt.Errorf("failed to set node as unschedulable: %v", err.Error())
// 	}
// 	return nil
// }
