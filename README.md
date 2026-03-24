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

dnsw devices list            # show all named devices
dnsw devices set <ip> <name> # name a device
dnsw devices remove <ip>     # remove a device name

dnsw interfaces              # list available network interfaces
dnsw config                  # show config file path and contents
```

### Device identification

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

## Wi-Fi limitation

On **Wi-Fi**, `dnsw` can only see DNS queries from the machine it's running on. This is because Wi-Fi encrypts each device's traffic separately, so your computer physically can't see packets from your phone or TV.

On a **wired network** (Ethernet), promiscuous mode can capture traffic from other devices on the same network segment.

To see DNS queries from **all** devices on your network, you have a few options:

- **Router DNS logs**: some routers (OpenWrt, pfSense, UniFi) can show DNS query logs directly
- **Pi-hole**: runs as your network's DNS server and logs all queries with a web dashboard
- **Run dnsw on your router**: if your router runs Linux (e.g. OpenWrt), you can run `dnsw` there where it can see all DNS traffic

## Why can't I see some websites?

Modern browsers use **DNS-over-HTTPS (DoH)**, which encrypts DNS queries inside regular HTTPS traffic. When DoH is active, DNS lookups bypass the standard UDP port 53, so `dnsw` can't see them.

### How to disable DoH to see all DNS traffic

**Chrome / Brave / Edge:**

1. Go to `chrome://settings/security` (or `brave://settings/security`, `edge://settings/privacy`)
2. Find **"Use secure DNS"**
3. Turn it **off**

**Firefox:**

1. Go to `about:preferences#general`
2. Scroll to **Network Settings**, click **Settings**
3. Uncheck **"Enable DNS over HTTPS"**

**Safari:**
Safari uses the system DNS settings by default and does **not** enable DoH, so it should work out of the box.

**macOS system-wide:**
macOS itself does not use DoH by default. System-level DNS (like `curl` in terminal) will appear in `dnsw` without any changes.

> After disabling DoH, all browser DNS queries go through standard UDP port 53 and will show up in `dnsw`.

## How it works

`dnsw` listens on your network interface for **UDP port 53** (standard DNS) traffic, parses the DNS query packets, and displays them in a formatted table.

- Only **DNS queries** are shown (not responses)
- `.local` and `.arpa` domains are filtered out (internal network lookups)
- Devices are identified via config file, reverse DNS, or MAC vendor lookup
- New devices are announced when first seen
- Duplicate queries from the same device within 2 seconds are merged into one line

For a beginner-friendly explanation of the networking concepts used in this project, see [NETWORK.md](NETWORK.md).

## License

[MIT](LICENSE)
