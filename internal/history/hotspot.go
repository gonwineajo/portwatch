package history

import "sort"

// Hotspot represents a host+port pair that has experienced frequent open/close
// transitions, indicating instability or repeated activity.
type Hotspot struct {
	Host   string
	Port   int
	Flips  int     // number of open<->close transitions
	Score  float64 // weighted activity score
}

// DetectHotspots returns host+port pairs ranked by transition frequency.
// minFlips sets the minimum number of open/close transitions to be included.
func DetectHotspots(entries []Entry, minFlips int) []Hotspot {
	type key struct {
		host string
		port int
	}

	flips := make(map[key]int)
	last := make(map[key]string) // last event type seen for this pair

	for _, e := range entries {
		if e.Event == EventNoChange {
			continue
		}
		for _, p := range e.Ports {
			k := key{host: e.Host, port: p}
			prev, seen := last[k]
			if seen && prev != e.Event {
				flips[k]++
			}
			last[k] = e.Event
		}
	}

	var result []Hotspot
	for k, f := range flips {
		if f < minFlips {
			continue
		}
		result = append(result, Hotspot{
			Host:  k.host,
			Port:  k.port,
			Flips: f,
			Score: float64(f),
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Score != result[j].Score {
			return result[i].Score > result[j].Score
		}
		if result[i].Host != result[j].Host {
			return result[i].Host < result[j].Host
		}
		return result[i].Port < result[j].Port
	})

	return result
}
