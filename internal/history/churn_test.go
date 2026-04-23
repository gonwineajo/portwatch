package history

import (
	"testing"
	"time"
)

var churnBase = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func churnEntry(host, event string, ports []int, offset int) Entry {
	return Entry{
		Timestamp: churnBase.Add(time.Duration(offset) * time.Hour),
		Host:      host,
		Event:     event,
		Ports:     ports,
	}
}

func TestAnalyseChurn_BasicOrder(t *testing.T) {
	entries := []Entry{
		churnEntry("a", EventOpened, []int{80}, 0),
		churnEntry("a", EventClosed, []int{80}, 1),
		churnEntry("a", EventOpened, []int{80}, 2),
		churnEntry("b", EventOpened, []int{443}, 0),
	}
	results := AnalyseChurn(entries)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Host != "a" {
		t.Errorf("expected host a first, got %s", results[0].Host)
	}
	if results[0].TotalFlips != 3 {
		t.Errorf("expected 3 flips for a, got %d", results[0].TotalFlips)
	}
	if results[0].UniquePorts != 1 {
		t.Errorf("expected 1 unique port for a, got %d", results[0].UniquePorts)
	}
}

func TestAnalyseChurn_SkipsNoChange(t *testing.T) {
	entries := []Entry{
		churnEntry("a", EventNoChange, []int{80}, 0),
		churnEntry("a", EventNoChange, []int{80}, 1),
	}
	results := AnalyseChurn(entries)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestAnalyseChurn_Score(t *testing.T) {
	entries := []Entry{
		churnEntry("x", EventOpened, []int{80, 443}, 0),
		churnEntry("x", EventClosed, []int{80}, 1),
	}
	results := AnalyseChurn(entries)
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	// 2 flips / 2 unique ports = 1.0
	if results[0].Score != 1.0 {
		t.Errorf("expected score 1.0, got %f", results[0].Score)
	}
}

func TestAnalyseChurn_Empty(t *testing.T) {
	results := AnalyseChurn(nil)
	if len(results) != 0 {
		t.Errorf("expected empty results")
	}
}

func TestAnalyseChurn_MultipleHosts(t *testing.T) {
	entries := []Entry{
		churnEntry("host1", EventOpened, []int{22, 80}, 0),
		churnEntry("host2", EventOpened, []int{22}, 0),
		churnEntry("host2", EventClosed, []int{22}, 1),
		churnEntry("host2", EventOpened, []int{22}, 2),
	}
	results := AnalyseChurn(entries)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	// host2: 3 flips / 1 unique = 3.0 > host1: 1 flip / 2 unique = 0.5
	if results[0].Host != "host2" {
		t.Errorf("expected host2 first, got %s", results[0].Host)
	}
}
