package history

import (
	"math"
	"time"
)

// DecayResult holds the computed relevance score for a host based on
// how recently and frequently port-change events occurred.
type DecayResult struct {
	Host      string
	Score     float64 // higher = more recently/frequently active
	LastEvent time.Time
	EventCount int
}

// DecayScores computes an exponential-decay relevance score for each host.
// Events closer to `now` contribute more to the score than older ones.
// halfLife controls how quickly past events lose relevance.
func DecayScores(entries []Entry, now time.Time, halfLife time.Duration) []DecayResult {
	if halfLife <= 0 {
		halfLife = 24 * time.Hour
	}

	type hostAcc struct {
		score     float64
		last      time.Time
		count     int
	}

	acc := make(map[string]*hostAcc)
	lambda := math.Log(2) / halfLife.Seconds()

	for _, e := range entries {
		if e.Event == EventNoChange {
			continue
		}
		age := now.Sub(e.Timestamp).Seconds()
		if age < 0 {
			age = 0
		}
		weight := math.Exp(-lambda * age)

		ha, ok := acc[e.Host]
		if !ok {
			ha = &hostAcc{}
			acc[e.Host] = ha
		}
		ha.score += weight
		ha.count++
		if e.Timestamp.After(ha.last) {
			ha.last = e.Timestamp
		}
	}

	results := make([]DecayResult, 0, len(acc))
	for host, ha := range acc {
		results = append(results, DecayResult{
			Host:       host,
			Score:      ha.score,
			LastEvent:  ha.last,
			EventCount: ha.count,
		})
	}

	// Sort descending by score.
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].Score > results[j-1].Score; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}
	return results
}
