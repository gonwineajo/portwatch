package history

import "time"

// FilterOptions defines criteria for filtering history entries.
type FilterOptions struct {
	Host    string
	Since   time.Time
	Until   time.Time
	Event   string // "opened" or "closed"
	Limit   int
}

// Filter returns entries matching all non-zero criteria in opts.
func Filter(entries []Entry, opts FilterOptions) []Entry {
	var out []Entry
	for _, e := range entries {
		if opts.Host != "" && e.Host != opts.Host {
			continue
		}
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		if !opts.Until.IsZero() && e.Timestamp.After(opts.Until) {
			continue
		}
		if opts.Event != "" && e.Event != opts.Event {
			continue
		}
		out = append(out, e)
		if opts.Limit > 0 && len(out) >= opts.Limit {
			break
		}
	}
	return out
}
