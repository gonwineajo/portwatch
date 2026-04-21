package history

import "time"

// AnomalyReport describes an unusual port event for a host.
type AnomalyReport struct {
	Host      string
	Port      int
	Event     string
	OccurredAt time.Time
	Reason    string
}

// DetectAnomalies scans entries for ports that appear or disappear outside
// their established recurring schedule. A port is considered anomalous when
// it fires an opened/closed event but has fewer than minCount prior
// occurrences in the history, making it a rare / unexpected change.
func DetectAnomalies(entries []Entry, minCount int) []AnomalyReport {
	if minCount <= 0 {
		minCount = 2
	}

	// Count how many times each (host, port, event) tuple has been seen.
	type key struct {
		host  string
		port  int
		event string
	}
	counts := make(map[key]int)
	for _, e := range entries {
		if e.Event == EventNoChange {
			continue
		}
		for _, p := range e.OpenedPorts {
			counts[key{e.Host, p, EventOpened}]++
		}
		for _, p := range e.ClosedPorts {
			counts[key{e.Host, p, EventClosed}]++
		}
	}

	var reports []AnomalyReport
	for _, e := range entries {
		if e.Event == EventNoChange {
			continue
		}
		check := func(ports []int, event string) {
			for _, p := range ports {
				k := key{e.Host, p, event}
				if counts[k] < minCount {
					reports = append(reports, AnomalyReport{
						Host:       e.Host,
						Port:       p,
						Event:      event,
						OccurredAt: e.ScannedAt,
						Reason:     "rare event: below minimum occurrence threshold",
					})
				}
			}
		}
		check(e.OpenedPorts, EventOpened)
		check(e.ClosedPorts, EventClosed)
	}
	return reports
}

// AnomaliesByHost groups anomaly reports by host.
func AnomaliesByHost(reports []AnomalyReport) map[string][]AnomalyReport {
	out := make(map[string][]AnomalyReport)
	for _, r := range reports {
		out[r.Host] = append(out[r.Host], r)
	}
	return out
}
