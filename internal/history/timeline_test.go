package history

import (
	"testing"
	"time"
)

func TestTimeline_BasicBuckets(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	entries := []Entry{
		{Timestamp: base, Host: "h1", Opened: []int{80}, Closed: []int{}},
		{Timestamp: base.Add(30 * time.Minute), Host: "h1", Opened: []int{443}, Closed: []int{}},
		{Timestamp: base.Add(90 * time.Minute), Host: "h1", Opened: []int{}, Closed: []int{80}},
	}

	buckets := Timeline(entries, time.Hour)
	if len(buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(buckets))
	}
	if buckets[0].Opened != 2 {
		t.Errorf("bucket 0 opened: want 2, got %d", buckets[0].Opened)
	}
	if buckets[1].Closed != 1 {
		t.Errorf("bucket 1 closed: want 1, got %d", buckets[1].Closed)
	}
}

func TestTimeline_Empty(t *testing.T) {
	if Timeline(nil, time.Hour) != nil {
		t.Error("expected nil for empty entries")
	}
}

func TestTimeline_ZeroWindow(t *testing.T) {
	base := time.Now()
	entries := []Entry{{Timestamp: base, Host: "h1", Opened: []int{80}}}
	if Timeline(entries, 0) != nil {
		t.Error("expected nil for zero window")
	}
}

func TestTimeline_SingleBucket(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	entries := []Entry{
		{Timestamp: base, Host: "h1", Opened: []int{80, 443}, Closed: []int{22}},
	}
	buckets := Timeline(entries, time.Hour)
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(buckets))
	}
	if buckets[0].Opened != 2 || buckets[0].Closed != 1 {
		t.Errorf("unexpected counts: %+v", buckets[0])
	}
}
