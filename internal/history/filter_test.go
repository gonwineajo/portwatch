package history

import (
	"testing"
	"time"
)

var base = time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC)

var filterEntries = []Entry{
	{Host: "host-a", Ports: []int{80}, Event: "opened", Timestamp: base},
	{Host: "host-b", Ports: []int{443}, Event: "opened", Timestamp: base.Add(time.Hour)},
	{Host: "host-a", Ports: []int{22}, Event: "closed", Timestamp: base.Add(2 * time.Hour)},
	{Host: "host-b", Ports: []int{8080}, Event: "closed", Timestamp: base.Add(3 * time.Hour)},
}

func TestFilter_ByHost(t *testing.T) {
	out := Filter(filterEntries, FilterOptions{Host: "host-a"})
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestFilter_ByEvent(t *testing.T) {
	out := Filter(filterEntries, FilterOptions{Event: "closed"})
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestFilter_Since(t *testing.T) {
	out := Filter(filterEntries, FilterOptions{Since: base.Add(90 * time.Minute)})
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestFilter_Until(t *testing.T) {
	out := Filter(filterEntries, FilterOptions{Until: base.Add(30 * time.Minute)})
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
}

func TestFilter_Limit(t *testing.T) {
	out := Filter(filterEntries, FilterOptions{Limit: 2})
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestFilter_Combined(t *testing.T) {
	out := Filter(filterEntries, FilterOptions{Host: "host-b", Event: "opened"})
	if len(out) != 1 || out[0].Ports[0] != 443 {
		t.Fatalf("unexpected result: %+v", out)
	}
}

func TestFilter_Empty(t *testing.T) {
	out := Filter([]Entry{}, FilterOptions{Host: "host-a"})
	if len(out) != 0 {
		t.Fatalf("expected 0, got %d", len(out))
	}
}
