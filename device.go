package main

import (
	"net"
	"strings"
	"sync"
)

// deviceCache maps IP addresses to friendly device names.
// We cache results so we only do the lookup once per IP.
var (
	deviceCache = map[string]string{}
	deviceMu    sync.RWMutex
)

// resolveDevice turns an IP address into a human-readable device name
// using reverse DNS (PTR lookup). If no name is found, the raw IP is returned.
func resolveDevice(ip string) string {
	deviceMu.RLock()
	if name, ok := deviceCache[ip]; ok {
		deviceMu.RUnlock()
		return name
	}
	deviceMu.RUnlock()

	// Reverse DNS: ask "what hostname is at this IP?"
	names, err := net.LookupAddr(ip)
	if err == nil && len(names) > 0 {
		// Result comes as "my-device.lan.", trim the trailing dot
		// and take just the first part (the hostname).
		name := strings.TrimSuffix(names[0], ".")
		parts := strings.Split(name, ".")
		deviceMu.Lock()
		deviceCache[ip] = parts[0]
		deviceMu.Unlock()
		return parts[0]
	}

	// No hostname found, fall back to showing the IP address.
	deviceMu.Lock()
	deviceCache[ip] = ip
	deviceMu.Unlock()
	return ip
}
