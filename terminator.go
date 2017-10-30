package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

func newTerminatorEvent() *TerminatorEvent {
	t := TerminatorEvent{}
	return &t
}

// TerminatorEvent event information needed to terminate a node
type TerminatorEvent struct {
	nodename string
}

func newTerminator(client *kubernetes.Clientset) *Terminator {
	t := Terminator{
		events: make(chan TerminatorEvent),
		client: client,
	}
	return &t
}

// Terminator handles node terminate events and handles the lifetime of the event
type Terminator struct {
	events chan TerminatorEvent
	client *kubernetes.Clientset
}

// Run terminator
func (t *Terminator) Run(stopCh chan struct{}) {
	for {
		select {
		case event := <-t.events:
			err := t.terminate(&event)
			if err != nil {
				log.Errorf("failed to terminate node %v %v", event.nodename, err.Error())
			}
		case _ = <-stopCh:
			log.Info("stopping updater runner")
			return
		}
	}
}

func (t *Terminator) terminate(e *TerminatorEvent) error {
	// TODO: implement drain node handling / eveicting of all pods on that node
	// TODO: implement wait for graceperiod before doing force terminate of nodes

	// set node ot no schedule
	log.Infof("terminator get event %v", e.nodename)
	err := setNodeUnschedulable(e.nodename, t.client)
	if err != nil {
		return fmt.Errorf("failed to patch node %v", err.Error())
	}
	// drain node

	// terninate node
	// TODO: implement actuall node termination (only for worker nodes) master should be skipped
	ec2Client := newEC2()
	err = ec2Client.awsTerminateInstance(e.nodename)
	if err != nil {
		return fmt.Errorf("failed to terminate node %v %v", e.nodename, err.Error())
	}

	// wait for new node (sleep)
	// release and go to next
	return nil
}
