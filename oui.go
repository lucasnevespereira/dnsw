package main

import "strings"

// ouiTable maps the first 3 bytes of a MAC address to a vendor name.
// Built at startup from the vendor list below.
var ouiTable = map[string]string{}

func init() {
	for _, v := range ouiVendors {
		for _, prefix := range v.prefixes {
			ouiTable[prefix] = v.name
		}
	}
}

// lookupVendor returns the vendor name for a MAC address, or "" if unknown.
func lookupVendor(mac string) string {
	if len(mac) < 8 {
		return ""
	}
	prefix := strings.ToUpper(mac[:8])
	return ouiTable[prefix]
}

// Common MAC address prefixes for consumer devices.
// Covers the devices most likely to appear on a home network.
// Source: IEEE OUI registry (standards.ieee.org/products-programs/regauth)
var ouiVendors = []struct {
	name     string
	prefixes []string
}{
	// Phones, tablets, laptops
	{"Apple", []string{
		"00:1C:B3", "3C:15:C2", "AC:DE:48", "A4:83:E7",
		"F0:18:98", "14:7D:DA", "6C:96:CF", "A8:51:6B",
		"DC:A9:04", "88:66:A5", "3C:06:30", "78:7B:8A",
		"D0:81:7A", "7C:D1:C3", "A8:88:08", "64:A2:F9",
		"F4:06:69", "28:6A:BA", "60:83:73", "DC:A4:CA",
	}},
	{"Samsung", []string{
		"00:1A:8A", "58:C3:8B", "AC:36:13", "84:25:DB",
		"FC:A1:3E", "50:01:D9", "C0:BD:D1", "78:47:1D",
		"E4:7D:BD", "BC:14:85", "8C:77:12", "D0:22:BE",
	}},
	{"Xiaomi", []string{
		"28:6C:07", "64:CC:2E", "74:23:44", "7C:1D:D9",
		"8C:DE:F9", "50:64:2B", "9C:99:A0",
	}},
	{"Huawei", []string{
		"00:18:82", "00:25:68", "20:F1:7C", "48:46:FB",
		"24:09:95", "88:28:B3", "70:8A:09",
	}},
	{"OnePlus", []string{
		"94:65:2D", "C0:EE:FB",
	}},
	{"Motorola", []string{
		"00:04:56", "00:0C:E5", "AC:37:43", "E4:90:7E",
	}},

	// Computers
	{"Dell", []string{
		"00:14:22", "18:03:73", "34:17:EB", "5C:26:0A",
	}},
	{"HP", []string{
		"00:17:A4", "3C:D9:2B", "64:51:06", "94:57:A5",
	}},
	{"Lenovo", []string{
		"00:06:1B", "28:D2:44", "54:E1:AD", "98:FA:9B",
	}},
	{"Intel", []string{
		"00:13:02", "34:13:E8", "68:17:29", "80:86:F2",
		"A4:34:D9", "3C:F0:11", "48:51:B7",
	}},

	// Smart home and IoT
	{"Google", []string{
		"F4:F5:E8", "54:60:09", "A4:77:33", "30:FD:38",
		"48:D6:D5", "94:EB:2C", "3C:5A:B4",
	}},
	{"Amazon", []string{
		"44:65:0D", "FC:65:DE", "68:54:FD", "40:B4:CD",
		"74:C2:46", "A0:02:DC", "38:F7:3D",
	}},
	{"Sonos", []string{
		"00:0E:58", "34:7E:5C", "48:A6:B8", "54:2A:1B",
	}},
	{"Philips", []string{
		"00:17:88",
	}},
	{"Ring", []string{
		"34:3E:A4", "94:A1:A2",
	}},

	// TV and streaming
	{"LG", []string{
		"00:1C:62", "00:1E:75", "00:22:A9", "10:68:3F",
		"58:A2:B5", "BC:F5:AC",
	}},
	{"Roku", []string{
		"B0:A7:37", "D8:31:34", "AC:3A:7A", "DC:3A:5E",
	}},
	{"TCL", []string{
		"C0:79:82", "FC:D2:B6",
	}},

	// Gaming
	{"PlayStation", []string{
		"00:04:1F", "00:1A:80", "00:24:8D", "AC:9B:0A",
		"78:C8:81", "F8:46:1C",
	}},
	{"Xbox", []string{
		"7C:ED:8D", "00:50:F2", "28:18:78", "60:45:BD",
		"98:5F:D3",
	}},
	{"Nintendo", []string{
		"00:17:AB", "00:1F:32", "00:22:D7", "00:24:F3",
		"7C:BB:8A", "E8:4E:CE",
	}},
	{"Valve", []string{
		"1E:F2:1A",
	}},

	// Networking equipment
	{"TP-Link", []string{
		"60:32:B1", "50:C7:BF", "54:C8:0F", "98:DA:C4",
		"B0:BE:76",
	}},
	{"Netgear", []string{
		"00:14:6C", "20:E5:2A", "44:94:FC", "84:1B:5E",
	}},
	{"Asus", []string{
		"00:1A:92", "04:D4:C4", "2C:4D:54", "50:46:5D",
	}},

	// Other
	{"Raspberry Pi", []string{
		"B8:27:EB", "DC:A6:32", "D8:3A:DD", "E4:5F:01",
	}},
	{"Microsoft", []string{
		"00:50:F2", "28:18:78", "7C:1E:52",
	}},
}
