package history

import "sort"

// FingerprintDiff compares two host signatures and returns which ports were
// added or removed between them.
type FingerprintDiff struct {
	Host    string
	Added   []int
	Removed []int
}

// DiffSignatures compares current signatures against a baseline set and
// returns per-host diffs for any host whose fingerprint has changed.
func DiffSignatures(baseline, current []HostSignature) []FingerprintDiff {
	baseMap := make(map[string]map[int]struct{})
	for _, s := range baseline {
		set := make(map[int]struct{}, len(s.Ports))
		for _, p := range s.Ports {
			set[p] = struct{}{}
		}
		baseMap[s.Host] = set
	}

	var diffs []FingerprintDiff
	for _, s := range current {
		base, ok := baseMap[s.Host]
		if !ok {
			// brand-new host — all ports are added
			added := make([]int, len(s.Ports))
			copy(added, s.Ports)
			sort.Ints(added)
			diffs = append(diffs, FingerprintDiff{Host: s.Host, Added: added})
			continue
		}

		curSet := make(map[int]struct{}, len(s.Ports))
		for _, p := range s.Ports {
			curSet[p] = struct{}{}
		}

		var added, removed []int
		for p := range curSet {
			if _, exists := base[p]; !exists {
				added = append(added, p)
			}
		}
		for p := range base {
			if _, exists := curSet[p]; !exists {
				removed = append(removed, p)
			}
		}

		if len(added) == 0 && len(removed) == 0 {
			continue
		}
		sort.Ints(added)
		sort.Ints(removed)
		diffs = append(diffs, FingerprintDiff{Host: s.Host, Added: added, Removed: removed})
	}

	sort.Slice(diffs, func(i, j int) bool { return diffs[i].Host < diffs[j].Host })
	return diffs
}
