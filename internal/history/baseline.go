package history

import "time"

// Baseline represents the known-good port state for a host at a point in time.
type Baseline struct {
	Host      string    `json:"host"`
	Ports     []int     `json:"ports"`
	CreatedAt time.Time `json:"created_at"`
	Note      string    `json:"note,omitempty"`
}

// SetBaseline creates a Baseline from the most recent scan entry for each host
// in the provided entries slice. Only entries with EventScan are considered.
func SetBaseline(entries []Entry) []Baseline {
	latest := make(map[string]Entry)
	for _, e := range entries {
		if e.Event != EventScan {
			continue
		}
		prev, ok := latest[e.Host]
		if !ok || e.Timestamp.After(prev.Timestamp) {
			latest[e.Host] = e
		}
	}

	baselines := make([]Baseline, 0, len(latest))
	for host, e := range latest {
		baselines = append(baselines, Baseline{
			Host:      host,
			Ports:     e.Ports,
			CreatedAt: e.Timestamp,
		})
	}
	return baselines
}

// DeviatesFromBaseline returns ports that are open in current but not in the
// baseline, and ports that are in the baseline but closed in current.
func DeviatesFromBaseline(baseline Baseline, current []int) (opened []int, closed []int) {
	baseSet := make(map[int]bool, len(baseline.Ports))
	for _, p := range baseline.Ports {
		baseSet[p] = true
	}
	currSet := make(map[int]bool, len(current))
	for _, p := range current {
		currSet[p] = true
	}
	for p := range currSet {
		if !baseSet[p] {
			opened = append(opened, p)
		}
	}
	for p := range baseSet {
		if !currSet[p] {
			closed = append(closed, p)
		}
	}
	return opened, closed
}
