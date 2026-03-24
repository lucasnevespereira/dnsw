package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func watchCmd() *cobra.Command {
	var iface string
	var noDedupe bool

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Start watching DNS queries on your network",
		Long:  "Capture DNS queries from all devices on your local network in real time.",
		RunE: func(cmd *cobra.Command, args []string) error {
			loadDeviceNames()

			// Auto-detect the default interface if none was specified.
			if iface == "" {
				detected := detectDefaultInterface()
				if detected == "" {
					fmt.Printf("\n%sCould not detect network interface.%s\n", red, reset)
					fmt.Printf("Run %sdnsw interfaces%s to see available options.\n\n", bold, reset)
					os.Exit(1)
				}
				iface = detected
			}

			handle, source, err := startCapture(iface)
			if err != nil {
				fmt.Printf("\n%sError:%s %v\n", red, reset, err)
				fmt.Printf("Try: %ssudo dnsw -i %s%s\n\n", bold, iface, reset)
				os.Exit(1)
			}
			defer handle.Close()

			printBanner(iface, isWifiInterface(iface))

			for packet := range source.Packets() {
				handlePacket(packet, noDedupe)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&iface, "interface", "i", "", "Network interface (e.g. en0, wlan0)")
	cmd.Flags().BoolVar(&noDedupe, "no-dedupe", false, "Show all DNS queries without merging duplicates")

	return cmd
}
