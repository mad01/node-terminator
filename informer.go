package main

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	lister_v1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"

	"github.com/mad01/node-terminator/pkg/annotations"
	"github.com/mad01/node-terminator/pkg/window"
)

var (
	errorNodeInDone                     = errors.New("node is in doneNodes")
	errorActiveTerminationsLimitReached = errors.New("activeTerminations limit reached")
	errorActiveTerminationsHasNode      = errors.New("activeTerminations has node")
	errorNodeAnnoationsMissing          = errors.New("missing annotations")
	errorNodeMaster                     = errors.New("node is master only supporting node type node")
)

type nodeControllerInput struct {
	waitInterval           time.Duration // TODO: need better name
	updateInterval         time.Duration
	kubeconfig             string
	concurrentTerminations int
}

// manages all nodes and sets annotations on N nodes to reboot at one time
type nodeController struct {
	client     *kubernetes.Clientset
	informer   cache.Controller
	indexer    cache.Indexer
	nodeLister lister_v1.NodeLister
	terminator *Terminator
}

func newNodeController(input *nodeControllerInput) *nodeController {
	client, err := k8sGetClient(input.kubeconfig)
	if err != nil {
		panic(fmt.Sprintf("failed to get client: %v", err.Error()))
	}

	c := &nodeController{
		client: client,
	}

	c.terminator = newTerminator(input.kubeconfig)
	c.terminator.eviction.waitInterval = input.waitInterval
	c.terminator.concurrentTerminations = input.concurrentTerminations

	indexer, informer := cache.NewIndexerInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (runtime.Object, error) {
				return client.Core().Nodes().List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return client.Core().Nodes().Watch(lo)
			},
		},
		// The types of objects this informer will return
		&v1.Node{},
		input.updateInterval,
		// Callback Functions to trigger on add/update/delete
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {},
			UpdateFunc: func(old, new interface{}) {
				node := new.(*v1.Node)
				_ = c.createNodeEvent(node, input)
			},
			DeleteFunc: func(obj interface{}) {},
		},
		cache.Indexers{},
	)

	c.informer = informer
	c.indexer = indexer
	c.nodeLister = lister_v1.NewNodeLister(indexer)

	return c
}

func (c *nodeController) Run(stopCh chan struct{}) {
	log.Info("Starting nodeController")

	go c.informer.Run(stopCh)
	go c.terminator.Run(stopCh)

	<-stopCh
	log.Info("Stopping nodeController")
}

func (c *nodeController) createNodeEvent(node *v1.Node, input *nodeControllerInput) error {
	if !c.terminator.doneNodes.Has(node.GetName()) {
	} else {
		return errorNodeInDone
	}

	if c.terminator.activeTerminations.Size() <= input.concurrentTerminations {
	} else {
		return errorActiveTerminationsLimitReached
	}

	if !c.terminator.activeTerminations.Has(node.GetName()) {
	} else {
		return errorActiveTerminationsHasNode
	}

	if annotations.CheckAnnotationsExists(node) == nil {
	} else {
		return errorNodeAnnoationsMissing
	}

	if !checkIfMaster(node) {
	} else {
		return errorNodeMaster
	}

	maintainWindow, _ := window.GetMaintenanceWindowFromAnnotations(node)
	if maintainWindow != nil {
		if maintainWindow.InMaintenanceWindow() == true {
			log.Infof("in maintainWindow starting with node %v window %v - %v :: current time %v",
				node.GetName(),
				maintainWindow.From(),
				maintainWindow.To(),
				time.Now(),
			)
			event := newTerminatorEvent(node.GetName())
			event.waitInterval = input.waitInterval
			c.terminator.events <- *event
			return nil
		}
	} else if maintainWindow == nil {
		log.Infof("maintainWindow not set starting termination of node %v", node.GetName())
		event := newTerminatorEvent(node.GetName())
		event.waitInterval = input.waitInterval
		c.terminator.events <- *event
		return nil
	}

	return fmt.Errorf("this should not have returned")
}
