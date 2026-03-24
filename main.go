package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	root := &cobra.Command{
		Use:     "dnsw",
		Short:   "Real-time DNS watcher for your local network",
		Version: version,
	}

	root.AddCommand(
		watchCmd(),
		devicesCmd(),
		configCmd(),
		interfacesCmd(),
	)

	// Running `dnsw` with no subcommand starts watching.
	root.RunE = watchCmd().RunE
	root.Flags().AddFlagSet(watchCmd().Flags())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
