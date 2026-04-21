package history

import (
	"testing"
	"time"
)

var velocityBase = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func velEntry(host string, event EventType, ports []int, offset time.Duration) Entry {
	return Entry{
		Host:      host,
		Event:     event,
		Ports:     ports,
		Timestamp: velocityBase.Add(offset),
	}
}

func TestVelocity_BasicRate(t *testing.T) {
	entries := []Entry{
		velEntry("host-a", EventOpened, []int{80}, 0),
		velEntry("host-a", EventOpened, []int{443}, time.Hour),
		velEntry("host-a", EventOpened, []int{8080}, 2*time.Hour),
	}
	results := Velocity(entries, time.Hour)
	if len(results) == 0 {
		t.Fatal("expected velocity results")
	}
	found := false
	for _, r := range results {
		if r.Host == "host-a" {
			found = true
			if r.ChangesPerWindow < 1 {
				t.Errorf("expected changes >= 1, got %f", r.ChangesPerWindow)
			}
		}
	}
	if !found {
		t.Error("host-a not found in results")
	}
}

func TestVelocity_Empty(t *testing.T) {
	results := Velocity(nil, time.Hour)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestVelocity_SkipsNoChange(t *testing.T) {
	entries := []Entry{
		velEntry("host-a", EventNoChange, []int{80}, 0),
		velEntry("host-a", EventNoChange, []int{80}, time.Hour),
	}
	results := Velocity(entries, time.Hour)
	for _, r := range results {
		if r.Host == "host-a" && r.TotalChanges > 0 {
			t.Errorf("no-change events should not count, got %d", r.TotalChanges)
		}
	}
}

func TestVelocity_MultipleHosts(t *testing.T) {
	entries := []Entry{
		velEntry("host-a", EventOpened, []int{80}, 0),
		velEntry("host-b", EventOpened, []int{443}, 0),
		velEntry("host-b", EventClosed, []int{22}, time.Hour),
	}
	results := Velocity(entries, time.Hour)
	hosts := map[string]bool{}
	for _, r := range results {
		hosts[r.Host] = true
	}
	if !hosts["host-a"] || !hosts["host-b"] {
		t.Error("expected both hosts in results")
	}
}

func TestVelocity_OrderedByRate(t *testing.T) {
	entries := []Entry{
		velEntry("slow", EventOpened, []int{80}, 0),
		velEntry("fast", EventOpened, []int{80}, 0),
		velEntry("fast", EventOpened, []int{443}, 30*time.Minute),
		velEntry("fast", EventClosed, []int{22}, time.Hour),
	}
	results := Velocity(entries, time.Hour)
	if len(results) < 2 {
		t.Fatal("expected at least 2 results")
	}
	if results[0].ChangesPerWindow < results[1].ChangesPerWindow {
		t.Error("expected results ordered by descending rate")
	}
}
