package main

import (
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	lister_v1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

// manages all nodes and sets annotations on N nodes to reboot at one time
type nodeController struct {
	client     kubernetes.Interface
	informer   cache.Controller
	indexer    cache.Indexer
	nodeLister lister_v1.NodeLister
}

func newNodeController(
	client kubernetes.Interface,
	namespace string,
	updateInterval time.Duration) *nodeController {

	c := &nodeController{
		client: client,
	}

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

	<-stopCh
	log.Info("Stopping nodeController")
}
