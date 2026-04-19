package history

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func snap(ports []int) snapshot.Snapshot {
	return snapshot.Snapshot{Host: "host1", Ports: ports}
}

func TestDiff_Opened(t *testing.T) {
	prev := snap([]int{80})
	curr := snap([]int{80, 443})
	e := Diff("host1", prev, curr, time.Now())
	if e.Event != "opened" {
		t.Errorf("expected opened, got %s", e.Event)
	}
	if len(e.Opened) != 1 || e.Opened[0] != 443 {
		t.Errorf("expected [443], got %v", e.Opened)
	}
}

func TestDiff_Closed(t *testing.T) {
	prev := snap([]int{80, 443})
	curr := snap([]int{80})
	e := Diff("host1", prev, curr, time.Now())
	if e.Event != "closed" {
		t.Errorf("expected closed, got %s", e.Event)
	}
	if len(e.Closed) != 1 || e.Closed[0] != 443 {
		t.Errorf("expected [443], got %v", e.Closed)
	}
}

func TestDiff_Unchanged(t *testing.T) {
	prev := snap([]int{80})
	curr := snap([]int{80})
	e := Diff("host1", prev, curr, time.Now())
	if e.Event != "unchanged" {
		t.Errorf("expected unchanged, got %s", e.Event)
	}
}

func TestHasChanges(t *testing.T) {
	if HasChanges(Entry{Event: "unchanged"}) {
		t.Error("expected no changes")
	}
	if !HasChanges(Entry{Event: "opened", Opened: []int{8080}}) {
		t.Error("expected changes")
	}
}
