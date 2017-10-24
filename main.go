package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Agent struct{}
type Controller struct{}
type FlagVars struct {
	nodeName   string
	kubeconfig string
	subcommand string
	agent      *Agent
	controller *Controller
	interval   time.Duration
}

func flags() *FlagVars {
	c := FlagVars{
		agent:      new(Agent),
		controller: new(Controller),
	}
	agent := flag.NewFlagSet("agent", flag.ExitOnError)
	agent.StringVar(&c.nodeName, "node.name", "", "name of node")
	agent.StringVar(&c.kubeconfig, "kube.config", "", "kubernetes config")
	agent.DurationVar(&c.interval, "update.interval", 10*time.Second, "time.Duration cache update intermal")

	controller := flag.NewFlagSet("controller", flag.ExitOnError)
	controller.StringVar(&c.nodeName, "node.name", "", "name of node")
	controller.StringVar(&c.kubeconfig, "kube.config", "", "kubernetes config")
	controller.DurationVar(&c.interval, "update.interval", 10*time.Second, "time.Duration cache update intermal")

	version := flag.NewFlagSet("version", flag.ExitOnError)
	HelpMsg := `k8s node updater.

Commands:
  version       Get version
  agent         Run agent
  controller    Run controller`

	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 {
		fmt.Println(HelpMsg)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "version":
		c.subcommand = "version"
		version.Parse(os.Args[2:])
	case "agent":
		c.subcommand = "agent"
		agent.Parse(os.Args[2:])
	case "controller":
		c.subcommand = "controller"
		controller.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if version.Parsed() {
		fmt.Println(getVersion())
	}
	if agent.Parsed() {
		fmt.Println("agent Parsed")
	}
	if controller.Parsed() {
		fmt.Println("controller Parsed")
	}
	return &c
}

func main() {
	LogInit(false)
	config := flags()

	client, err := k8sGetClient(config.kubeconfig)
	if err != nil {
		log.Error(fmt.Errorf("failed to get client: %v", err))
	}

	stopCh := make(chan struct{})
	defer close(stopCh)

	switch config.subcommand {
	case "agent":
		fmt.Println("agent subcommand")
		selector := fmt.Sprintf("kubernetes.io/hostname=%v", config.nodeName)
		agentController := newAgentNodeController(
			client,
			metav1.NamespaceAll,
			config.interval,
			metav1.ListOptions{LabelSelector: selector},
		)
		agentController.Run(stopCh)

	case "controller":
		fmt.Println("controller subcommand")
	}

}
