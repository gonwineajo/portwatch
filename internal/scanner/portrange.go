package scanner

import (
	"fmt"
	"strconv"
	"strings"
)

// ParsePortRange parses a port specification string into a slice of ints.
// Accepted formats: "80", "80,443", "8000-8080", or combinations like "22,80,8000-8080".
func ParsePortRange(spec string) ([]int, error) {
	var ports []int
	seen := make(map[int]bool)

	for _, part := range strings.Split(spec, ",") {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			bounds := strings.SplitN(part, "-", 2)
			lo, err1 := strconv.Atoi(bounds[0])
			hi, err2 := strconv.Atoi(bounds[1])
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("invalid range %q", part)
			}
			if lo > hi {
				return nil, fmt.Errorf("range start %d > end %d", lo, hi)
			}
			for p := lo; p <= hi; p++ {
				if !seen[p] {
					ports = append(ports, p)
					seen[p] = true
				}
			}
		} else {
			p, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid port %q", part)
			}
			if !seen[p] {
				ports = append(ports, p)
				seen[p] = true
			}
		}
	}
	return ports, nil
}
