package history

import "strings"

// SearchOptions filters history entries.
type SearchOptions struct {
	Host    string // substring match on host
	Port    int    // 0 means any
	Event   string // "opened", "closed", or "" for any
}

// Search returns entries matching all non-zero criteria in opts.
func Search(entries []Entry, opts SearchOptions) []Entry {
	var out []Entry
	for _, e := range entries {
		if opts.Host != "" && !strings.Contains(e.Host, opts.Host) {
			continue
		}
		if opts.Event != "" && e.Event != opts.Event {
			continue
		}
		if opts.Port != 0 && !containsPort(e, opts.Port) {
			continue
		}
		out = append(out, e)
	}
	return out
}

// SearchByHost returns all entries whose host contains the given substring.
// It is a convenience wrapper around Search for host-only queries.
func SearchByHost(entries []Entry, host string) []Entry {
	return Search(entries, SearchOptions{Host: host})
}

func containsPort(e Entry, port int) bool {
	for _, p := range e.Ports {
		if p == port {
			return true
		}
	}
	return false
}
