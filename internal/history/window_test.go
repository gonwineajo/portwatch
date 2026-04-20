package history

import (
	"testing"
	"time"
)

var windowNow = time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC)

func windowEntry(host, event string, minsAgo int, ports []int) Entry {
	return Entry{
		Timestamp: windowNow.Add(-time.Duration(minsAgo) * time.Minute),
		Host:      host,
		Event:     event,
		Ports:     ports,
	}
}

func TestRollingWindow_BasicCounts(t *testing.T) {
	entries := []Entry{
		windowEntry("host-a", "opened", 10, []int{80}),
		windowEntry("host-a", "opened", 20, []int{443}),
		windowEntry("host-a", "closed", 30, []int{22}),
		windowEntry("host-b", "opened", 5, []int{8080}),
	}

	stats := RollingWindow(entries, 60*time.Minute, windowNow)

	if len(stats) != 2 {
		t.Fatalf("expected 2 hosts, got %d", len(stats))
	}
	// host-a should be first (Total=3 > host-b Total=1)
	if stats[0].Host != "host-a" {
		t.Errorf("expected host-a first, got %s", stats[0].Host)
	}
	if stats[0].Opened != 2 || stats[0].Closed != 1 || stats[0].Total != 3 {
		t.Errorf("unexpected counts for host-a: %+v", stats[0])
	}
	if stats[1].Host != "host-b" || stats[1].Opened != 1 {
		t.Errorf("unexpected host-b stats: %+v", stats[1])
	}
}

func TestRollingWindow_ExcludesOldEntries(t *testing.T) {
	entries := []Entry{
		windowEntry("host-a", "opened", 30, []int{80}),
		windowEntry("host-a", "opened", 120, []int{443}), // outside window
	}

	stats := RollingWindow(entries, 60*time.Minute, windowNow)

	if len(stats) != 1 {
		t.Fatalf("expected 1 result, got %d", len(stats))
	}
	if stats[0].Opened != 1 {
		t.Errorf("expected 1 opened, got %d", stats[0].Opened)
	}
}

func TestRollingWindow_Empty(t *testing.T) {
	stats := RollingWindow([]Entry{}, 60*time.Minute, windowNow)
	if len(stats) != 0 {
		t.Errorf("expected empty result, got %d entries", len(stats))
	}
}

func TestRollingWindow_SinceUntilSet(t *testing.T) {
	entries := []Entry{
		windowEntry("host-a", "opened", 5, []int{80}),
	}
	window := 30 * time.Minute
	stats := RollingWindow(entries, window, windowNow)

	if len(stats) != 1 {
		t.Fatalf("expected 1 result")
	}
	expectedSince := windowNow.Add(-window)
	if !stats[0].Since.Equal(expectedSince) {
		t.Errorf("Since mismatch: got %v, want %v", stats[0].Since, expectedSince)
	}
	if !stats[0].Until.Equal(windowNow) {
		t.Errorf("Until mismatch: got %v, want %v", stats[0].Until, windowNow)
	}
}
