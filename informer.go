package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	lister_v1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

type nodeControllerInput struct {
	waitInterval           time.Duration // TODO: need better name
	updateInterval         time.Duration
	kubeconfig             string
	concurrentTerminations int
}

// manages all nodes and sets annotations on N nodes to reboot at one time
type nodeController struct {
	client                 *kubernetes.Clientset
	informer               cache.Controller
	indexer                cache.Indexer
	nodeLister             lister_v1.NodeLister
	terminator             *Terminator
	concurrentTerminations int
}

func newNodeController(input *nodeControllerInput) *nodeController {
	client, err := k8sGetClient(input.kubeconfig)
	if err != nil {
		panic(fmt.Sprintf("failed to get client: %v", err.Error()))
	}

	c := &nodeController{
		client:                 client,
		concurrentTerminations: input.concurrentTerminations,
	}

	c.terminator = newTerminator(input.kubeconfig)
	c.terminator.eviction.waitInterval = input.waitInterval

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
		updateInterval,
		// Callback Functions to trigger on add/update/delete
		cache.ResourceEventHandlerFuncs{
			AddFunc:    func(obj interface{}) {},
			UpdateFunc: func(old, new interface{}) {},
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

	// start up your worker threads based on concurrentTerminations
	for i := 0; i < c.concurrentTerminations; i++ {
		go wait.Until(c.terminator.Run, time.Second, stopCh)
	}

	<-stopCh
	log.Info("Stopping nodeController")
}
