package history

import "sort"

// DominanceResult holds the dominance score for a single port across all hosts.
type DominanceResult struct {
	Port      int
	HostCount int    // number of distinct hosts that have opened this port
	TotalOpen int    // total open events across all hosts
	Score     float64 // HostCount * TotalOpen normalised
}

// AnalyseDominance ranks ports by how widely and frequently they appear as
// opened across all scanned hosts. Only "opened" events are considered.
// minHosts filters out ports seen on fewer than minHosts distinct hosts.
func AnalyseDominance(entries []Entry, minHosts int) []DominanceResult {
	type portStats struct {
		hosts map[string]struct{}
		total int
	}

	stats := make(map[int]*portStats)

	for _, e := range entries {
		if e.Event != "opened" {
			continue
		}
		for _, p := range e.Ports {
			if _, ok := stats[p]; !ok {
				stats[p] = &portStats{hosts: make(map[string]struct{})}
			}
			stats[p].hosts[e.Host] = struct{}{}
			stats[p].total++
		}
	}

	var results []DominanceResult
	for port, s := range stats {
		hc := len(s.hosts)
		if hc < minHosts {
			continue
		}
		results = append(results, DominanceResult{
			Port:      port,
			HostCount: hc,
			TotalOpen: s.total,
			Score:     float64(hc) * float64(s.total),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		return results[i].Port < results[j].Port
	})

	return results
}
