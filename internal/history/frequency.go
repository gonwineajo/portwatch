package history

import "sort"

// PortFrequency holds the frequency statistics for a single port across all hosts.
type PortFrequency struct {
	Port       int
	OpenCount  int
	CloseCount int
	Hosts      []string
}

// AnalyseFrequency counts how many times each port was opened or closed
// across all entries. Only "opened" and "closed" events are considered.
// Results are sorted by OpenCount descending, then by Port ascending.
func AnalyseFrequency(entries []Entry) []PortFrequency {
	type portStats struct {
		openCount  int
		closeCount int
		hosts      map[string]struct{}
	}

	stats := make(map[int]*portStats)

	for _, e := range entries {
		if e.Event != EventOpened && e.Event != EventClosed {
			continue
		}
		ports := e.Ports
		if e.Event == EventClosed {
			ports = e.Ports
		}
		for _, p := range ports {
			if _, ok := stats[p]; !ok {
				stats[p] = &portStats{hosts: make(map[string]struct{})}
			}
			if e.Event == EventOpened {
				stats[p].openCount++
			} else {
				stats[p].closeCount++
			}
			stats[p].hosts[e.Host] = struct{}{}
		}
	}

	result := make([]PortFrequency, 0, len(stats))
	for port, s := range stats {
		hosts := make([]string, 0, len(s.hosts))
		for h := range s.hosts {
			hosts = append(hosts, h)
		}
		sort.Strings(hosts)
		result = append(result, PortFrequency{
			Port:       port,
			OpenCount:  s.openCount,
			CloseCount: s.closeCount,
			Hosts:      hosts,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].OpenCount != result[j].OpenCount {
			return result[i].OpenCount > result[j].OpenCount
		}
		return result[i].Port < result[j].Port
	})

	return result
}
