package history

import (
	"sort"
	"time"
)

// HostSummary holds aggregated change counts for a single host.
type HostSummary struct {
	Host        string
	Opened      int
	Closed      int
	LastChanged time.Time
}

// Summarize returns per-host change summaries from the provided entries,
// optionally filtered to only include entries after `since`.
func Summarize(entries []Entry, since time.Time) []HostSummary {
	type key = string
	type agg struct {
		opened  int
		closed  int
		lastAt  time.Time
	}

	m := make(map[key]*agg)

	for _, e := range entries {
		if e.Timestamp.Before(since) {
			continue
		}
		a, ok := m[e.Host]
		if !ok {
			a = &agg{}
			m[e.Host] = a
		}
		a.opened += len(e.Opened)
		a.closed += len(e.Closed)
		if e.Timestamp.After(a.lastAt) {
			a.lastAt = e.Timestamp
		}
	}

	result := make([]HostSummary, 0, len(m))
	for host, a := range m {
		result = append(result, HostSummary{
			Host:        host,
			Opened:      a.opened,
			Closed:      a.closed,
			LastChanged: a.lastAt,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Host < result[j].Host
	})

	return result
}
