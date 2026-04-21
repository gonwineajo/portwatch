package history

import (
	"testing"
	"time"
)

func heatEntry(host string, opened, closed []int, hour int) Entry {
	base := time.Date(2024, 1, 15, hour, 0, 0, 0, time.UTC)
	return Entry{
		Timestamp:   base,
		Host:        host,
		Event:       "scan",
		OpenedPorts: opened,
		ClosedPorts: closed,
	}
}

func TestHeatmap_BasicCounts(t *testing.T) {
	entries := []Entry{
		heatEntry("host-a", []int{80}, nil, 9),
		heatEntry("host-a", []int{443}, nil, 9),
		heatEntry("host-a", nil, []int{80}, 14),
	}

	cells := Heatmap(entries)
	if len(cells) != 2 {
		t.Fatalf("expected 2 cells, got %d", len(cells))
	}
	if cells[0].Hour != 9 || cells[0].Changes != 2 {
		t.Errorf("expected hour=9 changes=2, got hour=%d changes=%d", cells[0].Hour, cells[0].Changes)
	}
	if cells[1].Hour != 14 || cells[1].Changes != 1 {
		t.Errorf("expected hour=14 changes=1, got hour=%d changes=%d", cells[1].Hour, cells[1].Changes)
	}
}

func TestHeatmap_SkipsNoChange(t *testing.T) {
	entries := []Entry{
		heatEntry("host-b", nil, nil, 10),
		heatEntry("host-b", []int{22}, nil, 10),
	}

	cells := Heatmap(entries)
	if len(cells) != 1 {
		t.Fatalf("expected 1 cell, got %d", len(cells))
	}
}

func TestHeatmap_Empty(t *testing.T) {
	cells := Heatmap(nil)
	if len(cells) != 0 {
		t.Errorf("expected empty heatmap for nil input")
	}
}

func TestHeatmap_MultipleHosts(t *testing.T) {
	entries := []Entry{
		heatEntry("alpha", []int{80}, nil, 8),
		heatEntry("beta", []int{443}, nil, 8),
		heatEntry("alpha", nil, []int{80}, 8),
	}

	cells := Heatmap(entries)
	if len(cells) != 2 {
		t.Fatalf("expected 2 cells (one per host), got %d", len(cells))
	}
	for _, c := range cells {
		if c.Hour != 8 {
			t.Errorf("expected hour 8, got %d", c.Hour)
		}
	}
}

func TestHeatmap_SpansMultipleDays(t *testing.T) {
	// Entries on different calendar days but same hour should be aggregated by hour.
	day1 := Entry{
		Timestamp:   time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC),
		Host:        "host-c",
		Event:       "scan",
		OpenedPorts: []int{80},
	}
	day2 := Entry{
		Timestamp:   time.Date(2024, 1, 11, 9, 0, 0, 0, time.UTC),
		Host:        "host-c",
		Event:       "scan",
		OpenedPorts: []int{443},
	}

	cells := Heatmap([]Entry{day1, day2})
	if len(cells) != 1 {
		t.Fatalf("expected 1 cell for same hour across days, got %d", len(cells))
	}
	if cells[0].Hour != 9 || cells[0].Changes != 2 {
		t.Errorf("expected hour=9 changes=2, got hour=%d changes=%d", cells[0].Hour, cells[0].Changes)
	}
}

func TestPeakHour_ReturnsCorrectHour(t *testing.T) {
	entries := []Entry{
		heatEntry("h", []int{80, 443, 22}, nil, 9),
		heatEntry("h", []int{8080}, nil, 14),
	}

	peak := PeakHour(entries)
	if peak != 9 {
		t.Errorf("expected peak hour 9, got %d", peak)
	}
}

func TestPeakHour_Empty(t *testing.T) {
	if PeakHour(nil) != -1 {
		t.Errorf("expected -1 for empty entries")
	}
}
