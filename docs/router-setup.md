# Router Setup Guide

This guide explains how to configure your router so `dnsw proxy` can see DNS queries from **all** devices on your network (phones, TVs, tablets, etc).

You only need to change one setting: tell your router to use your computer as the DNS server.

## Before you start

1. Run `dnsw proxy` on your Mac. It will show your local IP:

```
DNS WATCHER (proxy mode)
────────────────────────────────────────────────────────
  DNS proxy  : 192.168.1.32:53
  Upstream   : 8.8.8.8:53
```

2. Note down your IP (in this example: `192.168.1.32`)

## Step 1: Open your router admin page

Open a browser and go to one of these addresses:

- `http://192.168.1.1` (most common)
- `http://192.168.0.1`
- `http://10.0.0.1`

You'll see a login page. The default credentials are usually printed on a sticker on the bottom of your router. Common defaults:

| Brand | Username | Password |
|-------|----------|----------|
| Most routers | admin | admin |
| Netgear | admin | password |
| TP-Link | admin | admin |
| Livebox (Orange) | admin | (on sticker) |
| Freebox | (no login) | mafreebox.freebox.fr |

## Step 2: Find DNS settings

The setting is in different places depending on your router, but it's always under LAN or DHCP settings:

- **Generic**: LAN > DHCP Server > DNS
- **Netgear**: Internet > DNS Address
- **TP-Link**: DHCP > DHCP Settings > Primary DNS
- **Livebox**: Advanced Settings > DHCP
- **Freebox**: Settings > DHCP > DNS
- **UniFi**: Settings > Networks > (your network) > DHCP > DNS Server

Look for fields called "DNS Server", "Primary DNS", or "DNS Address".

## Step 3: Set the DNS servers

Set two DNS servers:

| Field | Value | Why |
|-------|-------|-----|
| Primary DNS | Your Mac's IP (e.g. `192.168.1.32`) | Sends all queries through `dnsw proxy` |
| Secondary DNS | `8.8.8.8` | Fallback if `dnsw proxy` is not running |

The secondary DNS is important: if you stop `dnsw proxy` or your Mac goes to sleep, devices will automatically fall back to Google's DNS (`8.8.8.8`) and internet keeps working.

## Step 4: Save and wait

Save the settings. Devices will pick up the new DNS within about 1 minute. Some devices may need to reconnect to Wi-Fi or toggle airplane mode to get the new settings immediately.

## Verifying it works

After the setup, open a browser on your phone and visit any website. You should see the query appear in `dnsw proxy`:

```
  ★ new device  Apple-1  IP 192.168.1.73

23:15:01    Apple-1               ◈ SOCIAL      facebook.com
23:15:03    Apple-1               ♪ MUSIC       spclient.spotify.com
```

## Undoing the change

To go back to normal, just change the DNS settings in your router back to automatic (or remove the custom DNS entries). Your router will go back to using your ISP's default DNS.

You can also simply stop using `dnsw proxy`. As long as you set `8.8.8.8` as the secondary DNS, internet will keep working through Google's DNS.

## Troubleshooting

**"I can't access my router admin page"**
Try `192.168.0.1` or `10.0.0.1`. You can also find it by running `route get default` in Terminal and looking at the "gateway" line.

**"Devices still don't show up after changing DNS"**
Some devices cache DNS settings. Try:
- Toggle airplane mode on your phone
- Disconnect and reconnect Wi-Fi
- Wait 2 minutes

**"Internet stopped working on all devices"**
`dnsw proxy` is probably not running. Either:
- Start it again: `sudo dnsw proxy`
- Or revert the router DNS to automatic

**"dnsw proxy says 'address already in use'"**
Something else is using port 53 on your Mac. Check with `sudo lsof -i :53`. On macOS, you may need to disable the built-in mDNSResponder temporarily.
