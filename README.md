# dnsw

A real-time DNS watcher for your local network. See what every device on your home network is browsing, right from your terminal.

```
DNS WATCHER
────────────────────────────────────────────────────
  Interface : en0
  Mode      : promiscuous, capturing ALL devices on LAN
  Press     : Ctrl+C to stop
────────────────────────────────────────────────────

TIME        DEVICE                CATEGORY      DOMAIN
────────────────────────────────────────────────────────────────────────
22:50:21    My MacBook            ◈ SOCIAL      facebook.com
22:50:22    iPhone                ♪ MUSIC       spclient.spotify.com
22:50:24    My MacBook            ◎ DEV         github.com
22:50:30    Living Room TV        ▶ VIDEO       netflix.com
22:50:42    iPhone                ◈ SOCIAL      instagram.com
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

That's it. `sudo` is required to capture packets in promiscuous mode (seeing all devices, not just yours).

### Flags

| Flag             | Description                                                                |
| ---------------- | -------------------------------------------------------------------------- |
| `-i <interface>` | Use a specific network interface instead of auto-detecting                 |
| `--list`         | Show all available network interfaces                                      |
| `--no-dedupe`    | Show every DNS packet (by default, duplicate queries within 2s are merged) |

### Choosing a different interface

If auto-detection picks the wrong one, run `sudo ./dnsw --list` to see what's available:

```
Available interfaces:

  en0           192.168.1.32
  en1           192.168.1.45
  bridge0       192.168.2.1
```

Then specify it with `-i`:

```bash
./dnsw -i en1
```

### Device identification

`dnsw` automatically identifies devices on your network using three methods (in order):

1. **Your names** from `~/.config/dnsw/devices.json` (highest priority)
2. **Reverse DNS** (hostnames assigned by your router)
3. **MAC vendor detection** (reads the device manufacturer from its network address)

When a new device appears, you'll see a notice:

```
  ★ new device  Apple-1  IP 192.168.1.73  MAC AC:DE:48:00:11:22
```

Devices are auto-named by manufacturer: `Apple-1`, `Samsung-2`, `Google-3`, etc. This covers most phones, tablets, smart TVs, game consoles, and IoT devices.

To give devices custom names, edit `~/.config/dnsw/devices.json` (created automatically on first run):

```json
{
  "192.168.1.32": "My MacBook",
  "192.168.1.45": "iPhone",
  "192.168.1.50": "Living Room TV"
}
```

> **Tip**: watch the new device notices to see which IPs appear, then map them in the config.

> **Note**: modern phones (iOS 14+, Android 10+) use randomized MAC addresses by default. These won't match any known vendor, so they'll show as the raw IP unless you name them in `devices.json`.

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

`dnsw` puts your network interface into **promiscuous mode**, meaning it sees all packets on the local network, not just those addressed to your machine. It then filters for **UDP port 53** (standard DNS), parses the DNS query packets, and displays them in a formatted table.

- Only **DNS queries** are shown (not responses)
- `.local` and `.arpa` domains are filtered out (internal network lookups)
- Devices are identified via config file, reverse DNS, or MAC vendor lookup
- New devices are announced when first seen on the network
- Duplicate queries from the same device within 2 seconds are merged into one line

For a beginner-friendly explanation of the networking concepts used in this project, see [NETWORK.md](NETWORK.md).

## License

[MIT](LICENSE)
