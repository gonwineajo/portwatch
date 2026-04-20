package history

import (
	"testing"
	"time"
)

func scoreEntries() []Entry {
	now := time.Now()
	return []Entry{
		{Host: "host-a", Event: EventOpened, Ports: []int{80, 443}, Timestamp: now},
		{Host: "host-a", Event: EventOpened, Ports: []int{8080}, Timestamp: now.Add(time.Minute)},
		{Host: "host-a", Event: EventClosed, Ports: []int{80}, Timestamp: now.Add(2 * time.Minute)},
		{Host: "host-b", Event: EventOpened, Ports: []int{22}, Timestamp: now},
		{Host: "host-c", Event: EventClosed, Ports: []int{443}, Timestamp: now},
	}
}

func TestScoreHosts_Order(t *testing.T) {
	results := ScoreHosts(scoreEntries())
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0].Host != "host-a" {
		t.Errorf("expected host-a first, got %s", results[0].Host)
	}
}

func TestScoreHosts_HostA(t *testing.T) {
	results := ScoreHosts(scoreEntries())
	var ha RiskScore
	for _, r := range results {
		if r.Host == "host-a" {
			ha = r
			break
		}
	}
	// openCount=2 (*3=6), closeCount=1 (*1=1), uniquePorts=3 (*2=6) => 13
	if ha.OpenCount != 2 {
		t.Errorf("expected OpenCount 2, got %d", ha.OpenCount)
	}
	if ha.CloseCount != 1 {
		t.Errorf("expected CloseCount 1, got %d", ha.CloseCount)
	}
	if ha.UniquePortsOpened != 3 {
		t.Errorf("expected UniquePortsOpened 3, got %d", ha.UniquePortsOpened)
	}
	if ha.Score != 13 {
		t.Errorf("expected Score 13, got %d", ha.Score)
	}
}

func TestScoreHosts_HostC_OnlyClosed(t *testing.T) {
	results := ScoreHosts(scoreEntries())
	var hc RiskScore
	for _, r := range results {
		if r.Host == "host-c" {
			hc = r
			break
		}
	}
	if hc.OpenCount != 0 {
		t.Errorf("expected OpenCount 0, got %d", hc.OpenCount)
	}
	if hc.UniquePortsOpened != 0 {
		t.Errorf("expected UniquePortsOpened 0, got %d", hc.UniquePortsOpened)
	}
	// closeCount=1 => score=1
	if hc.Score != 1 {
		t.Errorf("expected Score 1, got %d", hc.Score)
	}
}

func TestScoreHosts_Empty(t *testing.T) {
	results := ScoreHosts([]Entry{})
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}
