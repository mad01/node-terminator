# node-terminator
[![Docker Repository on Quay](https://quay.io/repository/mad01/node-terminator/status "Docker Repository on Quay")](https://quay.io/repository/mad01/node-terminator)

a service to manage the upgrade lifecyckle of k8s nodes (supports aws)

the terminator looks for a annotation on the nodes. The terminator then based on the allowed concurrent terminations and wait inbetween terminations it will terminate N nodes at one time and track the progress of the nodes until all the nodes in a cluster have been terminated that have that annotation. 

the assumption for this to work is that you are using something like `kops` that have a remote state that is newer then the current state of the nodes in the cluster. for example the k8s version is `1.7.4` and you like to update all nodes to `1.7.9`. and you have updated the remote state that kops nodeup uses this would allow a cluster update of the worker nodes by a service running in the cluster


the terminator does the following
* looks if node exists in done nodes
* checks that max concurrent terminations is <= current
* checks it's not currently in progress
* checks that the needed annotation exits with the right value
* checks that it's not a master node
* adds the node to in progress list
* sets node to unschedulable
* drains the node and waits (follows the disruption budget)
* waits for N time to allow pods to terminate
* terminates the aws instance
* waits for N time to allow new node to join the cluster
* adds node to done nodes/removed from in progress
* done > looks for next termination event

## outside scope of service 

it's outside of the scope of this service to set the annotations on the nodes. The idea is that this will be managed by a external service or by using kubectl. a external service could set the annotations on nodes that does not match the desired state. You could for example only have nodes annotated during a time window i.e a maintenance window and only have nodes terminate slowly in that window

## annotation
`k8s.node.terminator.reboot` set to `true` as string in the nodes annotations
`k8s.node.terminator.fromTimeWindow` `k8s.node.terminator.toTimeWindow` if both from/to is set terminations will be done in that window in format `hh:mm AM/PM` . the termination window will mean that termination will be allowed every day in the window. The time is take with `time.Now()` and should base the timezone on local time. if the window not set or missing reboots will be done anytime

```yaml
apiVersion: v1
kind: Node
metadata:
  annotations:
    k8s.node.terminator.reboot: "true"
    k8s.node.terminator.fromTimeWindow: "02:01 AM"
    k8s.node.terminator.toTimeWindow: "05:01 AM"

```

## testing 
add annotation to node for termination
```
kubectl annotate node ip-172-20-118-57.eu-west-1.compute.internal k8s.node.terminator.reboot="true"
```


