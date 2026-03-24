package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func proxyCmd() *cobra.Command {
	var upstream string
	var noDedupe bool
	cmd := &cobra.Command{
		Use:   "proxy",
		Short: "Run a DNS proxy to see all devices on your network",
		Long: `Start a DNS proxy server on this machine. All devices that use
this machine as their DNS server will have their queries logged.

Set your router's DNS to this machine's IP to see traffic from
phones, TVs, and other devices on your network.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			loadDeviceNames()

			proxy := &dnsProxy{
				listenAddr: ":53",
				upstream:   upstream + ":53",
				noDedupe:   noDedupe,
			}

			// Handle Ctrl+C: show shutdown instructions.
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-sigCh
				printProxyShutdown()
				os.Exit(0)
			}()

			printProxyBanner(fmt.Sprintf("%s:53", localIP()), proxy.upstream)

			if err := proxy.run(); err != nil {
				fmt.Printf("\n%sError:%s %v\n", red, reset, err)
				fmt.Printf("Try: %ssudo dnsw proxy%s\n\n", bold, reset)
				os.Exit(1)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&upstream, "upstream", "8.8.8.8", "Upstream DNS server to forward queries to")
	cmd.Flags().BoolVar(&noDedupe, "no-dedupe", false, "Show all DNS queries without merging duplicates")

	return cmd
}
