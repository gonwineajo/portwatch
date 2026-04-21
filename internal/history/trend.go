package history

import (
	"sort"
	"time"
)

// TrendPoint represents the number of change events in a single time bucket.
type TrendPoint struct {
	BucketStart time.Time
	Opened      int
	Closed      int
	Total       int
}

// Trend aggregates opened/closed events into fixed-size time buckets,
// returning a slice of TrendPoints ordered by time.
func Trend(entries []Entry, window time.Duration) []TrendPoint {
	if len(entries) == 0 || window <= 0 {
		return nil
	}

	type key struct {
		bucket int64
	}

	type counts struct {
		opened, closed int
		start          time.Time
	}

	buckets := map[int64]*counts{}

	for _, e := range entries {
		if e.Event == EventNoChange {
			continue
		}
		b := e.Timestamp.Truncate(window).Unix()
		if _, ok := buckets[b]; !ok {
			buckets[b] = &counts{start: e.Timestamp.Truncate(window)}
		}
		switch e.Event {
		case EventOpened:
			buckets[b].opened++
		case EventClosed:
			buckets[b].closed++
		}
	}

	keys := make([]int64, 0, len(buckets))
	for k := range buckets {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	points := make([]TrendPoint, 0, len(keys))
	for _, k := range keys {
		c := buckets[k]
		points = append(points, TrendPoint{
			BucketStart: c.start,
			Opened:      c.opened,
			Closed:      c.closed,
			Total:       c.opened + c.closed,
		})
	}
	return points
}
