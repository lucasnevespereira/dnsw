package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

const (
	// DNS uses UDP port 53. This is the standard port all DNS queries go to.
	dnsPort = 53

	// Snapshot length: how many bytes to capture per packet.
	// 1600 bytes is enough for any DNS packet (they're usually ~100-500 bytes).
	snapshotLen = 1600
)

// detectDefaultInterface finds the network interface used for internet traffic.
// On macOS: parses "route get default" output.
// On Linux: parses "ip route show default" output.
func detectDefaultInterface() string {
	switch runtime.GOOS {
	case "darwin":
		out, err := exec.Command("route", "get", "default").Output()
		if err != nil {
			return ""
		}
		for _, line := range strings.Split(string(out), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "interface:") {
				return strings.TrimSpace(strings.TrimPrefix(line, "interface:"))
			}
		}
	case "linux":
		out, err := exec.Command("ip", "route", "show", "default").Output()
		if err != nil {
			return ""
		}
		// Output looks like: "default via 192.168.1.1 dev eth0 ..."
		fields := strings.Fields(string(out))
		for i, f := range fields {
			if f == "dev" && i+1 < len(fields) {
				return fields[i+1]
			}
		}
	}
	return ""
}

// startCapture opens the network interface and returns a packet source.
// It sets promiscuous mode so we can see traffic from all devices on the LAN,
// and applies a BPF filter to only receive DNS packets (UDP port 53).
func startCapture(iface string) (*pcap.Handle, *gopacket.PacketSource, error) {
	// Open the network interface in promiscuous mode (the `true` param).
	// Promiscuous mode means we see ALL packets on the network, not just
	// packets addressed to our machine. This is how we see other devices.
	handle, err := pcap.OpenLive(iface, snapshotLen, true, pcap.BlockForever)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open interface %s: %w", iface, err)
	}

	// BPF (Berkeley Packet Filter): a kernel-level filter so we only receive
	// UDP packets on port 53 instead of processing ALL network traffic.
	if err := handle.SetBPFFilter(fmt.Sprintf("udp port %d", dnsPort)); err != nil {
		handle.Close()
		return nil, nil, fmt.Errorf("failed to set BPF filter: %w", err)
	}

	source := gopacket.NewPacketSource(handle, handle.LinkType())
	return handle, source, nil
}

// listInterfaces shows all network interfaces with their IP addresses.
func listInterfaces() {
	fmt.Printf("\n%sAvailable interfaces:%s\n\n", bold, reset)
	devices, err := pcap.FindAllDevs()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	for _, d := range devices {
		if len(d.Addresses) == 0 {
			continue
		}
		for _, addr := range d.Addresses {
			ip := addr.IP.String()
			// Skip loopback addresses (your machine talking to itself).
			if strings.HasPrefix(ip, "127.") || strings.HasPrefix(ip, "::1") {
				continue
			}
			fmt.Printf("  %s%-12s%s  %s\n", bold, d.Name, reset, ip)
		}
	}
	fmt.Printf("\nUsage:   sudo dnsw -i en0\n\n")
}
