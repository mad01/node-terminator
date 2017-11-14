package main

import (
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

	"github.com/mad01/k8s-node-terminator/pkg/annotations"
	"github.com/mad01/k8s-node-terminator/pkg/window"
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
				if !c.terminator.doneNodes.Has(node.GetName()) {
					if c.terminator.activeTerminations.Size() <= input.concurrentTerminations {
						if !c.terminator.activeTerminations.Has(node.GetName()) {
							if annotations.CheckAnnotationsExists(node) == nil && !checkIfMaster(node) {
								maintainWindow, err := window.GetMaintenanceWindowFromAnnotations(node)
								if err != nil {
									if maintainWindow.InMaintenanceWindow() == true {
										event := newTerminatorEvent(node.GetName())
										event.waitInterval = input.waitInterval
										c.terminator.events <- *event
									}
								} else {
									event := newTerminatorEvent(node.GetName())
									event.waitInterval = input.waitInterval
									c.terminator.events <- *event
								}
							}
						}
					}
				}
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
