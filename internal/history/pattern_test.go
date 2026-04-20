package history

import (
	"testing"
	"time"
)

func patternEntry(host string, opened, closed []int, offset time.Duration) Entry {
	return Entry{
		Timestamp:   time.Now().Add(offset),
		Host:        host,
		Event:       "scan",
		OpenedPorts: opened,
		ClosedPorts: closed,
	}
}

func TestDetectPatterns_RecurringPort(t *testing.T) {
	entries := []Entry{
		patternEntry("host-a", []int{80}, []int{}, 0),
		patternEntry("host-a", []int{}, []int{80}, -1*time.Hour),
		patternEntry("host-a", []int{80}, []int{}, -2*time.Hour),
		patternEntry("host-a", []int{}, []int{80}, -3*time.Hour),
	}

	results := DetectPatterns(entries)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if !r.Recurring {
		t.Errorf("expected port 80 on host-a to be recurring")
	}
	if r.OpenCount != 2 || r.CloseCount != 2 {
		t.Errorf("expected open=2 close=2, got open=%d close=%d", r.OpenCount, r.CloseCount)
	}
}

func TestDetectPatterns_NonRecurring(t *testing.T) {
	entries := []Entry{
		patternEntry("host-b", []int{443}, []int{}, 0),
	}

	results := DetectPatterns(entries)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Recurring {
		t.Errorf("expected port 443 to not be recurring")
	}
}

func TestDetectPatterns_Empty(t *testing.T) {
	results := DetectPatterns(nil)
	if len(results) != 0 {
		t.Errorf("expected empty results for nil input")
	}
}

func TestRecurringOnly_Filters(t *testing.T) {
	patterns := []PatternResult{
		{Host: "h1", Port: 80, OpenCount: 3, CloseCount: 3, Recurring: true},
		{Host: "h1", Port: 443, OpenCount: 1, CloseCount: 0, Recurring: false},
		{Host: "h2", Port: 22, OpenCount: 2, CloseCount: 2, Recurring: true},
	}

	out := RecurringOnly(patterns)
	if len(out) != 2 {
		t.Fatalf("expected 2 recurring, got %d", len(out))
	}
	for _, r := range out {
		if !r.Recurring {
			t.Errorf("RecurringOnly returned non-recurring entry: %+v", r)
		}
	}
}

func TestDetectPatterns_MultipleHosts(t *testing.T) {
	entries := []Entry{
		patternEntry("alpha", []int{8080}, []int{}, 0),
		patternEntry("beta", []int{8080}, []int{}, 0),
		patternEntry("alpha", []int{}, []int{8080}, -1*time.Hour),
		patternEntry("alpha", []int{8080}, []int{}, -2*time.Hour),
		patternEntry("alpha", []int{}, []int{8080}, -3*time.Hour),
	}

	results := DetectPatterns(entries)
	recurring := RecurringOnly(results)
	if len(recurring) != 1 {
		t.Fatalf("expected 1 recurring pattern, got %d", len(recurring))
	}
	if recurring[0].Host != "alpha" {
		t.Errorf("expected recurring host to be alpha, got %s", recurring[0].Host)
	}
}
