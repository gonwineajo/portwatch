package history

import "sort"

// ClusterResult groups hosts that share a common set of open ports.
type ClusterResult struct {
	Ports []int
	Hosts []string
}

// ClusterByPorts groups hosts from the most recent scan entries by their
// open port fingerprint. Hosts sharing the exact same set of open ports
// are placed in the same cluster. Only "scan" events are considered.
func ClusterByPorts(entries []Entry) []ClusterResult {
	// latest scan per host
	latest := map[string]Entry{}
	for _, e := range entries {
		if e.Event != "scan" {
			continue
		}
		prev, ok := latest[e.Host]
		if !ok || e.Timestamp.After(prev.Timestamp) {
			latest[e.Host] = e
		}
	}

	// group hosts by port fingerprint
	type key = string
	groups := map[key][]string{}
	portSets := map[key][]int{}

	for host, e := range latest {
		sorted := make([]int, len(e.Ports))
		copy(sorted, e.Ports)
		sort.Ints(sorted)
		fp := portFingerprint(sorted)
		groups[fp] = append(groups[fp], host)
		portSets[fp] = sorted
	}

	results := make([]ClusterResult, 0, len(groups))
	for fp, hosts := range groups {
		sort.Strings(hosts)
		results = append(results, ClusterResult{
			Ports: portSets[fp],
			Hosts: hosts,
		})
	}

	// stable sort: larger clusters first, then by first host name
	sort.Slice(results, func(i, j int) bool {
		if len(results[i].Hosts) != len(results[j].Hosts) {
			return len(results[i].Hosts) > len(results[j].Hosts)
		}
		return results[i].Hosts[0] < results[j].Hosts[0]
	})

	return results
}
