package history

import (
	"testing"
	"time"
)

var replayBase = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func replayEntries() []Entry {
	return []Entry{
		{Timestamp: replayBase, Host: "host-a", Opened: []int{80, 443}, Closed: []int{}},
		{Timestamp: replayBase.Add(time.Hour), Host: "host-b", Opened: []int{22}, Closed: []int{80}},
		{Timestamp: replayBase.Add(2 * time.Hour), Host: "host-a", Opened: []int{}, Closed: []int{443}},
	}
}

func TestReplay_AllEvents(t *testing.T) {
	events := Replay(replayEntries(), ReplayOptions{})
	if len(events) != 5 {
		t.Fatalf("expected 5 events, got %d", len(events))
	}
}

func TestReplay_FilterHost(t *testing.T) {
	events := Replay(replayEntries(), ReplayOptions{Host: "host-a"})
	for _, e := range events {
		if e.Host != "host-a" {
			t.Errorf("unexpected host %s", e.Host)
		}
	}
	if len(events) != 3 {
		t.Fatalf("expected 3 events for host-a, got %d", len(events))
	}
}

func TestReplay_Since(t *testing.T) {
	events := Replay(replayEntries(), ReplayOptions{Since: replayBase.Add(30 * time.Minute)})
	if len(events) != 3 {
		t.Fatalf("expected 3 events after since, got %d", len(events))
	}
}

func TestReplay_Limit(t *testing.T) {
	events := Replay(replayEntries(), ReplayOptions{Limit: 2})
	if len(events) != 2 {
		t.Fatalf("expected 2 events with limit, got %d", len(events))
	}
}

func TestReplay_EventTypes(t *testing.T) {
	events := Replay(replayEntries(), ReplayOptions{Host: "host-b"})
	if len(events) != 2 {
		t.Fatalf("expected 2 events for host-b, got %d", len(events))
	}
	if events[0].Event != "opened" || events[1].Event != "closed" {
		t.Errorf("unexpected event types: %v", events)
	}
}
