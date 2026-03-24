package main

import "strings"

// category represents a visual grouping for a domain (icon + color).
type category struct {
	label string
	color string
}

// domainCategories maps known domain patterns to their visual category.
// When a DNS query matches one of these keywords, it gets a colored label
// in the output so you can quickly see what kind of traffic it is.
var domainCategories = []struct {
	keywords []string
	cat      category
}{
	// Video streaming
	{[]string{"youtube.com", "youtu.be", "ytimg.com", "googlevideo.com"}, category{"▶ VIDEO   ", cyan}},
	{[]string{"netflix.com", "nflxvideo.net", "nflximg.net"}, category{"▶ VIDEO   ", cyan}},
	{[]string{"twitch.tv", "twitchsvc.net", "jtvnw.net"}, category{"▶ VIDEO   ", cyan}},
	{[]string{"disneyplus.com", "disney-plus.net", "bamgrid.com"}, category{"▶ VIDEO   ", cyan}},
	{[]string{"primevideo.com", "aiv-cdn.net"}, category{"▶ VIDEO   ", cyan}},

	// Music
	{[]string{"spotify.com", "scdn.co", "spotifycdn.com"}, category{"♪ MUSIC   ", magenta}},
	{[]string{"deezer.com", "dzcdn.net"}, category{"♪ MUSIC   ", magenta}},
	{[]string{"soundcloud.com", "sndcdn.com"}, category{"♪ MUSIC   ", magenta}},

	// Social media
	{[]string{"facebook.com", "fbcdn.net", "fb.com"}, category{"◈ SOCIAL  ", yellow}},
	{[]string{"instagram.com", "cdninstagram.com"}, category{"◈ SOCIAL  ", yellow}},
	{[]string{"twitter.com", "twimg.com", "x.com"}, category{"◈ SOCIAL  ", yellow}},
	{[]string{"tiktok.com", "tiktokcdn.com"}, category{"◈ SOCIAL  ", yellow}},
	{[]string{"reddit.com", "redditmedia.com", "redditstatic.com"}, category{"◈ SOCIAL  ", yellow}},
	{[]string{"snapchat.com", "snap.com", "snapads.com"}, category{"◈ SOCIAL  ", yellow}},
	{[]string{"linkedin.com", "licdn.com"}, category{"◈ SOCIAL  ", yellow}},

	// Search engines
	{[]string{"google.com", "googleapis.com", "gstatic.com", "google.fr"}, category{"⌕ SEARCH  ", blue}},
	{[]string{"bing.com", "msn.com"}, category{"⌕ SEARCH  ", blue}},

	// Shopping
	{[]string{"amazon.com", "amazon.fr", "amazon.co"}, category{"⊞ SHOP    ", green}},
	{[]string{"ebay.com", "ebay.fr"}, category{"⊞ SHOP    ", green}},

	// Cloud infrastructure
	{[]string{"cloudflare.com", "cloudflare-dns.com"}, category{"⊡ CLOUD   ", green}},
	{[]string{"amazonaws.com", "aws.amazon.com"}, category{"⊡ CLOUD   ", green}},
	{[]string{"akamai.com", "akamaized.net", "akamaihd.net"}, category{"⊡ CLOUD   ", green}},

	// Ads & trackers
	{[]string{"doubleclick.net", "googlesyndication.com", "adnxs.com"}, category{"✗ ADS/TRCK", red}},
	{[]string{"googleadservices.com", "googletag", "adsrvr.org"}, category{"✗ ADS/TRCK", red}},
	{[]string{"facebook.net", "fbsbx.com"}, category{"✗ ADS/TRCK", red}},

	// Apple
	{[]string{"apple.com", "icloud.com", "mzstatic.com", "apple-dns.net"}, category{"◉ APPLE   ", white}},

	// Microsoft
	{[]string{"microsoft.com", "windows.com", "office.com", "live.com"}, category{"◉ MSFT    ", white}},

	// Development
	{[]string{"github.com", "githubusercontent.com", "gitlab.com"}, category{"◎ DEV     ", green}},
	{[]string{"stackoverflow.com", "stackexchange.com"}, category{"◎ DEV     ", green}},
	{[]string{"npmjs.org", "npmjs.com", "yarnpkg.com"}, category{"◎ DEV     ", green}},
	{[]string{"jetbrains.com"}, category{"◎ DEV     ", green}},
	{[]string{"anthropic.com", "claude.ai", "openai.com", "chatgpt.com"}, category{"◎ AI      ", green}},

	// Communication
	{[]string{"discord.com", "discordapp.com", "discord.gg", "discord.media"}, category{"◎ COMM    ", blue}},
	{[]string{"whatsapp.com", "whatsapp.net"}, category{"◎ COMM    ", blue}},
	{[]string{"slack.com", "slack-edge.com"}, category{"◎ COMM    ", blue}},
	{[]string{"telegram.org", "t.me"}, category{"◎ COMM    ", blue}},
	{[]string{"zoom.us", "zoom.com"}, category{"◎ COMM    ", blue}},

	// Gaming
	{[]string{"steampowered.com", "steamcommunity.com", "steamcdn"}, category{"⚔ GAMING  ", magenta}},
	{[]string{"epicgames.com", "unrealengine.com"}, category{"⚔ GAMING  ", magenta}},
	{[]string{"playstation.com", "playstation.net", "sonyentertainmentnetwork"}, category{"⚔ GAMING  ", magenta}},
	{[]string{"xbox.com", "xboxlive.com"}, category{"⚔ GAMING  ", magenta}},
	{[]string{"riotgames.com", "leagueoflegends.com"}, category{"⚔ GAMING  ", magenta}},
}

// categorize checks if a domain matches any known category.
// Returns the matching category, or "OTHER" if unknown.
func categorize(domain string) category {
	d := strings.ToLower(domain)
	for _, entry := range domainCategories {
		for _, kw := range entry.keywords {
			if strings.Contains(d, kw) {
				return entry.cat
			}
		}
	}
	return category{"· OTHER   ", gray}
}
