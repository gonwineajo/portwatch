package history

import "sort"

// ChurnResult holds the port churn score for a single host.
// Churn is defined as the total number of open+close events
// divided by the number of distinct ports seen, giving a
// normalised measure of how "unstable" a host's port surface is.
type ChurnResult struct {
	Host       string
	TotalFlips int     // total opened+closed events
	UniquePorts int    // distinct ports that changed
	Score      float64 // TotalFlips / UniquePorts (or 0)
}

// AnalyseChurn computes a churn score for every host found in
// entries. Only "opened" and "closed" events are counted.
// Results are returned sorted descending by Score.
func AnalyseChurn(entries []Entry) []ChurnResult {
	type hostData struct {
		flips int
		ports map[int]struct{}
	}
	acc := map[string]*hostData{}

	for _, e := range entries {
		if e.Event != EventOpened && e.Event != EventClosed {
			continue
		}
		hd, ok := acc[e.Host]
		if !ok {
			hd = &hostData{ports: map[int]struct{}{}}
			acc[e.Host] = hd
		}
		hd.flips++
		for _, p := range e.Ports {
			hd.ports[p] = struct{}{}
		}
	}

	results := make([]ChurnResult, 0, len(acc))
	for host, hd := range acc {
		score := 0.0
		if len(hd.ports) > 0 {
			score = float64(hd.flips) / float64(len(hd.ports))
		}
		results = append(results, ChurnResult{
			Host:        host,
			TotalFlips:  hd.flips,
			UniquePorts: len(hd.ports),
			Score:       score,
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
