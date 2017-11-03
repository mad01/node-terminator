package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	set "gopkg.in/fatih/set.v0"
	"k8s.io/client-go/kubernetes"
)

func newTerminatorEvent(nodename string) *TerminatorEvent {
	t := TerminatorEvent{
		nodename:     nodename,
		waitInterval: 1 * time.Minute,
	}
	return &t
}

// TerminatorEvent event information needed to terminate a node
type TerminatorEvent struct {
	nodename     string
	waitInterval time.Duration
}

func newTerminator(kubeconfig string) *Terminator {
	client, err := k8sGetClient(kubeconfig)
	if err != nil {
		panic(fmt.Sprintf("failed to get client: %v", err.Error()))
	}

	t := Terminator{
		events:             make(chan TerminatorEvent),
		client:             client,
		eviction:           newEviction(kubeconfig),
		activeTerminations: set.New(),
	}
	return &t
}

// Terminator handles node terminate events and handles the lifetime of the event
type Terminator struct {
	events             chan TerminatorEvent
	client             *kubernetes.Clientset
	eviction           *Eviction
	activeTerminations *set.Set
}

// Run terminator
func (t *Terminator) Run() {
	for {
		select {
		case event := <-t.events:
			err := t.terminate(&event)
			if err != nil {
				log.Errorf("failed to terminate node %v %v", event.nodename, err.Error())
			}
		}
	}
}

func (t *Terminator) terminate(event *TerminatorEvent) error {
	// set node ot no schedule
	log.Infof("terminator get event %v", event.nodename)
	t.activeTerminations.Add(event.nodename)
	err := setNodeUnschedulable(event.nodename, t.client)
	if err != nil {
		return fmt.Errorf("failed to patch node %v", err.Error())
	}

	// drain node
	log.Infof("starting drain of node %v", event.nodename)
	err = t.eviction.DrainNode(event.nodename)
	if err != nil {
		log.Errorf("failed to drain node %v %v", event.nodename, err.Error())
	}

	// terninate node
	// TODO: implement safeguard to never terminate master node only worker
	log.Infof("starting termination of ec2 instance with node name %v", event.nodename)
	ec2Client := newEC2()
	err = ec2Client.awsTerminateInstance(event.nodename)
	if err != nil {
		return fmt.Errorf("failed to terminate node %v %v", event.nodename, err.Error())
	}

	// wait for new node (sleep)
	log.Infof("waiting for %s node to terminate %v", event.waitInterval, event.nodename)
	time.Sleep(event.waitInterval)

	// release and go to next
	t.activeTerminations.Remove(event.nodename)
	log.Infof("done terminating node %v", event.nodename)

	return nil
}
