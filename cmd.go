package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func cmdTerminateNode() *cobra.Command {
	var nodename string

	var command = &cobra.Command{
		Use:   "terminate",
		Short: "terminate node with private dns name",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			client := newEC2()
			err := client.awsTerminateInstance(nodename)
			if err != nil {
				log.Errorf("failed to terminate instance %v %v", nodename, err.Error())
			}
		},
	}
	command.Flags().StringVar(&nodename, "node.name", "", "name of node")
	command.MarkFlagRequired("node.name")

	return command
}

func cmdPatchNode() *cobra.Command {
	var kubeconfig, nodename string
	var unschedulable bool

	var command = &cobra.Command{
		Use:   "patch",
		Short: "patch node",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := k8sGetClient(kubeconfig)
			if err != nil {
				log.Error(fmt.Errorf("failed to get client: %v", err))
			}
			if unschedulable {
				err := setNodeUnschedulable(nodename, client)
				if err != nil {
					fmt.Println(err.Error())
				}
			} else {
				err := setNodeSchedulable(nodename, client)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		},
	}
	command.Flags().StringVar(&kubeconfig, "kube.config", "", "path to kube config")
	command.Flags().BoolVar(&unschedulable, "unschedulable", false, "set node to unschedulable")
	command.Flags().StringVar(&nodename, "node.name", "", "name of node")
	command.MarkFlagRequired("kube.config")
	command.MarkFlagRequired("node.name")

	return command
}

func cmdCordinator() *cobra.Command {
	var updateInterval time.Duration
	var kubeconfig, nodename string
	var command = &cobra.Command{
		Use:   "cordinator",
		Short: "run cordinator",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := k8sGetClient(kubeconfig)
			if err != nil {
				log.Error(fmt.Errorf("failed to get client: %v", err))
			}
			fmt.Println(client.ServerVersion())
		},
	}

	command.Flags().DurationVar(&updateInterval, "update.interval", 10*time.Second, "time.Duration cache update interval")
	command.Flags().StringVar(&kubeconfig, "kube.config", "", "path to kube config")
	command.Flags().StringVar(&nodename, "node.name", "", "name of node")
	command.MarkFlagRequired("node.name")

	return command
}

func cmdVersion() *cobra.Command {
	var command = &cobra.Command{
		Use:   "version",
		Short: "get version",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(getVersion())
		},
	}
	return command
}

func runCmd() error {
	var rootCmd = &cobra.Command{Use: "k8s-node-updater"}
	rootCmd.AddCommand(cmdCordinator())
	rootCmd.AddCommand(cmdVersion())
	rootCmd.AddCommand(cmdPatchNode())
	rootCmd.AddCommand(cmdTerminateNode())

	err := rootCmd.Execute()
	if err != nil {
		return fmt.Errorf("%v", err.Error())
	}
	return nil
}
