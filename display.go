package main

import (
	"fmt"
	"strings"
)

// ANSI escape codes for terminal colors.
// These make the output colorful and scannable.
const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"
	gray    = "\033[90m"
)

const (
	// Max display widths to keep the table aligned.
	maxDomainDisplay = 45
	maxDeviceDisplay = 18
)

func printBanner(iface string) {
	fmt.Printf("\n%s%s DNS WATCHER%s\n", bold, cyan, reset)
	fmt.Printf("%s────────────────────────────────────────────────────%s\n", gray, reset)
	fmt.Printf("  Interface : %s%s%s\n", bold, iface, reset)
	fmt.Printf("  Mode      : promiscuous, capturing ALL devices on LAN\n")
	fmt.Printf("  Press     : %sCtrl+C%s to stop\n", bold, reset)
	fmt.Printf("%s────────────────────────────────────────────────────%s\n\n", gray, reset)
	fmt.Printf("%s%-10s  %-20s  %-6s  %-12s  %s%s\n", gray, "TIME", "DEVICE", "TYPE", "CATEGORY", "DOMAIN", reset)
	fmt.Printf("%s%s%s\n", gray, strings.Repeat("─", 78), reset)
}

func printQuery(timestamp, device, queryType string, cat category, domain string) {
	displayDomain := domain
	if len(displayDomain) > maxDomainDisplay {
		displayDomain = displayDomain[:maxDomainDisplay-3] + "..."
	}
	displayDevice := device
	if len(displayDevice) > maxDeviceDisplay {
		displayDevice = displayDevice[:maxDeviceDisplay-3] + "..."
	}

	fmt.Printf("%s%-10s%s  %-20s  %s%-6s%s  %s%s%s  %s%s%s\n",
		gray, timestamp, reset,
		displayDevice,
		bold, queryType, reset,
		cat.color, cat.label, reset,
		white, displayDomain, reset,
	)
}
