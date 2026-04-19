package history

import "time"

// ReplayOptions controls what events are replayed.
type ReplayOptions struct {
	Host  string
	Since time.Time
	Until time.Time
	Limit int
}

// ReplayEvent represents a single replayed port event.
type ReplayEvent struct {
	Timestamp time.Time
	Host      string
	Event     string
	Port      int
}

// Replay returns a flat list of individual port events from history entries,
// filtered by the provided options.
func Replay(entries []Entry, opts ReplayOptions) []ReplayEvent {
	var events []ReplayEvent

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
		for _, p := range e.Opened {
			events = append(events, ReplayEvent{Timestamp: e.Timestamp, Host: e.Host, Event: "opened", Port: p})
		}
		for _, p := range e.Closed {
			events = append(events, ReplayEvent{Timestamp: e.Timestamp, Host: e.Host, Event: "closed", Port: p})
		}
	}

	if opts.Limit > 0 && len(events) > opts.Limit {
		events = events[:opts.Limit]
	}
	return events
}
