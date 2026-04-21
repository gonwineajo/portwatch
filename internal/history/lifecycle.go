package history

// LifecycleSummary describes how long a port stayed open across its recorded history.
type LifecycleSummary struct {
	Host          string
	Port          int
	OpenCount     int
	CloseCount    int
	TotalOpenTime int64 // seconds between paired open/close events
	CurrentlyOpen bool
}

// Lifecycle computes open/close counts and total open durations from chains.
// Pairs each opened event with the next closed event to estimate open duration.
func Lifecycle(chains []Chain) []LifecycleSummary {
	out := make([]LifecycleSummary, 0, len(chains))
	for _, c := range chains {
		ls := LifecycleSummary{Host: c.Host, Port: c.Port}
		var lastOpen int64 = -1
		for _, step := range c.Steps {
			switch step.Event {
			case EventOpened:
				ls.OpenCount++
				lastOpen = step.Timestamp
				ls.CurrentlyOpen = true
			case EventClosed:
				ls.CloseCount++
				ls.CurrentlyOpen = false
				if lastOpen >= 0 {
					ls.TotalOpenTime += step.Timestamp - lastOpen
					lastOpen = -1
				}
			}
		}
		out = append(out, ls)
	}
	return out
}

// LongestOpen returns the LifecycleSummary with the greatest TotalOpenTime.
// Returns nil if the slice is empty.
func LongestOpen(summaries []LifecycleSummary) *LifecycleSummary {
	if len(summaries) == 0 {
		return nil
	}
	best := &summaries[0]
	for i := 1; i < len(summaries); i++ {
		if summaries[i].TotalOpenTime > best.TotalOpenTime {
			best = &summaries[i]
		}
	}
	return best
}
