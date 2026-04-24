package history

// OverlapResult describes the port overlap between two hosts.
type OverlapResult struct {
	HostA      string
	HostB      string
	Shared     []int
	OnlyA      []int
	OnlyB      []int
	JaccardSim float64
}

// AnalyseOverlap computes the port overlap between every pair of hosts based
// on their latest scan entries. Hosts with no scan entries are skipped.
func AnalyseOverlap(entries []Entry) []OverlapResult {
	portsByHost := latestScanPortsMap(entries)

	hosts := make([]string, 0, len(portsByHost))
	for h := range portsByHost {
		hosts = append(hosts, h)
	}
	sortStrings(hosts)

	var results []OverlapResult
	for i := 0; i < len(hosts); i++ {
		for j := i + 1; j < len(hosts); j++ {
			ha, hb := hosts[i], hosts[j]
			setA := portsByHost[ha]
			setB := portsByHost[hb]

			shared, onlyA, onlyB := diffSets(setA, setB)

			unionSize := len(shared) + len(onlyA) + len(onlyB)
			var jaccard float64
			if unionSize > 0 {
				jaccard = float64(len(shared)) / float64(unionSize)
			}

			results = append(results, OverlapResult{
				HostA:      ha,
				HostB:      hb,
				Shared:     shared,
				OnlyA:      onlyA,
				OnlyB:      onlyB,
				JaccardSim: jaccard,
			})
		}
	}
	return results
}

// latestScanPortsMap returns a map of host -> set of open ports from the most
// recent scan entry per host.
func latestScanPortsMap(entries []Entry) map[string]map[int]struct{} {
	latest := map[string]Entry{}
	for _, e := range entries {
		if e.Event != EventScan {
			continue
		}
		if prev, ok := latest[e.Host]; !ok || e.Timestamp.After(prev.Timestamp) {
			latest[e.Host] = e
		}
	}
	out := make(map[string]map[int]struct{}, len(latest))
	for h, e := range latest {
		s := make(map[int]struct{}, len(e.Ports))
		for _, p := range e.Ports {
			s[p] = struct{}{}
		}
		out[h] = s
	}
	return out
}

// diffSets returns shared ports, ports only in a, and ports only in b — all sorted.
func diffSets(a, b map[int]struct{}) (shared, onlyA, onlyB []int) {
	for p := range a {
		if _, ok := b[p]; ok {
			shared = append(shared, p)
		} else {
			onlyA = append(onlyA, p)
		}
	}
	for p := range b {
		if _, ok := a[p]; !ok {
			onlyB = append(onlyB, p)
		}
	}
	sortInts(shared)
	sortInts(onlyA)
	sortInts(onlyB)
	return
}

func sortStrings(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}

func sortInts(is []int) {
	for i := 1; i < len(is); i++ {
		for j := i; j > 0 && is[j] < is[j-1]; j-- {
			is[j], is[j-1] = is[j-1], is[j]
		}
	}
}
