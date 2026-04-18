package snapshot

import (
	"os"
	"sort"
	"testing"
	"time"
)

func makeSnap(host string, ports []int) PortSnapshot {
	return PortSnapshot{Host: host, Ports: ports, ScannedAt: time.Now()}
}

func TestSaveAndLoad(t *testing.T) {
	tmp, err := os.CreateTemp("", "snapshot-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	snap := makeSnap("localhost", []int{80, 443, 8080})
	if err := Save(tmp.Name(), snap); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(tmp.Name())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded.Host != snap.Host {
		t.Errorf("expected host %q, got %q", snap.Host, loaded.Host)
	}
	if len(loaded.Ports) != len(snap.Ports) {
		t.Errorf("expected %d ports, got %d", len(snap.Ports), len(loaded.Ports))
	}
}

func TestCompare_OpenedAndClosed(t *testing.T) {
	prev := makeSnap("host1", []int{80, 443, 22})
	curr := makeSnap("host1", []int{80, 8080})

	diff := Compare(prev, curr)

	if diff.Host != "host1" {
		t.Errorf("unexpected host: %s", diff.Host)
	}

	sort.Ints(diff.Opened)
	if len(diff.Opened) != 1 || diff.Opened[0] != 8080 {
		t.Errorf("expected opened [8080], got %v", diff.Opened)
	}

	sort.Ints(diff.Closed)
	if len(diff.Closed) != 2 || diff.Closed[0] != 22 || diff.Closed[1] != 443 {
		t.Errorf("expected closed [22 443], got %v", diff.Closed)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	prev := makeSnap("host1", []int{80, 443})
	curr := makeSnap("host1", []int{80, 443})
	diff := Compare(prev, curr)
	if len(diff.Opened) != 0 || len(diff.Closed) != 0 {
		t.Errorf("expected no changes, got opened=%v closed=%v", diff.Opened, diff.Closed)
	}
}
