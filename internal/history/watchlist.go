package history

// WatchlistEntry represents a port+host combination being actively watched.
type WatchlistEntry struct {
	Host string
	Port int
	Label string
}

// Watchlist holds a set of host:port pairs to monitor closely.
type Watchlist struct {
	entries []WatchlistEntry
}

// NewWatchlist creates an empty Watchlist.
func NewWatchlist() *Watchlist {
	return &Watchlist{}
}

// Add registers a host:port pair with an optional label.
func (w *Watchlist) Add(host string, port int, label string) {
	w.entries = append(w.entries, WatchlistEntry{Host: host, Port: port, Label: label})
}

// Remove deletes all entries matching the given host and port.
func (w *Watchlist) Remove(host string, port int) {
	filtered := w.entries[:0]
	for _, e := range w.entries {
		if e.Host != host || e.Port != port {
			filtered = append(filtered, e)
		}
	}
	w.entries = filtered
}

// Entries returns a copy of all watchlist entries.
func (w *Watchlist) Entries() []WatchlistEntry {
	out := make([]WatchlistEntry, len(w.entries))
	copy(out, w.entries)
	return out
}

// Match returns all history entries that involve a watched host:port pair.
func (w *Watchlist) Match(entries []Entry) []Entry {
	type key struct {
		host string
		port int
	}
	set := make(map[key]struct{}, len(w.entries))
	for _, we := range w.entries {
		set[key{we.Host, we.Port}] = struct{}{}
	}

	var out []Entry
	for _, e := range entries {
		for _, p := range e.Ports {
			if _, ok := set[key{e.Host, p}]; ok {
				out = append(out, e)
				break
			}
		}
	}
	return out
}
