package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func configCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Show config file location and contents",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := configPath()
			if err != nil {
				return fmt.Errorf("could not find config directory: %w", err)
			}

			fmt.Printf("\n  %sConfig directory%s   %s\n", bold, reset, filepath.Dir(path))
			fmt.Printf("  %sDevices file%s       %s\n\n", bold, reset, path)

			data, err := os.ReadFile(path)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Printf("  No config file yet.\n")
					fmt.Printf("  Run %sdnsw devices set <ip> <name>%s to create one.\n\n", bold, reset)
					return nil
				}
				return err
			}

			fmt.Printf("  %sdevices.json:%s\n\n", gray, reset)
			fmt.Printf("  %s\n", string(data))

			return nil
		},
	}
}
