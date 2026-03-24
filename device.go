package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// deviceNames holds user-defined names from the config file.
// Maps IP addresses to friendly names like "Dad's iPhone".
var deviceNames = map[string]string{}

// deviceCache stores resolved names so we only look up once per IP.
var (
	deviceCache = map[string]string{}
	deviceMu    sync.RWMutex
)

// macForIP records which MAC address belongs to which IP.
// Populated from the Ethernet layer of captured packets.
var (
	macForIP   = map[string]string{}
	macForIPMu sync.Mutex
)

// autoNames stores generated names per MAC so the same device
// always gets the same name within a session (e.g. "Apple-1").
var (
	autoNames   = map[string]string{}
	autoNamesMu sync.Mutex
)

// vendorCounter tracks how many devices of each vendor we've seen.
var vendorCounter = map[string]int{}

// loadDeviceNames reads user-defined device names from ~/.config/dnsw/devices.json.
// If the file doesn't exist, it creates a sample one so users know about it.
func loadDeviceNames() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return
	}

	dir := filepath.Join(configDir, "dnsw")
	path := filepath.Join(dir, "devices.json")

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			createSampleDevicesFile(dir, path)
		}
		return
	}

	if err := json.Unmarshal(data, &deviceNames); err != nil {
		fmt.Printf("%sWarning:%s could not parse %s: %v\n", yellow, reset, path, err)
	}
}

func createSampleDevicesFile(dir, path string) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return
	}

	sample := map[string]string{
		"192.168.1.32": "My MacBook",
		"192.168.1.45": "iPhone",
		"192.168.1.50": "Living Room TV",
	}

	data, err := json.MarshalIndent(sample, "", "  ")
	if err != nil {
		return
	}

	os.WriteFile(path, append(data, '\n'), 0644)
}

// registerMAC records the MAC address for a source IP.
// Called once per packet from the Ethernet layer.
func registerMAC(ip, mac string) {
	macForIPMu.Lock()
	if _, exists := macForIP[ip]; !exists {
		macForIP[ip] = mac
	}
	macForIPMu.Unlock()
}

// autoNameFromMAC generates a name like "Apple-1" from a MAC address.
// Returns "" if the vendor is unknown. Stable within a session:
// the same MAC always returns the same name.
func autoNameFromMAC(mac string) string {
	autoNamesMu.Lock()
	defer autoNamesMu.Unlock()

	if name, ok := autoNames[mac]; ok {
		return name
	}

	vendor := lookupVendor(mac)
	if vendor == "" {
		return ""
	}

	vendorCounter[vendor]++
	name := fmt.Sprintf("%s-%d", vendor, vendorCounter[vendor])
	autoNames[mac] = name
	return name
}

// resolveDevice turns an IP address into a human-readable device name.
// Resolution order:
//  1. User-defined names from devices.json
//  2. Reverse DNS (PTR lookup)
//  3. OUI vendor auto-name from MAC address (e.g. "Apple-1")
//  4. Raw IP address
//
// Returns the name and whether this is a newly seen device.
func resolveDevice(ip string) (string, bool) {
	// 1. User-defined name always wins.
	if name, ok := deviceNames[ip]; ok {
		deviceMu.RLock()
		_, seen := deviceCache[ip]
		deviceMu.RUnlock()
		if !seen {
			deviceMu.Lock()
			deviceCache[ip] = name
			deviceMu.Unlock()
			return name, true
		}
		return name, false
	}

	// 2. Already resolved in a previous packet.
	deviceMu.RLock()
	if name, ok := deviceCache[ip]; ok {
		deviceMu.RUnlock()
		return name, false
	}
	deviceMu.RUnlock()

	// 3. Reverse DNS: ask "what hostname is at this IP?"
	names, err := net.LookupAddr(ip)
	if err == nil && len(names) > 0 {
		// Result comes as "my-device.lan.", trim the trailing dot
		// and take just the first part (the hostname).
		name := strings.TrimSuffix(names[0], ".")
		parts := strings.Split(name, ".")
		deviceMu.Lock()
		deviceCache[ip] = parts[0]
		deviceMu.Unlock()
		return parts[0], true
	}

	// 4. Auto-name from MAC vendor (e.g. "Apple-1").
	macForIPMu.Lock()
	mac := macForIP[ip]
	macForIPMu.Unlock()

	if mac != "" {
		if name := autoNameFromMAC(mac); name != "" {
			deviceMu.Lock()
			deviceCache[ip] = name
			deviceMu.Unlock()
			return name, true
		}
	}

	// 5. Fallback to the raw IP.
	deviceMu.Lock()
	deviceCache[ip] = ip
	deviceMu.Unlock()
	return ip, true
}
