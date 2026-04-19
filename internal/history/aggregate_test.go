package history

import (
	"testing"
	"time"
)

var aggEntries = []Entry{
	{Host: "host-a", Event: "opened", Ports: []int{80, 443}, Timestamp: time.Now()},
	{Host: "host-a", Event: "closed", Ports: []int{22}, Timestamp: time.Now()},
	{Host: "host-b", Event: "opened", Ports: []int{80}, Timestamp: time.Now()},
	{Host: "host-b", Event: "opened", Ports: []int{8080}, Timestamp: time.Now()},
	{Host: "host-c", Event: "closed", Ports: []int{443}, Timestamp: time.Now()},
}

func TestAggregateByHost_Counts(t *testing.T) {
	result := AggregateByHost(aggEntries)
	if len(result) != 3 {
		t.Fatalf("expected 3 hosts, got %d", len(result))
	}
	// host-b has total=2, should be first
	if result[0].Host != "host-b" {
		t.Errorf("expected host-b first, got %s", result[0].Host)
	}
	if result[0].Opened != 2 || result[0].Closed != 0 {
		t.Errorf("unexpected counts for host-b: %+v", result[0])
	}
}

func TestAggregateByHost_HostA(t *testing.T) {
	result := AggregateByHost(aggEntries)
	var ha HostSummary
	for _, r := range result {
		if r.Host == "host-a" {
			ha = r
		}
	}
	if ha.Opened != 1 || ha.Closed != 1 || ha.Total != 2 {
		t.Errorf("unexpected host-a summary: %+v", ha)
	}
}

func TestAggregateByPort_Counts(t *testing.T) {
	result := AggregateByPort(aggEntries)
	if len(result) == 0 {
		t.Fatal("expected port frequencies")
	}
	// port 80 appears in host-a and host-b => count 2
	var freq80 int
	for _, r := range result {
		if r.Port == 80 {
			freq80 = r.Count
		}
	}
	if freq80 != 2 {
		t.Errorf("expected port 80 count=2, got %d", freq80)
	}
}

func TestAggregateByPort_Empty(t *testing.T) {
	result := AggregateByPort([]Entry{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}

func TestAggregateByHost_Empty(t *testing.T) {
	result := AggregateByHost([]Entry{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}
