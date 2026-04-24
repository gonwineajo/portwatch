package history

import (
	"sort"
	"time"
)

// ExposureResult holds the cumulative open duration for a port on a host.
type ExposureResult struct {
	Host     string
	Port     int
	OpenedAt time.Time
	ClosedAt time.Time // zero if still open
	Duration time.Duration
	StillOpen bool
}

// AnalyseExposure calculates how long each port has been (or was) open
// across all entries. It pairs opened/closed events per host+port and
// accumulates total exposure time. Ports still open are measured against
// `now`.
func AnalyseExposure(entries []Entry, now time.Time) []ExposureResult {
	type key struct {
		host string
		port int
	}

	openedAt := make(map[key]time.Time)
	totals := make(map[key]time.Duration)

	// Process in chronological order.
	sorted := make([]Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp.Before(sorted[j].Timestamp)
	})

	for _, e := range sorted {
		switch e.Event {
		case EventOpened:
			for _, p := range e.Ports {
				k := key{e.Host, p}
				if _, alreadyOpen := openedAt[k]; !alreadyOpen {
					openedAt[k] = e.Timestamp
				}
			}
		case EventClosed:
			for _, p := range e.Ports {
				k := key{e.Host, p}
				if t, ok := openedAt[k]; ok {
					totals[k] += e.Timestamp.Sub(t)
					delete(openedAt, k)
				}
			}
		}
	}

	// Build results — still-open ports measured to now.
	seen := make(map[key]bool)
	var results []ExposureResult

	for k, d := range totals {
		seen[k] = true
		stillOpen := false
		var closedAt time.Time
		if extra, ok := openedAt[k]; ok {
			d += now.Sub(extra)
			stillOpen = true
		} else {
			closedAt = now // approximate; real closed time tracked above
		}
		results = append(results, ExposureResult{
			Host: k.host, Port: k.port,
			Duration: d, StillOpen: stillOpen, ClosedAt: closedAt,
		})
	}

	// Ports opened but never closed (no prior total entry).
	for k, t := range openedAt {
		if seen[k] {
			continue
		}
		results = append(results, ExposureResult{
			Host: k.host, Port: k.port,
			OpenedAt: t,
			Duration: now.Sub(t),
			StillOpen: true,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Duration != results[j].Duration {
			return results[i].Duration > results[j].Duration
		}
		if results[i].Host != results[j].Host {
			return results[i].Host < results[j].Host
		}
		return results[i].Port < results[j].Port
	})
	return results
}
