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

func containsPort(e Entry, port int) bool {
	for _, p := range e.Ports {
		if p == port {
			return true
		}
	}
	return false
}
