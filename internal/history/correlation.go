package history

import "sort"

// PortCorrelation describes two ports that frequently appear together
// in the same scan event across one or more hosts.
type PortCorrelation struct {
	PortA     int
	PortB     int
	CoOccurrences int
	Hosts     []string
}

// CorrelateOpenPorts finds pairs of ports that are frequently opened
// together within the same scan entry. Only pairs seen at least
// minCount times are returned, sorted by co-occurrence count descending.
func CorrelateOpenPorts(entries []Entry, minCount int) []PortCorrelation {
	type pairKey struct{ a, b int }

	counts := map[pairKey]int{}
	hostSets := map[pairKey]map[string]struct{}{}

	for _, e := range entries {
		if e.Event != EventOpened && e.Event != EventScan {
			continue
		}
		ports := e.Ports
		for i := 0; i < len(ports); i++ {
			for j := i + 1; j < len(ports); j++ {
				a, b := ports[i], ports[j]
				if a > b {
					a, b = b, a
				}
				key := pairKey{a, b}
				counts[key]++
				if hostSets[key] == nil {
					hostSets[key] = map[string]struct{}{}
				}
				hostSets[key][e.Host] = struct{}{}
			}
		}
	}

	var result []PortCorrelation
	for key, count := range counts {
		if count < minCount {
			continue
		}
		hosts := make([]string, 0, len(hostSets[key]))
		for h := range hostSets[key] {
			hosts = append(hosts, h)
		}
		sort.Strings(hosts)
		result = append(result, PortCorrelation{
			PortA:         key.a,
			PortB:         key.b,
			CoOccurrences: count,
			Hosts:         hosts,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].CoOccurrences != result[j].CoOccurrences {
			return result[i].CoOccurrences > result[j].CoOccurrences
		}
		if result[i].PortA != result[j].PortA {
			return result[i].PortA < result[j].PortA
		}
		return result[i].PortB < result[j].PortB
	})

	return result
}
