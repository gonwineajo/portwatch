package history

import (
	"testing"
	"time"
)

var groupEntries = []Entry{
	{Host: "host-b", Event: EventOpened, Ports: []int{80}, Timestamp: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
	{Host: "host-a", Event: EventOpened, Ports: []int{443}, Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
	{Host: "host-a", Event: EventClosed, Ports: []int{22}, Timestamp: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
	{Host: "host-b", Event: EventClosed, Ports: []int{8080}, Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
}

func TestGroupByHost_Keys(t *testing.T) {
	groups := GroupByHost(groupEntries)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Key != "host-a" || groups[1].Key != "host-b" {
		t.Errorf("unexpected keys: %s, %s", groups[0].Key, groups[1].Key)
	}
}

func TestGroupByHost_Counts(t *testing.T) {
	groups := GroupByHost(groupEntries)
	for _, g := range groups {
		if len(g.Entries) != 2 {
			t.Errorf("host %s: expected 2 entries, got %d", g.Key, len(g.Entries))
		}
	}
}

func TestGroupByEvent_Keys(t *testing.T) {
	groups := GroupByEvent(groupEntries)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Key != "closed" || groups[1].Key != "opened" {
		t.Errorf("unexpected event keys: %s, %s", groups[0].Key, groups[1].Key)
	}
}

func TestGroupByDate_Keys(t *testing.T) {
	groups := GroupByDate(groupEntries)
	if len(groups) != 3 {
		t.Fatalf("expected 3 date groups, got %d", len(groups))
	}
	expected := []string{"2024-01-01", "2024-01-02", "2024-01-03"}
	for i, g := range groups {
		if g.Key != expected[i] {
			t.Errorf("index %d: expected %s, got %s", i, expected[i], g.Key)
		}
	}
}

func TestGroupByHost_Empty(t *testing.T) {
	groups := GroupByHost(nil)
	if len(groups) != 0 {
		t.Errorf("expected empty result for nil input, got %d groups", len(groups))
	}
}
