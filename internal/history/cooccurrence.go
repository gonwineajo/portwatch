package history

import "sort"

// CoOccurrence describes how often two ports appear open together across scans.
type CoOccurrence struct {
	PortA int
	PortB int
	Count int
	Hosts []string
}

// CoOccurrenceOptions controls filtering for CoOccurrenceMatrix.
type CoOccurrenceOptions struct {
	Host     string // filter to a single host; empty = all hosts
	MinCount int    // minimum co-occurrence count to include
}

// CoOccurrenceMatrix returns pairs of ports that appear open together in the
// same scan entry, ranked by co-occurrence count descending.
func CoOccurrenceMatrix(entries []Entry, opts CoOccurrenceOptions) []CoOccurrence {
	type pairKey struct{ a, b int }

	counts := make(map[pairKey]int)
	hostSets := make(map[pairKey]map[string]struct{})

	for _, e := range entries {
		if e.Event != EventScan {
			continue
		}
		if opts.Host != "" && e.Host != opts.Host {
			continue
		}
		ports := e.Ports
		for i := 0; i < len(ports); i++ {
			for j := i + 1; j < len(ports); j++ {
				a, b := ports[i], ports[j]
				if a > b {
					a, b = b, a
				}
				k := pairKey{a, b}
				counts[k]++
				if hostSets[k] == nil {
					hostSets[k] = make(map[string]struct{})
				}
				hostSets[k][e.Host] = struct{}{}
			}
		}
	}

	minCount := opts.MinCount
	if minCount <= 0 {
		minCount = 1
	}

	var result []CoOccurrence
	for k, count := range counts {
		if count < minCount {
			continue
		}
		hosts := make([]string, 0, len(hostSets[k]))
		for h := range hostSets[k] {
			hosts = append(hosts, h)
		}
		sort.Strings(hosts)
		result = append(result, CoOccurrence{
			PortA: k.a,
			PortB: k.b,
			Count: count,
			Hosts: hosts,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Count != result[j].Count {
			return result[i].Count > result[j].Count
		}
		if result[i].PortA != result[j].PortA {
			return result[i].PortA < result[j].PortA
		}
		return result[i].PortB < result[j].PortB
	})
	return result
}
