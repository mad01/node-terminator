package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func cmdCordinator() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "cordinator",
		Short: "run cordinator",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			var updateInterval time.Duration
			var kubeconfig, nodename string
			cmd.Flags().DurationVar(&updateInterval, "update.interval", 10*time.Second, "time.Duration cache update interval")
			cmd.Flags().StringVar(&kubeconfig, "kube.config", "", "path to kube config")
			cmd.Flags().StringVar(&nodename, "node.name", "", "name of node")
			cmd.MarkFlagRequired("node.name")

			// do stuff call external call or not ?

			client, err := k8sGetClient(kubeconfig)
			if err != nil {
				log.Error(fmt.Errorf("failed to get client: %v", err))
			}
			fmt.Println(client.ServerVersion())
		},
	}
	return cmd
}

func cmdVersion() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "get version",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(getVersion())
		},
	}
	return cmd
}

func runCmd() error {
	var rootCmd = &cobra.Command{
		Use:   "k8s-node-updater",
		Short: "manages k8s node upgrade",
		Long:  "",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	rootCmd.AddCommand(cmdCordinator())
	rootCmd.AddCommand(cmdVersion())

	err := rootCmd.Execute()
	if err != nil {
		return fmt.Errorf("%v", err.Error())
	}
	return nil
}
