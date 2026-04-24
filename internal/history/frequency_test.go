package history

import (
	"testing"
	"time"
)

func freqEntry(host string, event EventType, ports []int, t time.Time) Entry {
	return Entry{Host: host, Event: event, Ports: ports, Timestamp: t}
}

func TestAnalyseFrequency_BasicCounts(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		freqEntry("a", EventOpened, []int{80, 443}, now),
		freqEntry("b", EventOpened, []int{80}, now),
		freqEntry("a", EventClosed, []int{80}, now),
	}

	result := AnalyseFrequency(entries)

	if len(result) == 0 {
		t.Fatal("expected results, got none")
	}

	// port 80 should be first: 2 opens
	if result[0].Port != 80 {
		t.Errorf("expected port 80 first, got %d", result[0].Port)
	}
	if result[0].OpenCount != 2 {
		t.Errorf("expected OpenCount=2 for port 80, got %d", result[0].OpenCount)
	}
	if result[0].CloseCount != 1 {
		t.Errorf("expected CloseCount=1 for port 80, got %d", result[0].CloseCount)
	}
}

func TestAnalyseFrequency_SkipsNoChange(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		freqEntry("a", EventNoChange, []int{22}, now),
		freqEntry("a", EventOpened, []int{443}, now),
	}

	result := AnalyseFrequency(entries)

	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if result[0].Port != 443 {
		t.Errorf("expected port 443, got %d", result[0].Port)
	}
}

func TestAnalyseFrequency_Empty(t *testing.T) {
	result := AnalyseFrequency(nil)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestAnalyseFrequency_HostsDeduped(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		freqEntry("a", EventOpened, []int{8080}, now),
		freqEntry("a", EventOpened, []int{8080}, now.Add(time.Minute)),
		freqEntry("b", EventOpened, []int{8080}, now),
	}

	result := AnalyseFrequency(entries)

	if len(result) != 1 {
		t.Fatalf("expected 1 port, got %d", len(result))
	}
	if len(result[0].Hosts) != 2 {
		t.Errorf("expected 2 unique hosts, got %d", len(result[0].Hosts))
	}
	if result[0].OpenCount != 3 {
		t.Errorf("expected OpenCount=3, got %d", result[0].OpenCount)
	}
}
