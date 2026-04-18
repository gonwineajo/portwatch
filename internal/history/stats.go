package history

import "sort"

// PortFrequency holds a port number and how many times it appeared in history.
type PortFrequency struct {
	Port  int
	Count int
}

// TopPorts returns the most frequently seen opened ports across all entries,
// limited to n results. Pass host="" to include all hosts.
func TopPorts(entries []Entry, host string, n int) []PortFrequency {
	counts := make(map[int]int)
	for _, e := range entries {
		if host != "" && e.Host != host {
			continue
		}
		if e.Event != "opened" {
			continue
		}
		for _, p := range e.Ports {
			counts[p]++
		}
	}

	freqs := make([]PortFrequency, 0, len(counts))
	for port, count := range counts {
		freqs = append(freqs, PortFrequency{Port: port, Count: count})
	}

	sort.Slice(freqs, func(i, j int) bool {
		if freqs[i].Count != freqs[j].Count {
			return freqs[i].Count > freqs[j].Count
		}
		return freqs[i].Port < freqs[j].Port
	})

	if n > 0 && len(freqs) > n {
		return freqs[:n]
	}
	return freqs
}

// HostActivity returns the number of change events per host.
func HostActivity(entries []Entry) map[string]int {
	activity := make(map[string]int)
	for _, e := range entries {
		activity[e.Host]++
	}
	return activity
}
