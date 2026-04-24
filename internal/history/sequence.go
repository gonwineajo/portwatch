package history

import "sort"

// PortSequence describes the order in which ports were opened on a host.
type PortSequence struct {
	Host  string
	Ports []int // ordered by first-seen time
}

// BuildSequences returns the chronological open-port sequence for each host.
// Only "opened" events are considered; the first occurrence of each port
// determines its position in the sequence.
func BuildSequences(entries []Entry) []PortSequence {
	type hostData struct {
		seen  map[int]bool
		order []int
	}

	hosts := map[string]*hostData{}

	for _, e := range entries {
		if e.Event != EventOpened {
			continue
		}
		hd, ok := hosts[e.Host]
		if !ok {
			hd = &hostData{seen: map[int]bool{}}
			hosts[e.Host] = hd
		}
		for _, p := range e.Ports {
			if !hd.seen[p] {
				hd.seen[p] = true
				hd.order = append(hd.order, p)
			}
		}
	}

	keys := make([]string, 0, len(hosts))
	for k := range hosts {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make([]PortSequence, 0, len(keys))
	for _, k := range keys {
		out = append(out, PortSequence{
			Host:  k,
			Ports: hosts[k].order,
		})
	}
	return out
}

// SequenceForHost returns the open-port sequence for a single host.
// Returns nil if no opened events exist for that host.
func SequenceForHost(entries []Entry, host string) *PortSequence {
	seqs := BuildSequences(entries)
	for _, s := range seqs {
		if s.Host == host {
			return &s
		}
	}
	return nil
}
