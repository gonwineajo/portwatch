package history

import (
	"sort"
	"time"
)

// SnapshotDiffEntry records the port-level difference between two consecutive
// scans for a single host.
type SnapshotDiffEntry struct {
	Host    string
	At      time.Time
	Opened  []int
	Closed  []int
	Stable  []int // ports present in both snapshots
}

// SnapshotDiffs computes per-host diffs between consecutive scan entries in
// the provided history. Only entries with EventType "scan" are considered.
// Results are returned in chronological order.
func SnapshotDiffs(entries []Entry) []SnapshotDiffEntry {
	// Group scan entries by host, preserving order.
	byHost := make(map[string][]Entry)
	var hosts []string
	for _, e := range entries {
		if e.Event != "scan" {
			continue
		}
		if _, seen := byHost[e.Host]; !seen {
			hosts = append(hosts, e.Host)
		}
		byHost[e.Host] = append(byHost[e.Host], e)
	}
	sort.Strings(hosts)

	var result []SnapshotDiffEntry
	for _, host := range hosts {
		scans := byHost[host]
		for i := 1; i < len(scans); i++ {
			prev := toPortSet(scans[i-1].Ports)
			curr := toPortSet(scans[i].Ports)

			d := SnapshotDiffEntry{
				Host: host,
				At:   scans[i].Timestamp,
			}
			for p := range curr {
				if prev[p] {
					d.Stable = append(d.Stable, p)
				} else {
					d.Opened = append(d.Opened, p)
				}
			}
			for p := range prev {
				if !curr[p] {
					d.Closed = append(d.Closed, p)
				}
			}
			sort.Ints(d.Opened)
			sort.Ints(d.Closed)
			sort.Ints(d.Stable)
			result = append(result, d)
		}
	}
	return result
}

func toPortSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
