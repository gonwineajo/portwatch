package history

import "sort"

// PortFrequency holds a port number and how many times it appeared.
type PortFrequency struct {
	Port  int
	Count int
}

// HostSummary holds per-host open/close event totals.
type HostSummary struct {
	Host   string
	Opened int
	Closed int
	Total  int
}

// AggregateByHost returns a HostSummary for each distinct host in entries.
func AggregateByHost(entries []Entry) []HostSummary {
	type counts struct{ opened, closed int }
	m := map[string]*counts{}
	for _, e := range entries {
		if _, ok := m[e.Host]; !ok {
			m[e.Host] = &counts{}
		}
		if e.Event == "opened" {
			m[e.Host].opened++
		} else if e.Event == "closed" {
			m[e.Host].closed++
		}
	}
	out := make([]HostSummary, 0, len(m))
	for host, c := range m {
		out = append(out, HostSummary{
			Host:   host,
			Opened: c.opened,
			Closed: c.closed,
			Total:  c.opened + c.closed,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Total != out[j].Total {
			return out[i].Total > out[j].Total
		}
		return out[i].Host < out[j].Host
	})
	return out
}

// AggregateByPort returns port frequencies across all entries, sorted descending.
func AggregateByPort(entries []Entry) []PortFrequency {
	m := map[int]int{}
	for _, e := range entries {
		for _, p := range e.Ports {
			m[p]++
		}
	}
	out := make([]PortFrequency, 0, len(m))
	for port, count := range m {
		out = append(out, PortFrequency{Port: port, Count: count})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].Port < out[j].Port
	})
	return out
}
