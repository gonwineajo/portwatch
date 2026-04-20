package history

import (
	"testing"
	"time"
)

var t0 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func compactEntry(host string, offset int, ports []int, event EventType) Entry {
	return Entry{
		Timestamp: t0.Add(time.Duration(offset) * time.Minute),
		Host:      host,
		OpenPorts: ports,
		Event:     event,
	}
}

func TestCompact_RemovesDuplicateScans(t *testing.T) {
	entries := []Entry{
		compactEntry("host-a", 0, []int{80, 443}, EventScan),
		compactEntry("host-a", 1, []int{80, 443}, EventScan), // duplicate
		compactEntry("host-a", 2, []int{80, 443}, EventScan), // duplicate
	}
	out, res := Compact(entries)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if res.Removed != 2 {
		t.Errorf("expected 2 removed, got %d", res.Removed)
	}
}

func TestCompact_KeepsChangedPorts(t *testing.T) {
	entries := []Entry{
		compactEntry("host-a", 0, []int{80}, EventScan),
		compactEntry("host-a", 1, []int{80, 443}, EventScan), // ports changed
	}
	out, res := Compact(entries)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if res.Removed != 0 {
		t.Errorf("expected 0 removed, got %d", res.Removed)
	}
}

func TestCompact_KeepsOpenedClosedEvents(t *testing.T) {
	entries := []Entry{
		compactEntry("host-a", 0, []int{80}, EventScan),
		{Timestamp: t0.Add(1 * time.Minute), Host: "host-a", OpenPorts: []int{80}, Event: EventOpened},
		compactEntry("host-a", 2, []int{80}, EventScan), // same ports but after an event
	}
	out, _ := Compact(entries)
	// The EventOpened entry must survive regardless.
	var hasOpened bool
	for _, e := range out {
		if e.Event == EventOpened {
			hasOpened = true
		}
	}
	if !hasOpened {
		t.Error("expected EventOpened entry to be retained")
	}
}

func TestCompact_MultipleHosts(t *testing.T) {
	entries := []Entry{
		compactEntry("host-a", 0, []int{80}, EventScan),
		compactEntry("host-b", 0, []int{22}, EventScan),
		compactEntry("host-a", 1, []int{80}, EventScan), // dup for host-a
		compactEntry("host-b", 1, []int{22}, EventScan), // dup for host-b
	}
	out, res := Compact(entries)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if res.Removed != 2 {
		t.Errorf("expected 2 removed, got %d", res.Removed)
	}
}

func TestCompact_Empty(t *testing.T) {
	out, res := Compact(nil)
	if len(out) != 0 {
		t.Errorf("expected empty result")
	}
	if res.Before != 0 || res.After != 0 || res.Removed != 0 {
		t.Errorf("unexpected result: %+v", res)
	}
}
