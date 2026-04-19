package history

import "time"

// Bucket holds aggregated event counts for a time window.
type Bucket struct {
	Start  time.Time
	Opened int
	Closed int
}

// Timeline groups history entries into fixed-duration buckets.
func Timeline(entries []Entry, window time.Duration) []Bucket {
	if len(entries) == 0 || window <= 0 {
		return nil
	}

	start := entries[0].Timestamp.Truncate(window)
	end := entries[len(entries)-1].Timestamp

	var buckets []Bucket
	for t := start; !t.After(end); t = t.Add(window) {
		buckets = append(buckets, Bucket{Start: t})
	}

	for _, e := range entries {
		idx := int(e.Timestamp.Truncate(window).Sub(start) / window)
		if idx < 0 || idx >= len(buckets) {
			continue
		}
		buckets[idx].Opened += len(e.Opened)
		buckets[idx].Closed += len(e.Closed)
	}

	return buckets
}
