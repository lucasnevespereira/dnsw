package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	// Max size of a DNS packet over UDP (RFC 1035).
	maxDNSPacketSize = 512

	// How long to wait for an upstream DNS response.
	upstreamTimeout = 5 * time.Second
)

// dnsProxy receives DNS queries from devices on the network,
// logs them, and forwards them to an upstream DNS server.
type dnsProxy struct {
	listenAddr string
	upstream   string
	noDedupe   bool
}

// run starts the DNS proxy. Blocks until the listener is closed.
func (p *dnsProxy) run() error {
	addr, err := net.ResolveUDPAddr("udp", p.listenAddr)
	if err != nil {
		return fmt.Errorf("invalid listen address: %w", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("could not start DNS proxy on %s: %w", p.listenAddr, err)
	}
	defer conn.Close()

	buf := make([]byte, maxDNSPacketSize)

	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			// Connection closed (shutdown).
			if strings.Contains(err.Error(), "use of closed network connection") {
				return nil
			}
			continue
		}

		// Copy the packet so we can process it concurrently.
		packet := make([]byte, n)
		copy(packet, buf[:n])

		go p.handleQuery(conn, clientAddr, packet)
	}
}

// handleQuery logs a DNS query, forwards it upstream, and sends the response back.
func (p *dnsProxy) handleQuery(conn *net.UDPConn, clientAddr *net.UDPAddr, packet []byte) {
	if len(packet) < dnsHeaderSize {
		return
	}

	// Parse the query for logging.
	domain, _ := parseDNSQuery(packet)

	// Identify the device by its source IP.
	srcIP := clientAddr.IP.String()
	device, isNew := resolveDevice(srcIP)

	if isNew {
		printNewDevice(device, srcIP, "")
	}

	// Log the query (skip noise and dedup).
	if domain != "" && !strings.HasSuffix(domain, ".local") && !strings.HasSuffix(domain, ".arpa") {
		if p.noDedupe || !isDuplicate(device, domain) {
			cat := categorize(domain)
			now := time.Now().Format("15:04:05")
			printQuery(now, device, cat, domain)
		}
	}

	// Forward the raw packet to the upstream DNS server and relay the response.
	response, err := p.forward(packet)
	if err != nil {
		return
	}

	conn.WriteToUDP(response, clientAddr)
}

// forward sends a DNS packet to the upstream server and returns the response.
func (p *dnsProxy) forward(packet []byte) ([]byte, error) {
	upstream, err := net.ResolveUDPAddr("udp", p.upstream)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, upstream)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(upstreamTimeout))

	if _, err := conn.Write(packet); err != nil {
		return nil, err
	}

	buf := make([]byte, maxDNSPacketSize)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf[:n], nil
}

// localIP returns this machine's IP on the local network.
func localIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return "your-mac-ip"
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}

// printProxyBanner shows the proxy status and setup instructions.
func printProxyBanner(listenAddr, upstream string) {
	ip := localIP()

	fmt.Printf("\n%s%s DNS WATCHER (proxy mode)%s\n", bold, cyan, reset)
	fmt.Printf("%s────────────────────────────────────────────────────────%s\n", gray, reset)
	fmt.Printf("  DNS proxy  : %s%s%s\n", bold, listenAddr, reset)
	fmt.Printf("  Upstream   : %s\n", upstream)
	fmt.Printf("  Press      : %sCtrl+C%s to stop\n", bold, reset)
	fmt.Printf("%s────────────────────────────────────────────────────────%s\n\n", gray, reset)

	fmt.Printf("  %s%sTo see ALL devices on your network:%s\n\n", bold, yellow, reset)
	fmt.Printf("    1. Open your router admin page (usually %shttp://192.168.1.1%s)\n", bold, reset)
	fmt.Printf("    2. Find DHCP or LAN settings\n")
	fmt.Printf("    3. Set primary DNS to:   %s%s%s\n", bold, ip, reset)
	fmt.Printf("    4. Set secondary DNS to: %s8.8.8.8%s  (fallback if proxy stops)\n", bold, reset)
	fmt.Printf("    5. Save and wait ~1 min for devices to pick up the change\n\n")

	fmt.Printf("%s%-10s  %-20s  %-12s  %s%s\n", gray, "TIME", "DEVICE", "CATEGORY", "DOMAIN", reset)
	fmt.Printf("%s%s%s\n", gray, strings.Repeat("─", 72), reset)
}

// printProxyShutdown shows instructions when the proxy stops.
var printProxyShutdownOnce sync.Once

func printProxyShutdown() {
	printProxyShutdownOnce.Do(func() {
		fmt.Printf("\n\n%s────────────────────────────────────────────────────────%s\n", gray, reset)
		fmt.Printf("  DNS proxy stopped.\n\n")
		fmt.Printf("  Your devices will automatically use the fallback DNS (8.8.8.8)\n")
		fmt.Printf("  so internet should keep working. If not, set your router DNS\n")
		fmt.Printf("  back to automatic.\n")
		fmt.Printf("%s────────────────────────────────────────────────────────%s\n\n", gray, reset)
	})
}
