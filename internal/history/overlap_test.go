package history

import (
	"testing"
	"time"
)

var overlapBase = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func overlapEntry(host string, ports []int, offset int) Entry {
	return Entry{
		Timestamp: overlapBase.Add(time.Duration(offset) * time.Minute),
		Host:      host,
		Event:     EventScan,
		Ports:     ports,
	}
}

func TestAnalyseOverlap_SharedPorts(t *testing.T) {
	entries := []Entry{
		overlapEntry("host-a", []int{80, 443, 8080}, 0),
		overlapEntry("host-b", []int{80, 443, 9090}, 0),
	}
	results := AnalyseOverlap(entries)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if len(r.Shared) != 2 {
		t.Errorf("expected 2 shared ports, got %v", r.Shared)
	}
	if len(r.OnlyA) != 1 || r.OnlyA[0] != 8080 {
		t.Errorf("expected onlyA=[8080], got %v", r.OnlyA)
	}
	if len(r.OnlyB) != 1 || r.OnlyB[0] != 9090 {
		t.Errorf("expected onlyB=[9090], got %v", r.OnlyB)
	}
}

func TestAnalyseOverlap_JaccardFullOverlap(t *testing.T) {
	entries := []Entry{
		overlapEntry("host-a", []int{80, 443}, 0),
		overlapEntry("host-b", []int{80, 443}, 0),
	}
	results := AnalyseOverlap(entries)
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if results[0].JaccardSim != 1.0 {
		t.Errorf("expected jaccard=1.0, got %f", results[0].JaccardSim)
	}
}

func TestAnalyseOverlap_JaccardNoOverlap(t *testing.T) {
	entries := []Entry{
		overlapEntry("host-a", []int{80}, 0),
		overlapEntry("host-b", []int{443}, 0),
	}
	results := AnalyseOverlap(entries)
	if results[0].JaccardSim != 0.0 {
		t.Errorf("expected jaccard=0.0, got %f", results[0].JaccardSim)
	}
}

func TestAnalyseOverlap_SkipsNonScan(t *testing.T) {
	entries := []Entry{
		overlapEntry("host-a", []int{80}, 0),
		{Timestamp: overlapBase, Host: "host-b", Event: EventOpened, Ports: []int{80}},
	}
	// host-b has no scan entry, so only 1 host qualifies → no pairs
	results := AnalyseOverlap(entries)
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}

func TestAnalyseOverlap_UsesLatestScan(t *testing.T) {
	entries := []Entry{
		overlapEntry("host-a", []int{80}, 0),
		overlapEntry("host-a", []int{80, 443}, 10), // newer
		overlapEntry("host-b", []int{80, 443}, 0),
	}
	results := AnalyseOverlap(entries)
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if results[0].JaccardSim != 1.0 {
		t.Errorf("expected jaccard=1.0 using latest scan, got %f", results[0].JaccardSim)
	}
}

func TestAnalyseOverlap_Empty(t *testing.T) {
	results := AnalyseOverlap(nil)
	if len(results) != 0 {
		t.Errorf("expected empty results")
	}
}
