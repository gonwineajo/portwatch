package history

import "sort"

// DriftResult describes how much a host's open port set has changed
// relative to its established baseline.
type DriftResult struct {
	Host    string
	Added   []int
	Removed []int
	Score   float64 // (|added| + |removed|) / max(1, |baseline|)
}

// AnalyseDrift compares each host's latest scan against its baseline
// snapshot and returns a ranked list of hosts ordered by drift score
// (highest first). Hosts with no baseline are skipped.
func AnalyseDrift(entries []Entry) []DriftResult {
	baselines := SetBaseline(entries)
	latest := latestScanPortsMap(entries)

	var results []DriftResult

	for host, base := range baselines {
		current, ok := latest[host]
		if !ok {
			continue
		}

		baseSet := toIntSet(base.Ports)
		curSet := toIntSet(current)

		var added, removed []int
		for p := range curSet {
			if !baseSet[p] {
				added = append(added, p)
			}
		}
		for p := range baseSet {
			if !curSet[p] {
				removed = append(removed, p)
			}
		}

		sort.Ints(added)
		sort.Ints(removed)

		denominator := float64(len(baseSet))
		if denominator < 1 {
			denominator = 1
		}
		score := float64(len(added)+len(removed)) / denominator

		results = append(results, DriftResult{
			Host:    host,
			Added:   added,
			Removed: removed,
			Score:   score,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		return results[i].Host < results[j].Host
	})

	return results
}
