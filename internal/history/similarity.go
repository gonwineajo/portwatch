package history

import "sort"

// SimilarityResult holds the similarity score between two hosts.
type SimilarityResult struct {
	HostA      string
	HostB      string
	Common     []int
	OnlyA      []int
	OnlyB      []int
	Jaccard    float64
}

// ComputeSimilarity compares the open port sets of two hosts based on their
// latest scan entries and returns a SimilarityResult with a Jaccard index.
func ComputeSimilarity(entries []Entry, hostA, hostB string) SimilarityResult {
	portsA := latestScanPorts(entries, hostA)
	portsB := latestScanPorts(entries, hostB)

	setA := toIntSet(portsA)
	setB := toIntSet(portsB)

	var common, onlyA, onlyB []int

	for p := range setA {
		if setB[p] {
			common = append(common, p)
		} else {
			onlyA = append(onlyA, p)
		}
	}
	for p := range setB {
		if !setA[p] {
			onlyB = append(onlyB, p)
		}
	}

	sort.Ints(common)
	sort.Ints(onlyA)
	sort.Ints(onlyB)

	union := len(common) + len(onlyA) + len(onlyB)
	var jaccard float64
	if union > 0 {
		jaccard = float64(len(common)) / float64(union)
	}

	return SimilarityResult{
		HostA:   hostA,
		HostB:   hostB,
		Common:  common,
		OnlyA:   onlyA,
		OnlyB:   onlyB,
		Jaccard: jaccard,
	}
}

// AllPairSimilarity computes similarity for every unique pair of hosts found
// in entries, returning results sorted by Jaccard score descending.
func AllPairSimilarity(entries []Entry) []SimilarityResult {
	hostSet := map[string]struct{}{}
	for _, e := range entries {
		if e.Event == EventScan {
			hostSet[e.Host] = struct{}{}
		}
	}
	hosts := make([]string, 0, len(hostSet))
	for h := range hostSet {
		hosts = append(hosts, h)
	}
	sort.Strings(hosts)

	var results []SimilarityResult
	for i := 0; i < len(hosts); i++ {
		for j := i + 1; j < len(hosts); j++ {
			results = append(results, ComputeSimilarity(entries, hosts[i], hosts[j]))
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Jaccard > results[j].Jaccard
	})
	return results
}

func latestScanPorts(entries []Entry, host string) []int {
	for i := len(entries) - 1; i >= 0; i-- {
		e := entries[i]
		if e.Host == host && e.Event == EventScan {
			return e.Ports
		}
	}
	return nil
}

func toIntSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
