# Network Concepts

A beginner-friendly guide to the networking concepts used in this project.

## What is DNS?

DNS stands for **Domain Name System**. It's like a phone book for the internet.

When you type `facebook.com` in your browser, your computer doesn't know where Facebook's server is. It needs an IP address (like `157.240.1.35`) to connect. So it sends a **DNS query** to a DNS server asking: "What is the IP address of facebook.com?"

The DNS server responds with the IP address, and your browser can then connect.

This all happens in milliseconds, before the page even starts loading.

## How DNS queries travel

1. You type `youtube.com` in your browser
2. Your computer sends a small UDP packet to port 53 of your DNS server (usually your router)
3. The DNS server looks up the domain and sends back the IP address
4. Your browser connects to that IP address

`dnsw` sits in the middle and watches step 2: it sees the DNS query packets going out from devices on your network.

## What is a network interface?

A network interface is how your computer connects to a network. Each one has a name:

| Name | What it is |
|------|------------|
| `en0` | Wi-Fi adapter on macOS |
| `en1` | Ethernet or secondary adapter on macOS |
| `eth0` | Ethernet on Linux |
| `wlan0` | Wi-Fi on Linux |
| `lo0` / `lo` | Loopback (your machine talking to itself) |

When you run `dnsw -i en0`, you're telling it to listen on your Wi-Fi adapter.

## What is promiscuous mode?

Normally, your network interface only picks up packets meant for your machine. **Promiscuous mode** tells it to capture **all** packets on the network, including those from other devices (phones, TVs, other computers).

This is why `dnsw` needs `sudo`: promiscuous mode requires elevated privileges.

> Note: on modern switched networks and Wi-Fi, promiscuous mode may not capture traffic from every device. It depends on your router and network setup.

## What is UDP?

**UDP** (User Datagram Protocol) is one of the two main ways to send data over a network (the other being TCP).

UDP is simpler and faster than TCP: it just sends the data without checking if it arrived. DNS uses UDP because queries are small and speed matters. If a DNS query gets lost, the computer just sends another one.

**Port 53** is the standard port for DNS traffic. When `dnsw` filters for "udp port 53", it's saying: "only show me DNS traffic, ignore everything else."

## What is a BPF filter?

**BPF** stands for Berkeley Packet Filter. It's a way to tell the operating system kernel: "I only care about specific packets, don't bother sending me the rest."

Without a BPF filter, `dnsw` would receive every single packet on the network (web browsing, video streaming, file downloads...) and would have to check each one to see if it's DNS. That would waste a lot of CPU.

With the filter `udp port 53`, the kernel itself drops non-DNS packets before they even reach our program.

## How DNS packets are structured

A DNS query packet looks like this:

```
+------------------+
| Header (12 bytes)|  Transaction ID, flags, counts
+------------------+
| Question         |  The domain name being queried
+------------------+
```

### The header

The first 12 bytes always contain:
- **Transaction ID** (2 bytes): a random number to match queries with responses
- **Flags** (2 bytes): tells us if this is a query or a response
- **Counts** (8 bytes): how many questions/answers are in the packet

The most important flag is the **QR bit** (Query/Response). It's the very first bit of the flags field:
- `0` = this is a question ("what is the IP of facebook.com?")
- `1` = this is an answer ("facebook.com is at 157.240.1.35")

`dnsw` only shows queries (QR = 0), not responses.

### Domain name encoding

DNS doesn't store domain names as plain text. Instead, it uses **length-prefixed labels**:

```
facebook.com  is stored as:  [8]facebook[3]com[0]

  8 = length of "facebook"
  3 = length of "com"
  0 = end of name
```

Each segment starts with a byte indicating its length, and the whole name ends with a zero byte.

DNS also supports **compression pointers**: if the same domain name appears multiple times in a packet, later occurrences can point back to the first one instead of repeating it. A pointer is identified by its first byte having the two high bits set (`11xxxxxx` in binary, or `0xC0` in hex).

## DNS query types

When your computer queries a domain, it specifies what kind of information it wants:

