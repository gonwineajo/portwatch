package history

import "sort"

// CompactResult holds the outcome of a compaction operation.
type CompactResult struct {
	Before  int
	After   int
	Removed int
}

// Compact deduplicates consecutive entries for the same host that carry
// identical open-port sets, keeping only the first occurrence of each
// run. This reduces noise in long-running histories where ports are
// stable across many scans.
//
// Entries are processed in chronological order. The relative ordering
// of surviving entries is preserved.
//
// Only EventScan entries are eligible for deduplication; all other event
// types (e.g. EventOpen, EventClose) are always retained.
func Compact(entries []Entry) ([]Entry, CompactResult) {
	result := CompactResult{Before: len(entries)}
	if len(entries) == 0 {
		return entries, result
	}

	// Sort by timestamp ascending so we process runs in order.
	sorted := make([]Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp.Before(sorted[j].Timestamp)
	})

	// last tracks the most-recently-kept port fingerprint per host.
	last := make(map[string]string)
	out := make([]Entry, 0, len(sorted))

	for _, e := range sorted {
		fp := portFingerprint(e.OpenPorts)
		prev, seen := last[e.Host]
		if seen && prev == fp && e.Event == EventScan {
			// Identical stable scan — skip.
			continue
		}
		// Always update the fingerprint so that a non-scan event followed
		// by a scan with the same ports is still deduplicated correctly.
		last[e.Host] = fp
		out = append(out, e)
	}

	result.After = len(out)
	result.Removed = result.Before - result.After
	return out, result
}

// portFingerprint returns a stable string key for a slice of port numbers.
func portFingerprint(ports []int) string {
	if len(ports) == 0 {
		return ""
	}
	sorted := make([]int, len(ports))
	copy(sorted, ports)
	sort.Ints(sorted)
	return joinInts(sorted, ",")
}
