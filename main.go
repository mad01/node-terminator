package main

import (
	"flag"
	"fmt"
	"os"
)

type Agent struct{}
type Controller struct{}

func flags() {
	agent := flag.NewFlagSet("agent", flag.ExitOnError)
	controller := flag.NewFlagSet("controller", flag.ExitOnError)
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
		version.Parse(os.Args[2:])
	case "agent":
		agent.Parse(os.Args[2:])
	case "controller":
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

}

func main() {
	flags()
}
