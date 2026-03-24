package main

import (
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

const (
	// A DNS packet header is always 12 bytes. It contains:
	// - Transaction ID (2 bytes)
	// - Flags (2 bytes), tells us if it's a query or response
	// - Question/Answer counts (8 bytes)
	dnsHeaderSize = 12

	// In DNS, domain names can be "compressed" using pointers.
	// A pointer byte starts with the two high bits set (11xxxxxx = 0xC0).
	// This tells the parser "the rest of this name is at another offset."
	dnsPointerMask = 0xC0

	// The QR (Query/Response) bit is the highest bit in the flags field.
	// 0 = query (someone asking "what is facebook.com?")
	// 1 = response (the answer coming back)
	// We shift right by 15 to isolate this single bit.
	dnsQRBitShift = 15

	// Byte offsets within the DNS header where the flags field starts and ends.
	dnsFlagsStart = 2
	dnsFlagsEnd   = 4

	// After the domain name, there are 4 bytes: query type (2) + query class (2).
	dnsQueryFieldsSize = 4

	// DNS pointer is 2 bytes: the pointer marker byte + the offset byte.
	dnsPointerSize = 2
)

// DNS query types , each one asks for a different kind of record.
// A = IPv4 address, AAAA = IPv6 address, HTTPS = modern secure connection info, etc.
var queryTypes = map[uint16]string{
	1: "A", 2: "NS", 5: "CNAME", 6: "SOA", 12: "PTR", 15: "MX",
	16: "TXT", 28: "AAAA", 33: "SRV", 65: "HTTPS", 255: "ANY",
}

// parseDNSQuery extracts the domain name and query type from raw DNS packet bytes.
// DNS encodes domain names as length-prefixed labels:
//
//	"\x07example\x03com\x00" becomes "example.com"
//
// Each segment starts with its length, and the name ends with a zero byte.
func parseDNSQuery(data []byte) (domain string, queryType uint16) {
	if len(data) < dnsHeaderSize {
		return "", 0
	}

	// Skip the 12-byte header to get to the question section.
	offset := dnsHeaderSize
	var labels []string

	for offset < len(data) {
		labelLen := int(data[offset])

		// Zero byte means end of the domain name.
		if labelLen == 0 {
			offset++
			break
		}

		// If the two high bits are set, this is a compression pointer
		// (the name continues at a different offset in the packet).
		if labelLen&dnsPointerMask == dnsPointerMask {
			offset += dnsPointerSize
			break
		}

		offset++ // skip the length byte itself
		if offset+labelLen > len(data) {
			return "", 0
		}
		labels = append(labels, string(data[offset:offset+labelLen]))
		offset += labelLen
	}

	// After the domain name: 2 bytes for query type + 2 bytes for query class.
	if offset+dnsQueryFieldsSize > len(data) {
		return strings.Join(labels, "."), 0
	}
	queryType = binary.BigEndian.Uint16(data[offset : offset+2])
	return strings.Join(labels, "."), queryType
}

// handlePacket processes a single captured network packet.
// It extracts the source IP, checks it's a DNS query, parses it, and prints it.
func handlePacket(packet gopacket.Packet, noDedupe bool) {
	var srcIP string

	// A packet can be IPv4 or IPv6, we need to handle both.
	if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		srcIP = ip.SrcIP.String()
	} else if ipLayer := packet.Layer(layers.LayerTypeIPv6); ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv6)
		srcIP = ip.SrcIP.String()
	} else {
		return
	}

	// Extract the UDP layer. DNS runs over UDP (User Datagram Protocol).
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		return
	}
	udp, _ := udpLayer.(*layers.UDP)

	// Only look at outgoing queries (destination = DNS port 53).
	// Responses come FROM port 53, we don't care about those.
	if udp.DstPort != dnsPort {
		return
	}

	payload := udp.Payload
	if len(payload) < dnsHeaderSize {
		return
	}

	// Check the QR bit in the flags field:
	// 0 = query (we want these), 1 = response (skip).
	flags := binary.BigEndian.Uint16(payload[dnsFlagsStart:dnsFlagsEnd])
	isResponse := (flags>>dnsQRBitShift)&1 != 0
	if isResponse {
		return
	}

	domain, qt := parseDNSQuery(payload)

	// Filter out internal/noise domains:
	// .local = mDNS (devices discovering each other on LAN)
	// .arpa  = reverse DNS lookups (IP to hostname)
	if domain == "" || strings.HasSuffix(domain, ".local") || strings.HasSuffix(domain, ".arpa") {
		return
	}

	qtStr, ok := queryTypes[qt]
	if !ok {
		qtStr = fmt.Sprintf("%d", qt)
	}

	device := resolveDevice(srcIP)

	// Skip duplicate queries from the same device within a short window.
	if !noDedupe && isDuplicate(device, domain) {
		return
	}

	cat := categorize(domain)
	now := time.Now().Format("15:04:05")

	printQuery(now, device, qtStr, cat, domain)
}
