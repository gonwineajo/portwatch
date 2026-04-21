package history

import (
	"sort"
	"time"
)

// PortVelocity describes how rapidly a port's open/close state is changing
// for a given host over a sliding window.
type PortVelocity struct {
	Host      string
	Port      int
	Opened    int           // number of opened events in window
	Closed    int           // number of closed events in window
	Flips     int           // total state changes (opened + closed)
	Rate      float64       // flips per hour
	Window    time.Duration // window used for calculation
}

// VelocityOptions controls how velocity is calculated.
type VelocityOptions struct {
	// Window is the time duration to look back from the most recent entry.
	// Defaults to 24 hours if zero.
	Window time.Duration

	// MinFlips filters out ports with fewer state changes than this threshold.
	// Defaults to 1.
	MinFlips int

	// Host restricts calculation to a specific host. Empty means all hosts.
	Host string
}

// Velocity calculates the rate of port state changes per host/port pair
// within the configured time window. Results are sorted by descending rate.
func Velocity(entries []Entry, opts VelocityOptions) []PortVelocity {
	if opts.Window == 0 {
		opts.Window = 24 * time.Hour
	}
	if opts.MinFlips < 1 {
		opts.MinFlips = 1
	}

	// Determine the cutoff time from the latest entry timestamp.
	var latest time.Time
	for _, e := range entries {
		if e.Time.After(latest) {
			latest = e.Time
		}
	}
	if latest.IsZero() {
		return nil
	}
	cutoff := latest.Add(-opts.Window)

	type key struct {
		host string
		port int
	}
	type counts struct {
		opened int
		closed int
	}

	tally := make(map[key]*counts)

	for _, e := range entries {
		if e.Time.Before(cutoff) {
			continue
		}
		if e.Event == EventNoChange {
			continue
		}
		if opts.Host != "" && e.Host != opts.Host {
			continue
		}

		ports := e.OpenedPorts
		event := EventOpened
		if e.Event == EventClosed {
			ports = e.ClosedPorts
			event = EventClosed
		}

		for _, p := range ports {
			k := key{host: e.Host, port: p}
			if tally[k] == nil {
				tally[k] = &counts{}
			}
			if event == EventOpened {
				tally[k].opened++
			} else {
				tally[k].closed++
			}
		}
	}

	hours := opts.Window.Hours()
	if hours == 0 {
		hours = 1
	}

	var results []PortVelocity
	for k, c := range tally {
		flips := c.opened + c.closed
		if flips < opts.MinFlips {
			continue
		}
		results = append(results, PortVelocity{
			Host:   k.host,
			Port:   k.port,
			Opened: c.opened,
			Closed: c.closed,
			Flips:  flips,
			Rate:   float64(flips) / hours,
			Window: opts.Window,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Rate != results[j].Rate {
			return results[i].Rate > results[j].Rate
		}
		if results[i].Host != results[j].Host {
			return results[i].Host < results[j].Host
		}
		return results[i].Port < results[j].Port
	})

	return results
}
