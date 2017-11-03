package main

import (
	"k8s.io/api/core/v1"
)

const (
	nodeRoleLabel = "kubernetes.io/role"
)

func checkIfMaster(node *v1.Node) bool {
	labels := node.GetLabels()
	if _, ok := labels[nodeRoleLabel]; ok {
		if labels[nodeRoleLabel] == "master" {
			return true
		}
		if labels[nodeRoleLabel] == "node" {
			return false
		}
	}
	return false
}
