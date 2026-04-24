package history

import (
	"testing"
	"time"
)

func coEntry(host string, event EventType, ports []int, ts time.Time) Entry {
	return Entry{Host: host, Event: event, Ports: ports, Timestamp: ts}
}

var coBase = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestCoOccurrenceMatrix_BasicPair(t *testing.T) {
	entries := []Entry{
		coEntry("host-a", EventScan, []int{80, 443, 8080}, coBase),
		coEntry("host-a", EventScan, []int{80, 443}, coBase.Add(time.Hour)),
	}
	result := CoOccurrenceMatrix(entries, CoOccurrenceOptions{})
	if len(result) == 0 {
		t.Fatal("expected co-occurrence results")
	}
	// 80+443 should appear twice
	found := false
	for _, r := range result {
		if r.PortA == 80 && r.PortB == 443 {
			if r.Count != 2 {
				t.Errorf("expected count 2 for 80+443, got %d", r.Count)
			}
			found = true
		}
	}
	if !found {
		t.Error("expected pair 80+443 in result")
	}
}

func TestCoOccurrenceMatrix_SkipsNonScan(t *testing.T) {
	entries := []Entry{
		coEntry("host-a", EventOpened, []int{80, 443}, coBase),
		coEntry("host-a", EventClosed, []int{80, 443}, coBase.Add(time.Hour)),
	}
	result := CoOccurrenceMatrix(entries, CoOccurrenceOptions{})
	if len(result) != 0 {
		t.Errorf("expected no results for non-scan events, got %d", len(result))
	}
}

func TestCoOccurrenceMatrix_MinCountFilter(t *testing.T) {
	entries := []Entry{
		coEntry("host-a", EventScan, []int{80, 443}, coBase),
		coEntry("host-b", EventScan, []int{80, 8080}, coBase),
	}
	result := CoOccurrenceMatrix(entries, CoOccurrenceOptions{MinCount: 2})
	if len(result) != 0 {
		t.Errorf("expected no results above min count 2, got %d", len(result))
	}
}

func TestCoOccurrenceMatrix_HostFilter(t *testing.T) {
	entries := []Entry{
		coEntry("host-a", EventScan, []int{80, 443}, coBase),
		coEntry("host-b", EventScan, []int{80, 443}, coBase),
	}
	result := CoOccurrenceMatrix(entries, CoOccurrenceOptions{Host: "host-a"})
	if len(result) != 1 {
		t.Fatalf("expected 1 pair for host-a, got %d", len(result))
	}
	if len(result[0].Hosts) != 1 || result[0].Hosts[0] != "host-a" {
		t.Errorf("expected only host-a in Hosts, got %v", result[0].Hosts)
	}
}

func TestCoOccurrenceMatrix_OrderedByCount(t *testing.T) {
	entries := []Entry{
		coEntry("host-a", EventScan, []int{80, 443, 22}, coBase),
		coEntry("host-a", EventScan, []int{80, 443}, coBase.Add(time.Hour)),
		coEntry("host-a", EventScan, []int{80, 443}, coBase.Add(2*time.Hour)),
	}
	result := CoOccurrenceMatrix(entries, CoOccurrenceOptions{})
	if len(result) == 0 {
		t.Fatal("expected results")
	}
	if result[0].PortA != 80 || result[0].PortB != 443 {
		t.Errorf("expected 80+443 first (highest count), got %d+%d", result[0].PortA, result[0].PortB)
	}
}
