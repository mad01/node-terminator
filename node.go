package main

import (
	"encoding/json"
	"fmt"

	"k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/kubernetes"
)

func setNodeUnschedulable(nodename string, client *kubernetes.Clientset) error {
	node, err := client.Core().Nodes().Get(nodename, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get node %v %v", nodename, err.Error())
	}
	return nodeSchedulablePatch(node, true, client)
}

func setNodeSchedulable(nodename string, client *kubernetes.Clientset) error {
	node, err := client.Core().Nodes().Get(nodename, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get node %v %v", nodename, err.Error())
	}

	return nodeSchedulablePatch(node, false, client)
}

func nodeSchedulablePatch(node *v1.Node, schedulable bool, client *kubernetes.Clientset) error {
	// patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, v1.Node{})

	oldData, err := json.Marshal(node)
	if err != nil {
		return fmt.Errorf("failed to Marshal old node %v", err.Error())
	}
	node.Spec.Unschedulable = schedulable
	newData, err := json.Marshal(node)
	if err != nil {
		return fmt.Errorf("failed to Marshal new node %v", err.Error())
	}

	patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, v1.Node{})
	if err != nil {
		return fmt.Errorf("failed to create patch %v", err.Error())
	}

	_, err = client.Core().Nodes().Patch(node.GetName(), types.StrategicMergePatchType, patchBytes)
	if err != nil {
		return fmt.Errorf("failed to set node unschedulable=%v %v", schedulable, err.Error())
	}

	return nil
}
