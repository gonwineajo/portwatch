package history

import "time"

// PortPressure represents the cumulative "pressure" score for a port on a host,
// reflecting how frequently it transitions between open and closed states.
type PortPressure struct {
	Host      string
	Port      int
	Flips     int
	FirstSeen time.Time
	LastSeen  time.Time
	Score     float64 // flips per hour over the observed window
}

// PressureResult holds all pressure records, sorted by descending Score.
type PressureResult struct {
	Records []PortPressure
}

// AnalysePressure computes port-flip pressure for each (host, port) pair.
// Only entries with EventType "opened" or "closed" are considered.
// minFlips sets the minimum number of transitions required to be included.
func AnalysePressure(entries []Entry, minFlips int) PressureResult {
	type key struct {
		host string
		port int
	}

	type record struct {
		flips     int
		first     time.Time
		last      time.Time
	}

	data := make(map[key]*record)

	for _, e := range entries {
		if e.EventType != EventOpened && e.EventType != EventClosed {
			continue
		}
		for _, p := range e.Ports {
			k := key{host: e.Host, port: p}
			r, ok := data[k]
			if !ok {
				r = &record{first: e.Timestamp, last: e.Timestamp}
				data[k] = r
			}
			r.flips++
			if e.Timestamp.Before(r.first) {
				r.first = e.Timestamp
			}
			if e.Timestamp.After(r.last) {
				r.last = e.Timestamp
			}
		}
	}

	var result []PortPressure
	for k, r := range data {
		if r.flips < minFlips {
			continue
		}
		hours := r.last.Sub(r.first).Hours()
		var score float64
		if hours > 0 {
			score = float64(r.flips) / hours
		} else {
			score = float64(r.flips)
		}
		result = append(result, PortPressure{
			Host:      k.host,
			Port:      k.port,
			Flips:     r.flips,
			FirstSeen: r.first,
			LastSeen:  r.last,
			Score:     score,
		})
	}

	sortPortPressure(result)
	return PressureResult{Records: result}
}

func sortPortPressure(recs []PortPressure) {
	for i := 1; i < len(recs); i++ {
		for j := i; j > 0 && recs[j].Score > recs[j-1].Score; j-- {
			recs[j], recs[j-1] = recs[j-1], recs[j]
		}
	}
}
