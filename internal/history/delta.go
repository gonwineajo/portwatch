package history

import "sort"

// PortDelta holds aggregated open/close counts per port across all entries.
type PortDelta struct {
	Port       int
	TimesOpened int
	TimesClosed int
}

// PortDeltas computes aggregated open/close counts per port from a slice of entries.
func PortDeltas(entries []Entry) []PortDelta {
	type counts struct{ opened, closed int }
	m := map[int]*counts{}

	for _, e := range entries {
		for _, p := range e.Opened {
			if m[p] == nil {
				m[p] = &counts{}
			}
			m[p].opened++
		}
		for _, p := range e.Closed {
			if m[p] == nil {
				m[p] = &counts{}
			}
			m[p].closed++
		}
	}

	result := make([]PortDelta, 0, len(m))
	for port, c := range m {
		result = append(result, PortDelta{Port: port, TimesOpened: c.opened, TimesClosed: c.closed})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Port < result[j].Port })
	return result
}
