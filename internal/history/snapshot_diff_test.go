package history

import (
	"testing"
	"time"
)

var t0 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func sdEntry(host, event string, offset int, ports ...int) Entry {
	return Entry{
		Host:      host,
		Event:     event,
		Timestamp: t0.Add(time.Duration(offset) * time.Minute),
		Ports:     ports,
	}
}

func TestSnapshotDiffs_BasicOpened(t *testing.T) {
	entries := []Entry{
		sdEntry("host-a", "scan", 0, 80),
		sdEntry("host-a", "scan", 1, 80, 443),
	}
	diffs := SnapshotDiffs(entries)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	d := diffs[0]
	if len(d.Opened) != 1 || d.Opened[0] != 443 {
		t.Errorf("expected opened=[443], got %v", d.Opened)
	}
	if len(d.Closed) != 0 {
		t.Errorf("expected no closed ports, got %v", d.Closed)
	}
	if len(d.Stable) != 1 || d.Stable[0] != 80 {
		t.Errorf("expected stable=[80], got %v", d.Stable)
	}
}

func TestSnapshotDiffs_BasicClosed(t *testing.T) {
	entries := []Entry{
		sdEntry("host-b", "scan", 0, 80, 8080),
		sdEntry("host-b", "scan", 1, 80),
	}
	diffs := SnapshotDiffs(entries)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	d := diffs[0]
	if len(d.Closed) != 1 || d.Closed[0] != 8080 {
		t.Errorf("expected closed=[8080], got %v", d.Closed)
	}
}

func TestSnapshotDiffs_SkipsNonScan(t *testing.T) {
	entries := []Entry{
		sdEntry("host-c", "scan", 0, 22),
		sdEntry("host-c", "opened", 1, 443),
		sdEntry("host-c", "scan", 2, 22, 443),
	}
	diffs := SnapshotDiffs(entries)
	// Only two scan entries → one diff
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
}

func TestSnapshotDiffs_MultipleHosts(t *testing.T) {
	entries := []Entry{
		sdEntry("alpha", "scan", 0, 80),
		sdEntry("beta", "scan", 0, 22),
		sdEntry("alpha", "scan", 1, 80, 443),
		sdEntry("beta", "scan", 1, 22, 25),
	}
	diffs := SnapshotDiffs(entries)
	if len(diffs) != 2 {
		t.Fatalf("expected 2 diffs, got %d", len(diffs))
	}
	hosts := map[string]bool{}
	for _, d := range diffs {
		hosts[d.Host] = true
	}
	if !hosts["alpha"] || !hosts["beta"] {
		t.Errorf("expected diffs for alpha and beta, got %v", hosts)
	}
}

func TestSnapshotDiffs_Empty(t *testing.T) {
	diffs := SnapshotDiffs(nil)
	if len(diffs) != 0 {
		t.Errorf("expected empty result, got %v", diffs)
	}
}
