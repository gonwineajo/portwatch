package history

import (
	"fmt"
	"sort"
	"strings"
)

// Signature represents a fingerprint of open ports for a host at a point in time.
type Signature struct {
	Host      string
	Ports     []int
	Fingerprint string
}

// BuildSignatures computes a port signature for each host based on their
// most recent scan entry in the provided history.
func BuildSignatures(entries []Entry) []Signature {
	latest := make(map[string]Entry)
	for _, e := range entries {
		if e.Event != EventScan {
			continue
		}
		prev, ok := latest[e.Host]
		if !ok || e.Timestamp.After(prev.Timestamp) {
			latest[e.Host] = e
		}
	}

	var sigs []Signature
	for host, e := range latest {
		ports := make([]int, len(e.Ports))
		copy(ports, e.Ports)
		sort.Ints(ports)
		sigs = append(sigs, Signature{
			Host:        host,
			Ports:       ports,
			Fingerprint: computeFingerprint(ports),
		})
	}

	sort.Slice(sigs, func(i, j int) bool {
		return sigs[i].Host < sigs[j].Host
	})
	return sigs
}

// MatchSignature returns all Signature entries whose fingerprint matches
// the fingerprint of the given target host, excluding the target itself.
func MatchSignature(sigs []Signature, targetHost string) []Signature {
	var target *Signature
	for i := range sigs {
		if sigs[i].Host == targetHost {
			target = &sigs[i]
			break
		}
	}
	if target == nil {
		return nil
	}

	var matches []Signature
	for _, s := range sigs {
		if s.Host != targetHost && s.Fingerprint == target.Fingerprint {
			matches = append(matches, s)
		}
	}
	return matches
}

func computeFingerprint(ports []int) string {
	parts := make([]string, len(ports))
	for i, p := range ports {
		parts[i] = fmt.Sprintf("%d", p)
	}
	return strings.Join(parts, ",")
}
