// Package history provides utilities for recording, querying, and analysing
// port-scan history produced by portwatch.
//
// Watchlist
//
// The Watchlist type lets operators flag specific host:port pairs for closer
// attention. Entries can be persisted to disk with SaveWatchlist /
// LoadWatchlist and matched against a slice of history Entry values with
// Watchlist.Match.
//
// Example usage:
//
//	wl, _ := history.LoadWatchlist("/var/lib/portwatch/watchlist.json")
//	wl.Add("prod-db", 5432, "PostgreSQL")
//	_ = history.SaveWatchlist("/var/lib/portwatch/watchlist.json", wl)
//
//	entries, _ := history.Read("/var/lib/portwatch/history.json")
//	matched := wl.Match(entries)
package history
