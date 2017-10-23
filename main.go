package main

import (
	"flag"
	"fmt"
	"os"
)

type Agent struct {
	nodeName string
}
type Controller struct {
	nodeName string
}
type FlagVars struct {
	subcommand string
	agent      *Agent
	controller *Controller
}

func flags() *FlagVars {
	c := FlagVars{
		agent:      new(Agent),
		controller: new(Controller),
	}
	agent := flag.NewFlagSet("agent", flag.ExitOnError)
	agent.StringVar(&c.agent.nodeName, "node.name", "", "name of node")

	controller := flag.NewFlagSet("controller", flag.ExitOnError)
	controller.StringVar(&c.controller.nodeName, "node.name", "", "name of node")
	version := flag.NewFlagSet("version", flag.ExitOnError)

	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 {
		msg := `k8s node updater.

Commands:
  version       Get version
  agent         Run agent
  controller    Run controller`
		fmt.Println(msg)
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
	config := flags()
	switch config.subcommand {
	case "agent":
		fmt.Println("agent subcommand")
	case "controller":
		fmt.Println("controller subcommand")
	}

}
