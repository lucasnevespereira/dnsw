# dnsw

A real-time DNS watcher for your local network. See what your devices are browsing, right from your terminal.

```
DNS WATCHER
────────────────────────────────────────────────────────
  Interface : en0
  Mode      : Wi-Fi (captures this machine's DNS queries)
  Press     : Ctrl+C to stop
────────────────────────────────────────────────────────

  ★ new device  My MacBook  IP 192.168.1.32  MAC 7A:7F:22:36:CA:66

TIME        DEVICE                CATEGORY      DOMAIN
────────────────────────────────────────────────────────────────────────
22:50:21    My MacBook            ◈ SOCIAL      facebook.com
22:50:22    My MacBook            ♪ MUSIC       spclient.spotify.com
22:50:24    My MacBook            ◎ DEV         github.com
22:50:30    My MacBook            ◎ AI          api.anthropic.com
22:50:42    My MacBook            ◈ SOCIAL      instagram.com
```

## Table of contents

- [Install](#install)
- [Usage](#usage)
- [Commands](#commands)
- [Proxy mode (all devices)](#proxy-mode)
- [Device identification](#device-identification)
- [Categories](#categories)
- [How it works](#how-it-works)
- [Docs](#docs)
- [License](#license)

## Install

### Prerequisites

- **macOS** or **Linux**
- **Go 1.21+** ([install Go](https://go.dev/dl/))
- **libpcap** (packet capture library)

```bash
# macOS (libpcap is included with Xcode command line tools)
xcode-select --install

# Linux (Debian/Ubuntu)
sudo apt install libpcap-dev
```

### Build

```bash
git clone https://github.com/lucasnevespereira/dnsw.git
cd dnsw
go build -o dnsw .
```

## Usage

```bash
# Start watching DNS (auto-detects your network interface)
sudo ./dnsw
```

That's it. `sudo` is required to capture network packets.

### Commands

```bash
dnsw                         # start watching (default)
dnsw watch                   # same as above
dnsw watch -i en0            # use a specific network interface
dnsw watch --no-dedupe       # show all DNS packets without merging

dnsw proxy                   # run a DNS proxy to see ALL devices on your network
dnsw proxy --upstream 1.1.1.1  # use a custom upstream DNS server

dnsw devices list            # show all named devices
dnsw devices set <ip> <name> # name a device
dnsw devices remove <ip>     # remove a device name

dnsw interfaces              # list available network interfaces
dnsw config                  # show config file path and contents
```

## Proxy mode

On Wi-Fi, `dnsw watch` can only see DNS queries from the machine it runs on. Wi-Fi encrypts each device's traffic separately, so your computer can't see packets from your phone or TV.

`dnsw proxy` solves this by turning your Mac into a DNS server. All devices on the network send their DNS queries through your machine, so you see everything:

```
DNS WATCHER (proxy mode)
────────────────────────────────────────────────────────
  DNS proxy  : 192.168.1.32:53
  Upstream   : 8.8.8.8:53
────────────────────────────────────────────────────────

  ★ new device  Apple-1  IP 192.168.1.73

23:15:01    Apple-1               ◈ SOCIAL      facebook.com
23:15:03    Apple-1               ♪ MUSIC       spclient.spotify.com
23:15:05    Samsung-1             ▶ VIDEO       youtube.com
```

To use proxy mode, you need to point your router's DNS to your Mac's IP. This takes about 2 minutes. See the [Router Setup Guide](docs/router-setup.md) for step-by-step instructions.

**Safe to stop at any time.** Set `8.8.8.8` as the secondary DNS in your router, and internet keeps working even when `dnsw proxy` is not running.

## Device identification

`dnsw` identifies devices using three methods (in order):

1. **Your names** from `~/.config/dnsw/devices.json` (highest priority)
2. **Reverse DNS** (hostnames assigned by your router)
3. **MAC vendor detection** (reads the device manufacturer from its network address)

When a new device appears, you'll see a notice:

```
  ★ new device  Apple-1  IP 192.168.1.73  MAC AC:DE:48:00:11:22
```

Devices are auto-named by manufacturer: `Apple-1`, `Samsung-2`, `Google-3`, etc.

To give devices custom names:

```bash
dnsw devices set 192.168.1.32 "My MacBook"
dnsw devices set 192.168.1.45 "iPhone"
```

Or edit `~/.config/dnsw/devices.json` directly.

> **Note**: modern phones (iOS 14+, Android 10+) use randomized MAC addresses by default. These won't match any known vendor, so they'll show as the raw IP unless you name them with `dnsw devices set`.

## Categories

Domains are auto-categorized with icons:

| Icon         | Category        | Examples                            |
| ------------ | --------------- | ----------------------------------- |
| `▶ VIDEO`   | Streaming       | YouTube, Netflix, Twitch, Disney+   |
| `♪ MUSIC`    | Music           | Spotify, Deezer, SoundCloud         |
| `◈ SOCIAL`   | Social media    | Facebook, Instagram, TikTok, Reddit |
| `⌕ SEARCH`   | Search engines  | Google, Bing                        |
| `⊞ SHOP`     | Shopping        | Amazon, eBay                        |
| `◎ COMM`     | Communication   | Discord, WhatsApp, Slack, Zoom      |
| `◎ DEV`      | Development     | GitHub, GitLab, JetBrains           |
| `◎ AI`       | AI              | Claude, ChatGPT                     |
| `⚔ GAMING`  | Gaming          | Steam, Epic, PlayStation, Xbox      |
| `✗ ADS/TRCK` | Ads & trackers  | DoubleClick, Google Ads             |
| `◉ APPLE`    | Apple services  | iCloud, App Store                   |
| `◉ MSFT`     | Microsoft       | Office, Windows Update              |
| `⊡ CLOUD`    | Cloud infra     | Cloudflare, AWS, Akamai             |
| `· OTHER`    | Everything else |                                     |

## How it works

`dnsw` captures DNS queries (UDP port 53) on your network, identifies devices, categorizes domains, and displays everything in real time. Duplicate queries are merged and internal lookups (`.local`, `.arpa`) are filtered out.

> **Not seeing some websites?** Modern browsers use DNS-over-HTTPS (DoH) which bypasses standard DNS. See [docs/network.md](docs/network.md#why-some-domains-dont-show-up-dns-over-https) for how to disable it.

## Docs

- [Network Concepts](docs/network.md) - beginner-friendly explanation of DNS, packets, and how `dnsw` works under the hood
- [Router Setup Guide](docs/router-setup.md) - step-by-step instructions to configure your router for proxy mode

## License

[MIT](LICENSE)
