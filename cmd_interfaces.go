package main

import "github.com/spf13/cobra"

func interfacesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "interfaces",
		Short: "List available network interfaces",
		Run: func(cmd *cobra.Command, args []string) {
			listInterfaces()
		},
	}
}
