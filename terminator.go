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
	worker       int
}

// GetWorker returns the go rutine num of the worker
func (t *TerminatorEvent) GetWorker() string {
	return fmt.Sprintf("worker %d", t.worker)
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
		doneNodes:          set.New(),
	}
	return &t
}

// Terminator handles node terminate events and handles the lifetime of the event
type Terminator struct {
	events                 chan TerminatorEvent
	client                 *kubernetes.Clientset
	eviction               *Eviction
	activeTerminations     *set.Set
	doneNodes              *set.Set
	concurrentTerminations int
}

// Run terminator
func (t *Terminator) Run(stopCh chan struct{}) {
	// start up your worker threads based on concurrentTerminations
	for i := 0; i < t.concurrentTerminations; i++ {
		go t.worker(i, stopCh)
	}

}

func (t *Terminator) worker(num int, stopCh chan struct{}) {
	for {
		select {
		case event := <-t.events:
			event.worker = num
			err := t.terminate(&event)
			if err != nil {
				log.Errorf("%v failed to terminate node %v %v",
					event.GetWorker(), event.nodename, err.Error(),
				)
			}
		case _ = <-stopCh:
			log.Infof("stopping worker")
			return
		}
	}
}

func (t *Terminator) terminate(event *TerminatorEvent) error {
	// set node ot no schedule
	log.Infof("%v terminator got event %v", event.GetWorker(), event.nodename)
	t.activeTerminations.Add(event.nodename)

	// drain node
	log.Infof("%v starting drain of node %v", event.GetWorker(), event.nodename)
	err := t.eviction.DrainNode(event.nodename)
	if err != nil {
		t.activeTerminations.Remove(event.nodename)
		log.Errorf("%v failed to drain node %v %v",
			event.GetWorker(),
			event.nodename,
			err.Error(),
		)
	}

	// terninate node
	log.Infof("%v starting termination of ec2 instance with node name %v",
		event.GetWorker(),
		event.nodename,
	)
	ec2Client := newEC2()
	err = ec2Client.awsTerminateInstance(event.nodename)
	if err != nil {
		t.activeTerminations.Remove(event.nodename)
		return fmt.Errorf("%v failed to terminate node %v %v",
			event.GetWorker(),
			event.nodename,
			err.Error(),
		)
	}

	// wait for new node (sleep)
	log.Infof("%v waiting for %s node to terminate %v",
		event.GetWorker(),
		event.waitInterval,
		event.nodename,
	)
	time.Sleep(event.waitInterval)

	// release and go to next
	t.activeTerminations.Remove(event.nodename)
	t.doneNodes.Add(event.nodename)
	log.Infof("%v done terminating node %v", event.GetWorker(), event.nodename)

	return nil
}
