package main

import (
	"fmt"
	"os"
	"time"

	"github.com/coreos/go-systemd/login1"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	lister_v1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

const (
	annotationRebootOK             = "node.updater/ok"             // should be true/false as string
	annotationExpectKubeletVersion = "node.updater/kubeletVersion" // string format v<major>.<minor>.<patch>
)

type agentNodeController struct {
	client     kubernetes.Interface
	informer   cache.Controller
	indexer    cache.Indexer
	nodeLister lister_v1.NodeLister
}

func newAgentNodeController(
	client kubernetes.Interface,
	namespace string,
	updateInterval time.Duration,
	listOptions metav1.ListOptions) *agentNodeController {

	ac := &agentNodeController{
		client: client,
	}

	indexer, informer := cache.NewIndexerInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (runtime.Object, error) {
				return client.Core().Nodes().List(listOptions)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return client.Core().Nodes().Watch(listOptions)

			},
		},
		// The types of objects this informer will return
		&v1.Node{},
		updateInterval,
		// Callback Functions to trigger on add/update/delete
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {},
			UpdateFunc: func(old, new interface{}) {
				if key, err := cache.MetaNamespaceKeyFunc(new); err == nil {
					log.Debugf("updateFunc key: %v", key)
					newNode := new.(*v1.Node)
					oldNode := old.(*v1.Node)
					if newNode.ResourceVersion != oldNode.ResourceVersion {
						reboot, err := rebootOK(newNode)
						if err != nil {
							log.Errorf("%v", err.Error())
						}
						if reboot {
							rebooting := rebootNode()
							log.Infof("rebooting node: %v", rebooting)
						}

						fmt.Printf("node %v time: %v\n", newNode.Name, time.Now().Unix())
					}
				}
			},
			DeleteFunc: func(obj interface{}) {},
		},
		cache.Indexers{},
	)

	ac.informer = informer
	ac.indexer = indexer
	ac.nodeLister = lister_v1.NewNodeLister(indexer)

	return ac
}

func (c *agentNodeController) Run(stopCh chan struct{}) {
	log.Info("Starting agentNodeController")

	go c.informer.Run(stopCh)

	<-stopCh
	log.Info("Stopping agentNodeController")
}

func rebootNode() bool {
	if os.Geteuid() != 0 {
		fmt.Fprintln(os.Stderr, "Must be root to initiate reboot.")
		return false
	}

	lgn, err := login1.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing login connection:", err)
		return false
	}

	lgn.Reboot(false)

	return true
}

func rebootOK(node *v1.Node) (bool, error) {
	annotations := node.GetAnnotations()
	if _, ok := annotations[annotationExpectKubeletVersion]; !ok {
		return false, fmt.Errorf("missing annotation %v", annotationExpectKubeletVersion)
	}
	if _, ok := annotations[annotationRebootOK]; !ok {
		return false, fmt.Errorf("missing annotation %v", annotationRebootOK)
	}

	if annotations[annotationExpectKubeletVersion] == node.Status.NodeInfo.KubeletVersion {
		return false, nil
	}
	if annotations[annotationRebootOK] == "true" {
		log.Info("ok to reboot node")
		return true, nil
	}
	return false, nil
}
