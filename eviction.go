package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mad01/k8s-node-terminator/pkg/kutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubernetes/pkg/kubectl/cmd"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

// Upgrade nodes, with greater parallelism
// We run nodes in series, even if they are in separate instance groups
// typically they will not being separate instance groups. If you roll the nodes in parallel
// you can get into a scenario where you can evict multiple statefulset pods from the same
// statefulset at the same time. Further improvements needs to be made to protect from this as
// well.

func drainNode(nodename string, client *kubernetes.Clientset) error {
	// look at eviction policy
	// client.Core().Pods().Evict
	return nil
}

func newEviction(kubeconfig string) (*Eviction, error) {
	config, err := k8sGetClientConfig(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get kube rest config: %v", err.Error())
	}

	e := Eviction{
		DrainInterval: 5 * time.Minute,
		ClientConfig:  kutil.NewClientConfig(config, metav1.NamespaceAll),
	}
	return &e, nil
}

// Eviction struct
type Eviction struct {
	DrainInterval time.Duration
	ClientConfig  clientcmd.ClientConfig
}

// DrainNode drains a K8s node.
func (e *Eviction) DrainNode(nodename string) error {

	f := cmdutil.NewFactory(e.ClientConfig)

	// TODO: Send out somewhere else, also DrainOptions has errout
	out := os.Stdout
	errOut := os.Stderr

	options := &cmd.DrainOptions{
		Factory:          f,
		Out:              out,
		IgnoreDaemonsets: true,
		Force:            true,
		DeleteLocalData:  true,
		ErrOut:           errOut,
	}

	cmd := &cobra.Command{
		Use: "cordon NODE",
	}
	args := []string{nodename}
	err := options.SetupDrain(cmd, args)
	if err != nil {
		return fmt.Errorf("error setting up drain: %v", err)
	}

	err = options.RunCordonOrUncordon(true)
	if err != nil {
		return fmt.Errorf("error cordoning node node: %v", err)
	}

	err = options.RunDrain()
	if err != nil {
		return fmt.Errorf("error draining node: %v", err)
	}

	if e.DrainInterval > time.Second*0 {
		log.Infof("Waiting for %s for pods to stabilize after draining.", e.DrainInterval)
		time.Sleep(e.DrainInterval)
	}

	return nil
}
