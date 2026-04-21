package history

import "time"

// BurstResult describes a burst of port-change events detected for a host
// within a short window of time.
type BurstResult struct {
	Host      string
	WindowEnd time.Time
	Count     int
	Ports     []int
	Events    []string
}

// DetectBursts scans entries for hosts that produced more than threshold
// port-change events (opened or closed) within the given window duration.
// Only the first qualifying burst window per host is returned.
func DetectBursts(entries []Entry, window time.Duration, threshold int) []BurstResult {
	if len(entries) == 0 || window <= 0 || threshold <= 0 {
		return nil
	}

	// Group change events by host, sorted by time (entries assumed ordered).
	byHost := make(map[string][]Entry)
	for _, e := range entries {
		if e.Event == EventOpened || e.Event == EventClosed {
			byHost[e.Host] = append(byHost[e.Host], e)
		}
	}

	var results []BurstResult

	for host, evs := range byHost {
		// Sliding window: find first window where count >= threshold.
		for i := 0; i < len(evs); i++ {
			start := evs[i].Timestamp
			end := start.Add(window)
			var ports []int
			var events []string

			for j := i; j < len(evs); j++ {
				if evs[j].Timestamp.After(end) {
					break
				}
				ports = append(ports, evs[j].Ports...)
				events = append(events, string(evs[j].Event))
			}

			if len(events) >= threshold {
				results = append(results, BurstResult{
					Host:      host,
					WindowEnd: end,
					Count:     len(events),
					Ports:     ports,
					Events:    events,
				})
				break // one burst result per host
			}
		}
	}

	return results
}
