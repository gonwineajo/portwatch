package history

import (
	"testing"
	"time"
)

func driftEntry(host, event string, ports []int, ts time.Time) Entry {
	return Entry{
		Timestamp: ts,
		Host:      host,
		Event:     event,
		Ports:     ports,
	}
}

func TestAnalyseDrift_BasicDrift(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		driftEntry("host-a", EventScan, []int{80, 443}, now.Add(-2*time.Hour)),
		driftEntry("host-a", EventScan, []int{80, 443, 8080}, now),
	}

	results := AnalyseDrift(entries)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if r.Host != "host-a" {
		t.Errorf("unexpected host: %s", r.Host)
	}
	if len(r.Added) != 1 || r.Added[0] != 8080 {
		t.Errorf("expected added=[8080], got %v", r.Added)
	}
	if len(r.Removed) != 0 {
		t.Errorf("expected no removed ports, got %v", r.Removed)
	}
}

func TestAnalyseDrift_RemovedPorts(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		driftEntry("host-b", EventScan, []int{22, 80, 443}, now.Add(-time.Hour)),
		driftEntry("host-b", EventScan, []int{80}, now),
	}

	results := AnalyseDrift(entries)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if len(r.Removed) != 2 {
		t.Errorf("expected 2 removed ports, got %v", r.Removed)
	}
}

func TestAnalyseDrift_NoDrift(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		driftEntry("host-c", EventScan, []int{80, 443}, now.Add(-time.Hour)),
		driftEntry("host-c", EventScan, []int{80, 443}, now),
	}

	results := AnalyseDrift(entries)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Score != 0 {
		t.Errorf("expected score 0 for no drift, got %f", results[0].Score)
	}
}

func TestAnalyseDrift_SkipsNoBaseline(t *testing.T) {
	now := time.Now()
	// Only one scan — SetBaseline picks it as baseline, latest is same scan.
	// Two scans needed for meaningful drift; a host with only non-scan events
	// should be skipped entirely.
	entries := []Entry{
		driftEntry("host-d", EventOpened, []int{80}, now),
	}

	results := AnalyseDrift(entries)
	if len(results) != 0 {
		t.Errorf("expected 0 results for host with no scan events, got %d", len(results))
	}
}

func TestAnalyseDrift_OrderedByScore(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		// host-x: small drift
		driftEntry("host-x", EventScan, []int{80, 443}, now.Add(-time.Hour)),
		driftEntry("host-x", EventScan, []int{80, 443, 8080}, now),
		// host-y: larger drift
		driftEntry("host-y", EventScan, []int{22, 80, 443, 3306}, now.Add(-time.Hour)),
		driftEntry("host-y", EventScan, []int{80}, now),
	}

	results := AnalyseDrift(entries)
	if len(results) < 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Score < results[1].Score {
		t.Errorf("results not ordered by descending score: %v", results)
	}
}
