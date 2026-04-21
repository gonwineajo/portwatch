// Package history provides hotspot detection for identifying unstable host+port
// pairs that frequently transition between open and closed states.
//
// # Hotspot Detection
//
// DetectHotspots analyses a sequence of history entries and counts the number
// of open<->close transitions ("flips") for each unique host+port pair. Pairs
// that exceed a minimum flip threshold are returned as Hotspot records, sorted
// by descending activity score.
//
// Example usage:
//
//	entries, _ := history.Read("portwatch.json")
//	hotspots := history.DetectHotspots(entries, 2)
//	for _, h := range hotspots {
//		fmt.Printf("%s:%d — %d flips\n", h.Host, h.Port, h.Flips)
//	}
//
// Hotspots are useful for identifying services that are frequently restarted,
// flapping, or subject to repeated external probing.
package history
