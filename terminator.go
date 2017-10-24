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
	for {
		select {
		case _ = <-stopCh:
			log.Info("stopping updater runner")
			return
		}
	}
}
