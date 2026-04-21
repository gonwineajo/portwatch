package history

import (
	"testing"
	"time"
)

var cadenceBase = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func cadenceEntry(host string, port int, event string, offset time.Duration) Entry {
	return Entry{
		Host:      host,
		Event:     event,
		Ports:     []int{port},
		Timestamp: cadenceBase.Add(offset),
	}
}

func TestAnalyseCadence_RegularPort(t *testing.T) {
	// Port 80 opens every hour — very regular.
	entries := []Entry{
		cadenceEntry("host-a", 80, "opened", 0),
		cadenceEntry("host-a", 80, "opened", time.Hour),
		cadenceEntry("host-a", 80, "opened", 2*time.Hour),
		cadenceEntry("host-a", 80, "opened", 3*time.Hour),
	}
	res := AnalyseCadence(entries)
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	if !res[0].Regular {
		t.Errorf("expected port 80 to be marked regular")
	}
	if res[0].Occurrences != 4 {
		t.Errorf("expected 4 occurrences, got %d", res[0].Occurrences)
	}
	if res[0].AvgInterval != time.Hour {
		t.Errorf("expected avg interval 1h, got %v", res[0].AvgInterval)
	}
}

func TestAnalyseCadence_IrregularPort(t *testing.T) {
	// Highly irregular gaps.
	entries := []Entry{
		cadenceEntry("host-b", 443, "opened", 0),
		cadenceEntry("host-b", 443, "opened", 10*time.Minute),
		cadenceEntry("host-b", 443, "opened", 5*time.Hour),
		cadenceEntry("host-b", 443, "opened", 5*time.Hour+11*time.Minute),
	}
	res := AnalyseCadence(entries)
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	if res[0].Regular {
		t.Errorf("expected port 443 to be marked irregular")
	}
}

func TestAnalyseCadence_SkipsNonOpened(t *testing.T) {
	entries := []Entry{
		cadenceEntry("host-c", 22, "closed", 0),
		cadenceEntry("host-c", 22, "closed", time.Hour),
		cadenceEntry("host-c", 22, "scan", 2*time.Hour),
	}
	res := AnalyseCadence(entries)
	if len(res) != 0 {
		t.Errorf("expected no results for non-opened events, got %d", len(res))
	}
}

func TestAnalyseCadence_SingleOccurrenceExcluded(t *testing.T) {
	entries := []Entry{
		cadenceEntry("host-d", 8080, "opened", 0),
	}
	res := AnalyseCadence(entries)
	if len(res) != 0 {
		t.Errorf("expected no results for single occurrence, got %d", len(res))
	}
}

func TestAnalyseCadence_MultipleHosts(t *testing.T) {
	entries := []Entry{
		cadenceEntry("alpha", 80, "opened", 0),
		cadenceEntry("alpha", 80, "opened", time.Hour),
		cadenceEntry("beta", 80, "opened", 0),
		cadenceEntry("beta", 80, "opened", 2*time.Hour),
	}
	res := AnalyseCadence(entries)
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
	if res[0].Host != "alpha" || res[1].Host != "beta" {
		t.Errorf("unexpected host order: %v, %v", res[0].Host, res[1].Host)
	}
}