| Type | Name | What it asks for |
|------|------|------------------|
| A | Address | IPv4 address (e.g. `93.184.216.34`) |
| AAAA | IPv6 Address | IPv6 address (e.g. `2606:2800:220:1:...`) |
| CNAME | Canonical Name | "This domain is an alias for that domain" |
| MX | Mail Exchange | Where to send email for this domain |
| TXT | Text | Arbitrary text data (often used for verification) |
| NS | Name Server | Which DNS server is responsible for this domain |
| SRV | Service | Location of a specific service |
| HTTPS | HTTPS binding | Modern record for secure connection info |
| PTR | Pointer | Reverse lookup: IP address to hostname |

Most web browsing generates **A** and **AAAA** queries. You'll often see both at the same time because the browser tries IPv4 and IPv6 simultaneously.

## Reverse DNS (PTR lookups)

Normal DNS converts a domain name to an IP: `facebook.com` to `157.240.1.35`.

**Reverse DNS** does the opposite: it converts an IP address to a hostname. This is how `dnsw` shows device names instead of raw IPs.

Your router usually assigns hostnames to devices on the network (like `my-iphone` or `living-room-tv`). When `dnsw` sees a DNS query from `192.168.1.42`, it does a reverse DNS lookup to find out the device's name.

## Why some domains don't show up (DNS-over-HTTPS)

Traditional DNS sends queries as plain text over UDP port 53. Anyone on the network can see them (that's exactly what `dnsw` does).

**DNS-over-HTTPS (DoH)** is a newer approach where DNS queries are encrypted and sent inside regular HTTPS traffic on port 443. This makes them invisible to tools like `dnsw` because they look like normal web traffic.

Most modern browsers enable DoH by default. This is why you might visit `facebook.com` but not see it in `dnsw`. The browser asked a DoH server directly, bypassing your router's DNS entirely.

To see all DNS traffic, you need to disable DoH in your browser settings (see the README for instructions per browser).

## Filtered domains

`dnsw` hides two types of internal DNS traffic:

- **`.local` domains**: these are mDNS (multicast DNS) queries used by devices to discover each other on the local network. For example, your printer might announce itself as `my-printer.local`. These are not internet browsing.

- **`.arpa` domains**: these are reverse DNS lookups. When `dnsw` itself resolves a device IP to a hostname, that generates `.arpa` queries. Showing these would create a noisy feedback loop.

## Glossary

| Term | Meaning |
|------|---------|
| **BPF** | Berkeley Packet Filter. A kernel-level mechanism to filter network packets efficiently before they reach your program. |
| **DHCP** | Dynamic Host Configuration Protocol. How devices get their IP address from the router when they join the network. |
| **DNS** | Domain Name System. Translates human-readable domain names to IP addresses. |
| **DoH** | DNS-over-HTTPS. Encrypts DNS queries inside HTTPS, making them invisible to network sniffers. |
| **DoT** | DNS-over-TLS. Similar to DoH but uses a dedicated port (853) instead of mixing with web traffic. |
| **IP address** | A numerical label assigned to each device on a network. IPv4 looks like `192.168.1.1`, IPv6 looks like `fe80::1`. |
| **LAN** | Local Area Network. Your home network: your router and all devices connected to it. |
| **Loopback** | A virtual network interface (`127.0.0.1` or `::1`) where your machine talks to itself. |
| **mDNS** | Multicast DNS. Lets devices find each other on the local network without a central DNS server. |
| **Packet** | A small chunk of data sent over the network. Network traffic is split into many packets. |
| **Port** | A number that identifies a specific service on a machine. DNS uses port 53, HTTPS uses port 443. |
| **Promiscuous mode** | A mode where the network interface captures all packets on the network, not just those addressed to your machine. |
| **PTR record** | A DNS record for reverse lookups: maps an IP address back to a hostname. |
| **Reverse DNS** | Looking up the hostname for a given IP address (the opposite of normal DNS). |
| **Snapshot length** | How many bytes to capture from each packet. DNS packets are small, so we don't need to capture full packets. |
| **TCP** | Transmission Control Protocol. A reliable way to send data (checks if it arrived). Used by HTTP, email, etc. |
| **UDP** | User Datagram Protocol. A fast, simple way to send data without delivery guarantees. Used by DNS. |
