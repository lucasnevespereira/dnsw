package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
)

func devicesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "devices",
		Short: "List or rename your devices",
	}

	cmd.AddCommand(devicesListCmd())
	cmd.AddCommand(devicesSetCmd())
	cmd.AddCommand(devicesRemoveCmd())

	return cmd
}

func devicesListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Show all named devices",
		RunE: func(cmd *cobra.Command, args []string) error {
			loadDeviceNames()

			if len(deviceNames) == 0 {
				fmt.Printf("\n  No devices named yet.\n")
				fmt.Printf("  Run %sdnsw devices set <ip> <name>%s to add one.\n\n", bold, reset)
				return nil
			}

			// Sort IPs for consistent output.
			ips := make([]string, 0, len(deviceNames))
			for ip := range deviceNames {
				ips = append(ips, ip)
			}
			sort.Strings(ips)

			fmt.Printf("\n  %s%-18s  %s%s\n", gray, "IP", "NAME", reset)
			fmt.Printf("  %s%s%s\n", gray, "──────────────────────────────────────", reset)
			for _, ip := range ips {
				fmt.Printf("  %-18s  %s%s%s\n", ip, bold, deviceNames[ip], reset)
			}
			fmt.Println()

			return nil
		},
	}
}

func devicesSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "set <ip> <name>",
		Short:   "Name a device by its IP address",
		Example: "  dnsw devices set 192.168.1.45 \"Mom's iPhone\"",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ip, name := args[0], args[1]

			loadDeviceNames()
			deviceNames[ip] = name

			if err := saveDeviceNames(); err != nil {
				return fmt.Errorf("could not save: %w", err)
			}

			fmt.Printf("\n  %s✓%s  %s%s%s is now %s%s%s\n\n", green, reset, gray, ip, reset, bold, name, reset)
			return nil
		},
	}
}

func devicesRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "remove <ip>",
		Short:   "Remove a device name",
		Example: "  dnsw devices remove 192.168.1.45",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ip := args[0]

			loadDeviceNames()

			if _, ok := deviceNames[ip]; !ok {
				fmt.Printf("\n  No device named for %s.\n\n", ip)
				return nil
			}

			delete(deviceNames, ip)

			if err := saveDeviceNames(); err != nil {
				return fmt.Errorf("could not save: %w", err)
			}

			fmt.Printf("\n  %s✓%s  Removed name for %s\n\n", green, reset, ip)
			return nil
		},
	}
}

// saveDeviceNames writes the current deviceNames map to the config file.
func saveDeviceNames() error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(deviceNames, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, append(data, '\n'), 0644)
}
