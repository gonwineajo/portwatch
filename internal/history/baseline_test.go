package history

import (
	"sort"
	"testing"
	"time"
)

var baselineTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func baselineEntry(host string, event EventType, ports []int, offset int) Entry {
	return Entry{
		Host:      host,
		Event:     event,
		Ports:     ports,
		Timestamp: baselineTime.Add(time.Duration(offset) * time.Minute),
	}
}

func TestSetBaseline_PicksLatestScan(t *testing.T) {
	entries := []Entry{
		baselineEntry("host-a", EventScan, []int{80, 443}, 0),
		baselineEntry("host-a", EventScan, []int{80, 443, 8080}, 10),
		baselineEntry("host-a", EventOpened, []int{8080}, 5),
	}
	baselines := SetBaseline(entries)
	if len(baselines) != 1 {
		t.Fatalf("expected 1 baseline, got %d", len(baselines))
	}
	if len(baselines[0].Ports) != 3 {
		t.Errorf("expected 3 ports, got %v", baselines[0].Ports)
	}
}

func TestSetBaseline_MultipleHosts(t *testing.T) {
	entries := []Entry{
		baselineEntry("host-a", EventScan, []int{80}, 0),
		baselineEntry("host-b", EventScan, []int{22, 443}, 0),
	}
	baselines := SetBaseline(entries)
	if len(baselines) != 2 {
		t.Errorf("expected 2 baselines, got %d", len(baselines))
	}
}

func TestSetBaseline_IgnoresNonScan(t *testing.T) {
	entries := []Entry{
		baselineEntry("host-a", EventOpened, []int{8080}, 0),
		baselineEntry("host-a", EventClosed, []int{443}, 1),
	}
	baselines := SetBaseline(entries)
	if len(baselines) != 0 {
		t.Errorf("expected 0 baselines, got %d", len(baselines))
	}
}

func TestDeviatesFromBaseline_Opened(t *testing.T) {
	bl := Baseline{Host: "host-a", Ports: []int{80, 443}}
	opened, closed := DeviatesFromBaseline(bl, []int{80, 443, 8080})
	if len(opened) != 1 || opened[0] != 8080 {
		t.Errorf("expected [8080] opened, got %v", opened)
	}
	if len(closed) != 0 {
		t.Errorf("expected no closed ports, got %v", closed)
	}
}

func TestDeviatesFromBaseline_Closed(t *testing.T) {
	bl := Baseline{Host: "host-a", Ports: []int{80, 443, 22}}
	opened, closed := DeviatesFromBaseline(bl, []int{80, 443})
	sort.Ints(closed)
	if len(closed) != 1 || closed[0] != 22 {
		t.Errorf("expected [22] closed, got %v", closed)
	}
	if len(opened) != 0 {
		t.Errorf("expected no opened ports, got %v", opened)
	}
}

func TestDeviatesFromBaseline_NoChange(t *testing.T) {
	bl := Baseline{Host: "host-a", Ports: []int{80, 443}}
	opened, closed := DeviatesFromBaseline(bl, []int{80, 443})
	if len(opened) != 0 || len(closed) != 0 {
		t.Errorf("expected no deviations, got opened=%v closed=%v", opened, closed)
	}
}
