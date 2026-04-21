package history

import (
	"testing"
	"time"
)

var normBase = time.Date(2024, 3, 1, 12, 0, 0, 0, time.UTC)

func normEntry(host, event string, ports, opened, closed []int, offset int) Entry {
	return Entry{
		Host:        host,
		Event:       event,
		Ports:       ports,
		OpenedPorts: opened,
		ClosedPorts: closed,
		ScannedAt:   normBase.Add(time.Duration(offset) * time.Minute),
	}
}

func TestNormalize_SortPorts(t *testing.T) {
	entries := []Entry{
		normEntry("h", EventOpened, []int{443, 80, 22}, []int{443, 80}, nil, 0),
	}
	out := Normalize(entries, NormalizeOptions{SortPorts: true})
	if out[0].Ports[0] != 22 {
		t.Errorf("expected Ports[0]=22, got %d", out[0].Ports[0])
	}
	if out[0].OpenedPorts[0] != 80 {
		t.Errorf("expected OpenedPorts[0]=80, got %d", out[0].OpenedPorts[0])
	}
}

func TestNormalize_DeduplicateScans(t *testing.T) {
	ports := []int{80, 443}
	entries := []Entry{
		normEntry("h", EventNoChange, ports, nil, nil, 0),
		normEntry("h", EventNoChange, ports, nil, nil, 1), // duplicate
		normEntry("h", EventNoChange, ports, nil, nil, 2), // duplicate
	}
	out := Normalize(entries, NormalizeOptions{DeduplicateScans: true})
	if len(out) != 1 {
		t.Errorf("expected 1 entry after dedup, got %d", len(out))
	}
}

func TestNormalize_DropEmpty(t *testing.T) {
	entries := []Entry{
		normEntry("h", EventOpened, nil, nil, nil, 0), // empty opened/closed → drop
		normEntry("h", EventOpened, nil, []int{80}, nil, 1), // has opened → keep
	}
	out := Normalize(entries, NormalizeOptions{DropEmpty: true})
	if len(out) != 1 {
		t.Errorf("expected 1 entry, got %d", len(out))
	}
	if out[0].OpenedPorts[0] != 80 {
		t.Errorf("expected port 80, got %d", out[0].OpenedPorts[0])
	}
}

func TestNormalize_DoesNotMutateOriginal(t *testing.T) {
	orig := []Entry{
		normEntry("h", EventOpened, []int{443, 80}, []int{443, 80}, nil, 0),
	}
	Normalize(orig, NormalizeOptions{SortPorts: true})
	if orig[0].Ports[0] != 443 {
		t.Errorf("original entry mutated: Ports[0]=%d", orig[0].Ports[0])
	}
}

func TestNormalize_Empty(t *testing.T) {
	out := Normalize(nil, NormalizeOptions{SortPorts: true, DeduplicateScans: true})
	if len(out) != 0 {
		t.Errorf("expected empty result, got %d entries", len(out))
	}
}
