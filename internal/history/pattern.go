package history

import "sort"

// PatternResult holds recurrence info for a port on a host.
type PatternResult struct {
	Host      string
	Port      int
	OpenCount int
	CloseCount int
	Recurring bool // true if opened+closed more than once
}

// DetectPatterns identifies ports that repeatedly open and close across history,
// which may indicate flapping services or scheduled processes.
func DetectPatterns(entries []Entry) []PatternResult {
	type key struct {
		host string
		port int
	}

	opens := make(map[key]int)
	closes := make(map[key]int)

	for _, e := range entries {
		for _, p := range e.OpenedPorts {
			opens[key{e.Host, p}]++
		}
		for _, p := range e.ClosedPorts {
			closes[key{e.Host, p}]++
		}
	}

	seen := make(map[key]bool)
	var results []PatternResult

	for k, oc := range opens {
		seen[k] = true
		cc := closes[k]
		results = append(results, PatternResult{
			Host:       k.host,
			Port:       k.port,
			OpenCount:  oc,
			CloseCount: cc,
			Recurring:  oc > 1 && cc > 1,
		})
	}

	for k, cc := range closes {
		if seen[k] {
			continue
		}
		results = append(results, PatternResult{
			Host:       k.host,
			Port:       k.port,
			OpenCount:  0,
			CloseCount: cc,
			Recurring:  false,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Host != results[j].Host {
			return results[i].Host < results[j].Host
		}
		return results[i].Port < results[j].Port
	})

	return results
}

// RecurringOnly filters PatternResults to those marked as recurring.
func RecurringOnly(patterns []PatternResult) []PatternResult {
	var out []PatternResult
	for _, p := range patterns {
		if p.Recurring {
			out = append(out, p)
		}
	}
	return out
}
