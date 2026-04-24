package history

import "sort"

// ReachResult summarises how many distinct hosts a port has been observed
// open on, giving a sense of how "far" that port has reached across the
// monitored fleet.
type ReachResult struct {
	Port  int
	Hosts []string // sorted list of distinct hosts
	Count int      // len(Hosts)
}

// AnalyseReach returns, for every port that has ever been seen in an "opened"
// or "scan" event, the set of distinct hosts it appeared on.  Results are
// ordered by Count descending (highest reach first), then by port number.
//
// Only entries whose EventType is EventOpened or EventScan are considered;
// EventClosed and EventNoChange are ignored so that a port that was briefly
// open on many hosts still shows up, but a port that is merely absent does
// not inflate the numbers.
func AnalyseReach(entries []Entry) []ReachResult {
	type key struct {
		port int
		host string
	}

	seen := make(map[key]struct{})
	portHosts := make(map[int]map[string]struct{})

	for _, e := range entries {
		if e.EventType != EventOpened && e.EventType != EventScan {
			continue
		}
		for _, p := range e.Ports {
			k := key{port: p, host: e.Host}
			if _, exists := seen[k]; exists {
				continue
			}
			seen[k] = struct{}{}
			if portHosts[p] == nil {
				portHosts[p] = make(map[string]struct{})
			}
			portHosts[p][e.Host] = struct{}{}
		}
	}

	results := make([]ReachResult, 0, len(portHosts))
	for port, hostSet := range portHosts {
		hosts := make([]string, 0, len(hostSet))
		for h := range hostSet {
			hosts = append(hosts, h)
		}
		sort.Strings(hosts)
		results = append(results, ReachResult{
			Port:  port,
			Hosts: hosts,
			Count: len(hosts),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Count != results[j].Count {
			return results[i].Count > results[j].Count
		}
		return results[i].Port < results[j].Port
	})

	return results
}

// ReachForPort is a convenience wrapper that returns the ReachResult for a
// single port, and a boolean indicating whether any data was found.
func ReachForPort(entries []Entry, port int) (ReachResult, bool) {
	for _, r := range AnalyseReach(entries) {
		if r.Port == port {
			return r, true
		}
	}
	return ReachResult{}, false
}
