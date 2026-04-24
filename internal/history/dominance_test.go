package history

import (
	"testing"
	"time"
)

func domEntry(host, event string, ports []int, minutesAgo int) Entry {
	return Entry{
		Host:      host,
		Event:     event,
		Ports:     ports,
		Timestamp: time.Now().Add(-time.Duration(minutesAgo) * time.Minute),
	}
}

func TestAnalyseDominance_BasicRanking(t *testing.T) {
	entries := []Entry{
		domEntry("host-a", "opened", []int{80, 443}, 10),
		domEntry("host-b", "opened", []int{80, 443}, 9),
		domEntry("host-c", "opened", []int{80}, 8),
		domEntry("host-d", "opened", []int{443}, 7),
	}

	results := AnalyseDominance(entries, 1)
	if len(results) == 0 {
		t.Fatal("expected results, got none")
	}
	// port 80 opened on 3 hosts, port 443 on 3 hosts — both score equally
	// just verify both appear
	ports := make(map[int]bool)
	for _, r := range results {
		ports[r.Port] = true
	}
	if !ports[80] || !ports[443] {
		t.Errorf("expected ports 80 and 443 in results, got %v", results)
	}
}

func TestAnalyseDominance_MinHostsFilter(t *testing.T) {
	entries := []Entry{
		domEntry("host-a", "opened", []int{22, 80}, 5),
		domEntry("host-b", "opened", []int{80}, 4),
	}

	// port 22 only on 1 host, port 80 on 2
	results := AnalyseDominance(entries, 2)
	for _, r := range results {
		if r.Port == 22 {
			t.Errorf("port 22 should have been filtered out by minHosts=2")
		}
	}
	if len(results) != 1 || results[0].Port != 80 {
		t.Errorf("expected only port 80, got %v", results)
	}
}

func TestAnalyseDominance_SkipsNonOpened(t *testing.T) {
	entries := []Entry{
		domEntry("host-a", "closed", []int{8080}, 5),
		domEntry("host-b", "no_change", []int{8080}, 4),
		domEntry("host-a", "opened", []int{443}, 3),
	}

	results := AnalyseDominance(entries, 1)
	for _, r := range results {
		if r.Port == 8080 {
			t.Errorf("port 8080 should not appear — only closed/no_change events")
		}
	}
}

func TestAnalyseDominance_Empty(t *testing.T) {
	results := AnalyseDominance(nil, 1)
	if len(results) != 0 {
		t.Errorf("expected empty results for nil input")
	}
}

func TestAnalyseDominance_ScoreReflectsFrequency(t *testing.T) {
	entries := []Entry{
		domEntry("host-a", "opened", []int{80}, 10),
		domEntry("host-a", "opened", []int{80}, 8),
		domEntry("host-a", "opened", []int{80}, 6),
		domEntry("host-b", "opened", []int{443}, 5),
	}

	results := AnalyseDominance(entries, 1)
	var r80, r443 *DominanceResult
	for i := range results {
		if results[i].Port == 80 {
			r80 = &results[i]
		}
		if results[i].Port == 443 {
			r443 = &results[i]
		}
	}
	if r80 == nil || r443 == nil {
		t.Fatal("expected both ports in results")
	}
	if r80.TotalOpen != 3 {
		t.Errorf("expected TotalOpen=3 for port 80, got %d", r80.TotalOpen)
	}
	if r80.Score <= r443.Score {
		t.Errorf("port 80 should score higher than port 443")
	}
}
