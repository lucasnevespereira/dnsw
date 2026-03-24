package main

import (
	"sync"
	"time"
)

const (
	// When a device queries the same domain multiple times within this window,
	// we only show it once. DNS clients often send A + AAAA queries simultaneously
	// (IPv4 + IPv6), which would otherwise clutter the output.
	dedupeWindow = 2 * time.Second

	// Cleanup threshold: when the cache grows beyond this, we prune old entries.
	dedupeMaxEntries = 500

	// Entries older than this get cleaned up during pruning.
	dedupeCleanupAge = 10 * time.Second
)

// recentQueries tracks when we last displayed a query from a given device+domain pair.
var (
	recentQueries = map[string]time.Time{}
	dedupeMu      sync.Mutex
)

// isDuplicate returns true if we've already shown this device+domain combo recently.
func isDuplicate(device, domain string) bool {
	key := device + "|" + domain
	now := time.Now()

	dedupeMu.Lock()
	defer dedupeMu.Unlock()

	if last, ok := recentQueries[key]; ok && now.Sub(last) < dedupeWindow {
		recentQueries[key] = now
		return true
	}
	recentQueries[key] = now

	// Periodically clean old entries to prevent unbounded memory growth.
	if len(recentQueries) > dedupeMaxEntries {
		for k, t := range recentQueries {
			if now.Sub(t) > dedupeCleanupAge {
				delete(recentQueries, k)
			}
		}
	}

	return false
}
