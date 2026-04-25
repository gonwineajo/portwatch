package history

import "math"

// EntropyResult holds the port-change entropy score for a single host.
// Higher entropy indicates more unpredictable / varied port activity.
type EntropyResult struct {
	Host    string
	Entropy float64 // Shannon entropy in bits
	Events  int     // total opened/closed events considered
}

// AnalyseEntropy computes the Shannon entropy of port-open/close events per
// host. Each distinct port is treated as a symbol; entropy measures how
// evenly activity is spread across ports.
//
// Only "opened" and "closed" events are considered. Hosts with fewer than
// minEvents events are excluded. Results are sorted descending by entropy.
func AnalyseEntropy(entries []Entry, minEvents int) []EntropyResult {
	type counter struct {
		total  int
		byPort map[int]int
	}

	hosts := map[string]*counter{}

	for _, e := range entries {
		if e.Event != EventOpened && e.Event != EventClosed {
			continue
		}
		c, ok := hosts[e.Host]
		if !ok {
			c = &counter{byPort: map[int]int{}}
			hosts[e.Host] = c
		}
		for _, p := range e.Ports {
			c.byPort[p]++
			c.total++
		}
	}

	var results []EntropyResult
	for host, c := range hosts {
		if c.total < minEvents {
			continue
		}
		var h float64
		for _, cnt := range c.byPort {
			p := float64(cnt) / float64(c.total)
			if p > 0 {
				h -= p * math.Log2(p)
			}
		}
		results = append(results, EntropyResult{
			Host:    host,
			Entropy: h,
			Events:  c.total,
		})
	}

	// sort descending by entropy, then alphabetically for stability
	for i := 1; i < len(results); i++ {
		for j := i; j > 0; j-- {
			a, b := results[j-1], results[j]
			if a.Entropy < b.Entropy || (a.Entropy == b.Entropy && a.Host > b.Host) {
				results[j-1], results[j] = results[j], results[j-1]
			}
		}
	}
	return results
}
