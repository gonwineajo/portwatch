package history

import "sort"

// RiskLevel categorises the threat level of a host based on port activity.
type RiskLevel string

const (
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"
)

// RiskReport summarises the risk profile of a single host.
type RiskReport struct {
	Host        string
	Level       RiskLevel
	OpenedCount int
	ClosedCount int
	UniquePorts int
	Score       float64
}

// AssessRisk evaluates each host's risk level based on port-change frequency
// and the number of distinct ports involved.
func AssessRisk(entries []Entry) []RiskReport {
	type hostStats struct {
		opened int
		closed int
		ports  map[int]struct{}
	}

	hosts := map[string]*hostStats{}

	for _, e := range entries {
		if e.Event != EventOpened && e.Event != EventClosed {
			continue
		}
		s, ok := hosts[e.Host]
		if !ok {
			s = &hostStats{ports: map[int]struct{}{}}
			hosts[e.Host] = s
		}
		if e.Event == EventOpened {
			s.opened++
		} else {
			s.closed++
		}
		for _, p := range e.Ports {
			s.ports[p] = struct{}{}
		}
	}

	var reports []RiskReport
	for host, s := range hosts {
		score := float64(s.opened*2+s.closed) + float64(len(s.ports))*0.5
		var level RiskLevel
		switch {
		case score >= 10:
			level = RiskHigh
		case score >= 4:
			level = RiskMedium
		default:
			level = RiskLow
		}
		reports = append(reports, RiskReport{
			Host:        host,
			Level:       level,
			OpenedCount: s.opened,
			ClosedCount: s.closed,
			UniquePorts: len(s.ports),
			Score:       score,
		})
	}

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].Score > reports[j].Score
	})
	return reports
}
