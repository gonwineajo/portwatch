package history

import (
	"testing"
	"time"
)

var staleBase = time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC)

func staleEntry(host string, event EventType, ports []int, t time.Time) Entry {
	return Entry{Host: host, Event: event, Ports: ports, Timestamp: t}
}

func TestAnalyseStaleness_BasicOrder(t *testing.T) {
	now := staleBase
	entries := []Entry{
		staleEntry("host-a", EventScan, []int{80}, now.Add(-72*time.Hour)),
		staleEntry("host-b", EventScan, []int{443}, now.Add(-1*time.Hour)),
	}

	results := AnalyseStaleness(entries, 48*time.Hour, now)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	// host-a is older, should appear first
	if results[0].Host != "host-a" {
		t.Errorf("expected host-a first, got %s", results[0].Host)
	}
}

func TestAnalyseStaleness_IsStaleFlag(t *testing.T) {
	now := staleBase
	entries := []Entry{
		staleEntry("host-a", EventScan, []int{80}, now.Add(-72*time.Hour)),
		staleEntry("host-b", EventScan, []int{443}, now.Add(-1*time.Hour)),
	}

	results := AnalyseStaleness(entries, 48*time.Hour, now)

	byHost := make(map[string]StalenessResult)
	for _, r := range results {
		byHost[r.Host] = r
	}

	if !byHost["host-a"].IsStale {
		t.Error("expected host-a to be stale")
	}
	if byHost["host-b"].IsStale {
		t.Error("expected host-b to not be stale")
	}
}

func TestAnalyseStaleness_SkipsNonScan(t *testing.T) {
	now := staleBase
	entries := []Entry{
		staleEntry("host-a", EventOpened, []int{80}, now.Add(-1*time.Hour)),
		staleEntry("host-a", EventClosed, []int{443}, now.Add(-2*time.Hour)),
	}

	results := AnalyseStaleness(entries, time.Hour, now)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestAnalyseStaleness_UsesLatestScan(t *testing.T) {
	now := staleBase
	entries := []Entry{
		staleEntry("host-a", EventScan, []int{80}, now.Add(-100*time.Hour)),
		staleEntry("host-a", EventScan, []int{80, 443}, now.Add(-5*time.Hour)),
	}

	results := AnalyseStaleness(entries, 48*time.Hour, now)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].IsStale {
		t.Error("expected host-a to not be stale (latest scan is recent)")
	}
}

func TestAnalyseStaleness_Empty(t *testing.T) {
	results := AnalyseStaleness(nil, time.Hour, staleBase)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}
