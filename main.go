package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	iface := flag.String("i", "", "Network interface (e.g. en0, wlan0)")
	list := flag.Bool("list", false, "List available network interfaces")
	noDedupe := flag.Bool("no-dedupe", false, "Show all DNS queries (don't merge duplicates)")
	flag.Parse()

	if *list {
		listInterfaces()
		os.Exit(0)
	}

	// Auto-detect the default interface if none was specified.
	if *iface == "" {
		detected := detectDefaultInterface()
		if detected == "" {
			fmt.Printf("\n%sCould not detect network interface.%s\n", red, reset)
			fmt.Printf("Use %s-i <interface>%s to specify one, or %s--list%s to see available interfaces.\n\n", bold, reset, bold, reset)
			os.Exit(1)
		}
		*iface = detected
	}

	handle, source, err := startCapture(*iface)
	if err != nil {
		fmt.Printf("\n%sError:%s %v\n", red, reset, err)
		fmt.Printf("Try: %ssudo dnsw -i %s%s\n\n", bold, *iface, reset)
		os.Exit(1)
	}
	defer handle.Close()

	printBanner(*iface)

	for packet := range source.Packets() {
		handlePacket(packet, *noDedupe)
	}
}
