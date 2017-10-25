package main

import log "github.com/sirupsen/logrus"

func newTerminatorEvent() *TerminatorEvent {
	t := TerminatorEvent{}
	return &t
}

// TerminatorEvent event information needed to terminate a node
type TerminatorEvent struct {
	nodename string
}

func newTerminator() *Terminator {
	t := Terminator{
		events: make(chan TerminatorEvent),
	}
	return &t
}

// Terminator handles node terminate events and handles the lifetime of the event
type Terminator struct {
	events chan TerminatorEvent
}

// Run terminator
func (t *Terminator) Run(stopCh chan struct{}) {
	// TODO: implement node no schedule
	// TODO: implement drain node handling / eveicting of all pods on that node
	// TODO: implement actuall node termination (only for worker nodes) master should be skipped
	// TODO: implement wait for graceperiod before doing force terminate of nodes
	for {
		select {
		case event := <-t.events:
			log.Infof("terminator get event %v", event.nodename)
		case _ = <-stopCh:
			log.Info("stopping updater runner")
			return
		}
	}
}
