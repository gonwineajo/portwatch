package history

import "sort"

// HeatCell represents activity intensity for a host at a given hour-of-day.
type HeatCell struct {
	Host    string
	Hour    int // 0-23
	Changes int
}

// Heatmap builds a per-host, per-hour-of-day activity map from history entries.
// It counts the number of port change events (opened or closed) per hour bucket.
func Heatmap(entries []Entry) []HeatCell {
	type key struct {
		host string
		hour int
	}

	counts := make(map[key]int)

	for _, e := range entries {
		changes := len(e.OpenedPorts) + len(e.ClosedPorts)
		if changes == 0 {
			continue
		}
		hour := e.Timestamp.Hour()
		counts[key{e.Host, hour}]++
	}

	var cells []HeatCell
	for k, c := range counts {
		cells = append(cells, HeatCell{
			Host:    k.host,
			Hour:    k.hour,
			Changes: c,
		})
	}

	sort.Slice(cells, func(i, j int) bool {
		if cells[i].Host != cells[j].Host {
			return cells[i].Host < cells[j].Host
		}
		return cells[i].Hour < cells[j].Hour
	})

	return cells
}

// PeakHour returns the hour of day with the most total change activity across all hosts.
// Returns -1 if there are no entries with changes.
func PeakHour(entries []Entry) int {
	hourTotals := make(map[int]int)
	for _, e := range entries {
		changes := len(e.OpenedPorts) + len(e.ClosedPorts)
		if changes == 0 {
			continue
		}
		hourTotals[e.Timestamp.Hour()] += changes
	}

	peak, max := -1, 0
	for h, c := range hourTotals {
		if c > max {
			max = c
			peak = h
		}
	}
	return peak
}
