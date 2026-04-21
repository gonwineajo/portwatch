package history

import (
	"sort"
	"time"
)

// PortForecast represents a predicted port event for a host.
type PortForecast struct {
	Host      string
	Port      int
	Event     EventType
	LikelyAt  time.Time
	Confidence float64 // 0.0 to 1.0
}

// Forecast predicts future port events based on recurring patterns in entries.
// It uses average interval between past occurrences to project the next event time.
func Forecast(entries []Entry, now time.Time) []PortForecast {
	patterns := DetectPatterns(entries)
	recurring := RecurringOnly(patterns)

	var forecasts []PortForecast

	for _, p := range recurring {
		if len(p.SeenAt) < 2 {
			continue
		}

		times := make([]time.Time, len(p.SeenAt))
		copy(times, p.SeenAt)
		sort.Slice(times, func(i, j int) bool {
			return times[i].Before(times[j])
		})

		var totalInterval time.Duration
		for i := 1; i < len(times); i++ {
			totalInterval += times[i].Sub(times[i-1])
		}
		avgInterval := totalInterval / time.Duration(len(times)-1)

		last := times[len(times)-1]
		nextTime := last.Add(avgInterval)

		// Confidence increases with more observations, capped at 0.95
		confidence := float64(p.Count) / float64(p.Count+2)
		if confidence > 0.95 {
			confidence = 0.95
		}

		forecasts = append(forecasts, PortForecast{
			Host:       p.Host,
			Port:       p.Port,
			Event:      p.Event,
			LikelyAt:   nextTime,
			Confidence: confidence,
		})
	}

	sort.Slice(forecasts, func(i, j int) bool {
		return forecasts[i].LikelyAt.Before(forecasts[j].LikelyAt)
	})

	return forecasts
}

// ForecastByHost returns forecasts filtered to a specific host.
func ForecastByHost(entries []Entry, host string, now time.Time) []PortForecast {
	all := Forecast(entries, now)
	var out []PortForecast
	for _, f := range all {
		if f.Host == host {
			out = append(out, f)
		}
	}
	return out
}
