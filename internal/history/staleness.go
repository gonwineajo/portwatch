package history

import (
	"sort"
	"time"
)

// StalenessResult holds staleness information for a host.
type StalenessResult struct {
	Host        string
	LastSeen    time.Time
	Staleness   time.Duration
	OpenPorts   []int
	IsStale     bool
}

// AnalyseStaleness returns staleness information for each host based on
// the age of their most recent scan entry. A host is considered stale
// when its last scan is older than the given threshold.
func AnalyseStaleness(entries []Entry, threshold time.Duration, now time.Time) []StalenessResult {
	type hostState struct {
		lastSeen  time.Time
		openPorts []int
	}

	states := make(map[string]*hostState)

	for _, e := range entries {
		if e.Event != EventScan {
			continue
		}
		s, ok := states[e.Host]
		if !ok || e.Timestamp.After(s.lastSeen) {
			states[e.Host] = &hostState{
				lastSeen:  e.Timestamp,
				openPorts: e.Ports,
			}
		}
	}

	results := make([]StalenessResult, 0, len(states))
	for host, s := range states {
		age := now.Sub(s.lastSeen)
		results = append(results, StalenessResult{
			Host:      host,
			LastSeen:  s.lastSeen,
			Staleness: age,
			OpenPorts: s.openPorts,
			IsStale:   age > threshold,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Staleness > results[j].Staleness
	})

	return results
}
