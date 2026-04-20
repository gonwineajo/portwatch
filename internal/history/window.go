package history

import "time"

// WindowStats holds aggregated port event counts over a rolling time window.
type WindowStats struct {
	Host    string
	Opened  int
	Closed  int
	Total   int
	Since   time.Time
	Until   time.Time
}

// RollingWindow returns per-host stats for entries within the given duration
// looking back from `now`. Pass time.Now() for live use.
func RollingWindow(entries []Entry, window time.Duration, now time.Time) []WindowStats {
	since := now.Add(-window)

	type counts struct {
		opened int
		closed int
	}

	hostMap := make(map[string]*counts)

	for _, e := range entries {
		if e.Timestamp.Before(since) || e.Timestamp.After(now) {
			continue
		}
		c, ok := hostMap[e.Host]
		if !ok {
			c = &counts{}
			hostMap[e.Host] = c
		}
		switch e.Event {
		case "opened":
			c.opened++
		case "closed":
			c.closed++
		}
	}

	results := make([]WindowStats, 0, len(hostMap))
	for host, c := range hostMap {
		results = append(results, WindowStats{
			Host:   host,
			Opened: c.opened,
			Closed: c.closed,
			Total:  c.opened + c.closed,
			Since:  since,
			Until:  now,
		})
	}

	// Sort by Total descending, then Host ascending for stable output.
	for i := 1; i < len(results); i++ {
		for j := i; j > 0; j-- {
			a, b := results[j-1], results[j]
			if a.Total < b.Total || (a.Total == b.Total && a.Host > b.Host) {
				results[j-1], results[j] = b, a
			} else {
				break
			}
		}
	}

	return results
}
