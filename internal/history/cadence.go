package history

import (
	"sort"
	"time"
)

// CadenceResult describes how regularly a port opens on a given host.
type CadenceResult struct {
	Host        string
	Port        int
	Occurrences int
	AvgInterval time.Duration
	MinInterval time.Duration
	MaxInterval time.Duration
	Regular     bool // true if stddev is within 20% of mean
}

// AnalyseCadence computes the open-event cadence for each (host, port) pair
// found in entries. A pair is considered "regular" when the coefficient of
// variation of its inter-arrival intervals is ≤ 0.20.
func AnalyseCadence(entries []Entry) []CadenceResult {
	type key struct {
		host string
		port int
	}
	times := make(map[key][]time.Time)

	for _, e := range entries {
		if e.Event != "opened" {
			continue
		}
		for _, p := range e.Ports {
			k := key{e.Host, p}
			times[k] = append(times[k], e.Timestamp)
		}
	}

	var results []CadenceResult
	for k, ts := range times {
		if len(ts) < 2 {
			continue
		}
		sort.Slice(ts, func(i, j int) bool { return ts[i].Before(ts[j]) })

		intervals := make([]float64, len(ts)-1)
		var sum float64
		for i := 1; i < len(ts); i++ {
			d := ts[i].Sub(ts[i-1]).Seconds()
			intervals[i-1] = d
			sum += d
		}
		mean := sum / float64(len(intervals))

		var minD, maxD float64 = intervals[0], intervals[0]
		var variance float64
		for _, v := range intervals {
			if v < minD {
				minD = v
			}
			if v > maxD {
				maxD = v
			}
			diff := v - mean
			variance += diff * diff
		}
		variance /= float64(len(intervals))

		var cv float64
		if mean > 0 {
			cv = sqrt64(variance) / mean
		}

		results = append(results, CadenceResult{
			Host:        k.host,
			Port:        k.port,
			Occurrences: len(ts),
			AvgInterval: time.Duration(mean * float64(time.Second)),
			MinInterval: time.Duration(minD * float64(time.Second)),
			MaxInterval: time.Duration(maxD * float64(time.Second)),
			Regular:     cv <= 0.20,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Host != results[j].Host {
			return results[i].Host < results[j].Host
		}
		return results[i].Port < results[j].Port
	})
	return results
}

// sqrt64 is a small helper to avoid importing math in this file.
func sqrt64(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x
	for i := 0; i < 50; i++ {
		z -= (z*z - x) / (2 * z)
	}
	return z
}
