package history

import (
	"testing"
	"time"
)

var burstBase = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func burstEntry(host string, event EventType, offsetSec int, ports ...int) Entry {
	return Entry{
		Host:      host,
		Event:     event,
		Timestamp: burstBase.Add(time.Duration(offsetSec) * time.Second),
		Ports:     ports,
	}
}

func TestDetectBursts_DetectsBurst(t *testing.T) {
	entries := []Entry{
		burstEntry("host-a", EventOpened, 0, 80),
		burstEntry("host-a", EventOpened, 5, 443),
		burstEntry("host-a", EventClosed, 8, 22),
	}

	results := DetectBursts(entries, 30*time.Second, 3)
	if len(results) != 1 {
		t.Fatalf("expected 1 burst, got %d", len(results))
	}
	if results[0].Host != "host-a" {
		t.Errorf("expected host-a, got %s", results[0].Host)
	}
	if results[0].Count != 3 {
		t.Errorf("expected count 3, got %d", results[0].Count)
	}
}

func TestDetectBursts_BelowThreshold(t *testing.T) {
	entries := []Entry{
		burstEntry("host-b", EventOpened, 0, 80),
		burstEntry("host-b", EventOpened, 5, 443),
	}

	results := DetectBursts(entries, 30*time.Second, 3)
	if len(results) != 0 {
		t.Fatalf("expected 0 bursts below threshold, got %d", len(results))
	}
}

func TestDetectBursts_SkipsNoChange(t *testing.T) {
	entries := []Entry{
		burstEntry("host-c", EventNoChange, 0, 80),
		burstEntry("host-c", EventNoChange, 1, 443),
		burstEntry("host-c", EventNoChange, 2, 22),
	}

	results := DetectBursts(entries, 30*time.Second, 2)
	if len(results) != 0 {
		t.Fatalf("expected 0 bursts for no-change events, got %d", len(results))
	}
}

func TestDetectBursts_Empty(t *testing.T) {
	results := DetectBursts(nil, 30*time.Second, 2)
	if results != nil {
		t.Errorf("expected nil for empty input")
	}
}

func TestDetectBursts_MultipleHosts(t *testing.T) {
	entries := []Entry{
		burstEntry("host-x", EventOpened, 0, 80),
		burstEntry("host-x", EventClosed, 2, 22),
		burstEntry("host-y", EventOpened, 0, 9000),
		burstEntry("host-y", EventOpened, 1, 9001),
		burstEntry("host-y", EventClosed, 3, 9002),
	}

	results := DetectBursts(entries, 10*time.Second, 2)
	hosts := make(map[string]bool)
	for _, r := range results {
		hosts[r.Host] = true
	}
	if !hosts["host-x"] {
		t.Errorf("expected host-x in results")
	}
	if !hosts["host-y"] {
		t.Errorf("expected host-y in results")
	}
}
