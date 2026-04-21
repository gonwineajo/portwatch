package history

import (
	"testing"
	"time"
)

var corrBase = time.Now()

func corrEntry(host string, event EventType, ports []int, offset int) Entry {
	return Entry{
		Host:      host,
		Event:     event,
		Ports:     ports,
		Timestamp: corrBase.Add(time.Duration(offset) * time.Minute),
	}
}

func TestCorrelateOpenPorts_BasicPair(t *testing.T) {
	entries := []Entry{
		corrEntry("host-a", EventScan, []int{80, 443}, 0),
		corrEntry("host-a", EventScan, []int{80, 443}, 1),
		corrEntry("host-a", EventScan, []int{80, 443}, 2),
	}
	result := CorrelateOpenPorts(entries, 2)
	if len(result) != 1 {
		t.Fatalf("expected 1 correlation, got %d", len(result))
	}
	if result[0].PortA != 80 || result[0].PortB != 443 {
		t.Errorf("unexpected pair: %d/%d", result[0].PortA, result[0].PortB)
	}
	if result[0].CoOccurrences != 3 {
		t.Errorf("expected 3 co-occurrences, got %d", result[0].CoOccurrences)
	}
}

func TestCorrelateOpenPorts_BelowMinCount(t *testing.T) {
	entries := []Entry{
		corrEntry("host-a", EventScan, []int{22, 8080}, 0),
	}
	result := CorrelateOpenPorts(entries, 3)
	if len(result) != 0 {
		t.Errorf("expected no results below minCount, got %d", len(result))
	}
}

func TestCorrelateOpenPorts_SkipsClosedEvent(t *testing.T) {
	entries := []Entry{
		corrEntry("host-a", EventClosed, []int{80, 443}, 0),
		corrEntry("host-a", EventClosed, []int{80, 443}, 1),
		corrEntry("host-a", EventClosed, []int{80, 443}, 2),
	}
	result := CorrelateOpenPorts(entries, 1)
	if len(result) != 0 {
		t.Errorf("expected closed events to be skipped, got %d results", len(result))
	}
}

func TestCorrelateOpenPorts_MultipleHosts(t *testing.T) {
	entries := []Entry{
		corrEntry("host-a", EventScan, []int{80, 443}, 0),
		corrEntry("host-b", EventScan, []int{80, 443}, 1),
		corrEntry("host-a", EventScan, []int{80, 443}, 2),
	}
	result := CorrelateOpenPorts(entries, 2)
	if len(result) != 1 {
		t.Fatalf("expected 1 correlation, got %d", len(result))
	}
	if len(result[0].Hosts) != 2 {
		t.Errorf("expected 2 hosts, got %d", len(result[0].Hosts))
	}
}

func TestCorrelateOpenPorts_Empty(t *testing.T) {
	result := CorrelateOpenPorts([]Entry{}, 1)
	if len(result) != 0 {
		t.Errorf("expected empty result for empty input")
	}
}
