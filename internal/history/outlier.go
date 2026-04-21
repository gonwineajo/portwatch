package history

import "sort"

// OutlierResult holds a host and its outlier score relative to peers.
type OutlierResult struct {
	Host       string
	OpenPorts  int
	MeanPeers  float64
	Deviation  float64
	IsOutlier  bool
}

// DetectOutliers compares each host's open port count from the latest scan
// against the mean of all hosts. Hosts whose count deviates more than
// threshold standard deviations from the mean are flagged as outliers.
func DetectOutliers(entries []Entry, threshold float64) []OutlierResult {
	if threshold <= 0 {
		threshold = 2.0
	}

	// Collect latest scan per host.
	latest := map[string]Entry{}
	for _, e := range entries {
		if e.Event != EventScan {
			continue
		}
		prev, ok := latest[e.Host]
		if !ok || e.Timestamp.After(prev.Timestamp) {
			latest[e.Host] = e
		}
	}

	if len(latest) == 0 {
		return nil
	}

	// Build counts slice.
	type hostCount struct {
		host  string
		count int
	}
	counts := make([]hostCount, 0, len(latest))
	for host, e := range latest {
		counts = append(counts, hostCount{host, len(e.Ports)})
	}

	// Compute mean.
	sum := 0
	for _, hc := range counts {
		sum += hc.count
	}
	mean := float64(sum) / float64(len(counts))

	// Compute std deviation.
	var variance float64
	for _, hc := range counts {
		d := float64(hc.count) - mean
		variance += d * d
	}
	variance /= float64(len(counts))
	stddev := math.Sqrt(variance)

	results := make([]OutlierResult, 0, len(counts))
	for _, hc := range counts {
		dev := float64(hc.count) - mean
		if stddev < 0.0001 {
			dev = 0
		}
		abs := dev
		if abs < 0 {
			abs = -abs
		}
		results = append(results, OutlierResult{
			Host:      hc.host,
			OpenPorts: hc.count,
			MeanPeers: mean,
			Deviation: dev,
			IsOutlier: stddev >= 0.0001 && abs/stddev >= threshold,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		di, dj := results[i].Deviation, results[j].Deviation
		if di < 0 {
			di = -di
		}
		if dj < 0 {
			dj = -dj
		}
		return di > dj
	})
	return results
}
