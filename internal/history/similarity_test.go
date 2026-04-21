package history

import (
	"testing"
	"time"
)

var simBase = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func simEntry(host string, event EventType, ports []int, offset int) Entry {
	return Entry{
		Host:      host,
		Event:     event,
		Ports:     ports,
		Timestamp: simBase.Add(time.Duration(offset) * time.Minute),
	}
}

func TestComputeSimilarity_FullOverlap(t *testing.T) {
	entries := []Entry{
		simEntry("a", EventScan, []int{80, 443}, 0),
		simEntry("b", EventScan, []int{80, 443}, 1),
	}
	r := ComputeSimilarity(entries, "a", "b")
	if r.Jaccard != 1.0 {
		t.Fatalf("expected jaccard 1.0, got %f", r.Jaccard)
	}
	if len(r.Common) != 2 {
		t.Fatalf("expected 2 common ports, got %d", len(r.Common))
	}
	if len(r.OnlyA) != 0 || len(r.OnlyB) != 0 {
		t.Fatal("expected no exclusive ports")
	}
}

func TestComputeSimilarity_NoOverlap(t *testing.T) {
	entries := []Entry{
		simEntry("a", EventScan, []int{22, 80}, 0),
		simEntry("b", EventScan, []int{443, 8080}, 1),
	}
	r := ComputeSimilarity(entries, "a", "b")
	if r.Jaccard != 0.0 {
		t.Fatalf("expected jaccard 0.0, got %f", r.Jaccard)
	}
	if len(r.Common) != 0 {
		t.Fatalf("expected 0 common ports, got %d", len(r.Common))
	}
}

func TestComputeSimilarity_PartialOverlap(t *testing.T) {
	entries := []Entry{
		simEntry("a", EventScan, []int{80, 443, 22}, 0),
		simEntry("b", EventScan, []int{80, 443, 8080}, 1),
	}
	r := ComputeSimilarity(entries, "a", "b")
	// common={80,443}, onlyA={22}, onlyB={8080} => 2/4 = 0.5
	if r.Jaccard != 0.5 {
		t.Fatalf("expected jaccard 0.5, got %f", r.Jaccard)
	}
}

func TestComputeSimilarity_UsesLatestScan(t *testing.T) {
	entries := []Entry{
		simEntry("a", EventScan, []int{22}, 0),
		simEntry("a", EventScan, []int{80, 443}, 2),
		simEntry("b", EventScan, []int{80, 443}, 1),
	}
	r := ComputeSimilarity(entries, "a", "b")
	if r.Jaccard != 1.0 {
		t.Fatalf("expected latest scan used, jaccard 1.0, got %f", r.Jaccard)
	}
}

func TestAllPairSimilarity_OrderedByJaccard(t *testing.T) {
	entries := []Entry{
		simEntry("a", EventScan, []int{80, 443}, 0),
		simEntry("b", EventScan, []int{80, 443}, 1),
		simEntry("c", EventScan, []int{22}, 2),
	}
	results := AllPairSimilarity(entries)
	if len(results) != 3 {
		t.Fatalf("expected 3 pairs, got %d", len(results))
	}
	if results[0].Jaccard < results[1].Jaccard {
		t.Fatal("results should be sorted descending by jaccard")
	}
	if results[0].HostA != "a" || results[0].HostB != "b" {
		t.Fatalf("top pair should be a/b, got %s/%s", results[0].HostA, results[0].HostB)
	}
}

func TestComputeSimilarity_EmptyHosts(t *testing.T) {
	r := ComputeSimilarity([]Entry{}, "a", "b")
	if r.Jaccard != 0.0 {
		t.Fatalf("expected 0 jaccard for empty entries, got %f", r.Jaccard)
	}
}
