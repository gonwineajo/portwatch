package history

import (
	"testing"
	"time"
)

var trendBase = time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)

func trendEntry(event EventType, ports []int, offset time.Duration) Entry {
	return Entry{
		Host:      "host-a",
		Event:     event,
		Ports:     ports,
		Timestamp: trendBase.Add(offset),
	}
}

func TestTrend_BasicBuckets(t *testing.T) {
	entries := []Entry{
		trendEntry(EventOpened, []int{80}, 0),
		trendEntry(EventOpened, []int{443}, 30*time.Minute),
		trendEntry(EventClosed, []int{22}, 2*time.Hour),
	}
	points := Trend(entries, time.Hour)
	if len(points) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(points))
	}
	if points[0].Opened != 2 {
		t.Errorf("expected 2 opened in first bucket, got %d", points[0].Opened)
	}
	if points[1].Closed != 1 {
		t.Errorf("expected 1 closed in second bucket, got %d", points[1].Closed)
	}
}

func TestTrend_SkipsNoChange(t *testing.T) {
	entries := []Entry{
		trendEntry(EventNoChange, []int{80}, 0),
		trendEntry(EventOpened, []int{443}, 0),
	}
	points := Trend(entries, time.Hour)
	if len(points) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(points))
	}
	if points[0].Total != 1 {
		t.Errorf("expected total 1, got %d", points[0].Total)
	}
}

func TestTrend_Empty(t *testing.T) {
	points := Trend(nil, time.Hour)
	if len(points) != 0 {
		t.Errorf("expected empty, got %d", len(points))
	}
}

func TestTrend_ZeroWindow(t *testing.T) {
	entries := []Entry{
		trendEntry(EventOpened, []int{80}, 0),
	}
	points := Trend(entries, 0)
	if len(points) != 0 {
		t.Errorf("expected empty for zero window, got %d", len(points))
	}
}

func TestTrend_OrderedByTime(t *testing.T) {
	entries := []Entry{
		trendEntry(EventClosed, []int{22}, 3*time.Hour),
		trendEntry(EventOpened, []int{80}, 0),
		trendEntry(EventOpened, []int{443}, time.Hour),
	}
	points := Trend(entries, time.Hour)
	for i := 1; i < len(points); i++ {
		if points[i].BucketStart.Before(points[i-1].BucketStart) {
			t.Error("trend points not ordered by time")
		}
	}
}
