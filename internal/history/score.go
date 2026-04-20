package history

import "sort"

// RiskScore represents a computed risk level for a host based on its history.
type RiskScore struct {
	Host      string
	Score     int
	OpenCount int
	CloseCount int
	UniquePortsOpened int
}

// ScoreHosts computes a simple risk score for each host in the provided entries.
// Score is weighted: each opened port event adds 3 points, each closed port
// event adds 1 point, and unique ports opened contribute an additional 2 points each.
func ScoreHosts(entries []Entry) []RiskScore {
	type hostData struct {
		openCount  int
		closeCount int
		ports      map[int]struct{}
	}

	hosts := make(map[string]*hostData)

	for _, e := range entries {
		if _, ok := hosts[e.Host]; !ok {
			hosts[e.Host] = &hostData{ports: make(map[int]struct{})}
		}
		hd := hosts[e.Host]
		switch e.Event {
		case EventOpened:
			hd.openCount++
			for _, p := range e.Ports {
				hd.ports[p] = struct{}{}
			}
		case EventClosed:
			hd.closeCount++
		}
	}

	results := make([]RiskScore, 0, len(hosts))
	for host, hd := range hosts {
		score := hd.openCount*3 + hd.closeCount*1 + len(hd.ports)*2
		results = append(results, RiskScore{
			Host:              host,
			Score:             score,
			OpenCount:         hd.openCount,
			CloseCount:        hd.closeCount,
			UniquePortsOpened: len(hd.ports),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		return results[i].Host < results[j].Host
	})

	return results
}
