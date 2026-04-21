package history

import (
	"testing"
)

func TestLifecycle_OpenClose(t *testing.T) {
	entries := []Entry{
		chainEntry("host-a", EventOpened, []int{80}, nil, 1000),
		chainEntry("host-a", EventClosed, nil, []int{80}, 1500),
	}
	chains := BuildChains(entries)
	summaries := Lifecycle(chains)
	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(summaries))
	}
	s := summaries[0]
	if s.OpenCount != 1 {
		t.Errorf("expected OpenCount=1, got %d", s.OpenCount)
	}
	if s.CloseCount != 1 {
		t.Errorf("expected CloseCount=1, got %d", s.CloseCount)
	}
	if s.TotalOpenTime != 500 {
		t.Errorf("expected TotalOpenTime=500, got %d", s.TotalOpenTime)
	}
	if s.CurrentlyOpen {
		t.Error("expected CurrentlyOpen=false")
	}
}

func TestLifecycle_StillOpen(t *testing.T) {
	entries := []Entry{
		chainEntry("host-a", EventOpened, []int{443}, nil, 2000),
	}
	chains := BuildChains(entries)
	summaries := Lifecycle(chains)
	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(summaries))
	}
	if !summaries[0].CurrentlyOpen {
		t.Error("expected CurrentlyOpen=true")
	}
	if summaries[0].TotalOpenTime != 0 {
		t.Errorf("expected TotalOpenTime=0, got %d", summaries[0].TotalOpenTime)
	}
}

func TestLifecycle_MultipleOpenClose(t *testing.T) {
	entries := []Entry{
		chainEntry("host-a", EventOpened, []int{22}, nil, 1000),
		chainEntry("host-a", EventClosed, nil, []int{22}, 1200),
		chainEntry("host-a", EventOpened, []int{22}, nil, 2000),
		chainEntry("host-a", EventClosed, nil, []int{22}, 2300),
	}
	chains := BuildChains(entries)
	summaries := Lifecycle(chains)
	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(summaries))
	}
	s := summaries[0]
	if s.OpenCount != 2 {
		t.Errorf("expected OpenCount=2, got %d", s.OpenCount)
	}
	if s.TotalOpenTime != 500 {
		t.Errorf("expected TotalOpenTime=500, got %d", s.TotalOpenTime)
	}
}

func TestLongestOpen_ReturnsMax(t *testing.T) {
	summaries := []LifecycleSummary{
		{Host: "host-a", Port: 80, TotalOpenTime: 100},
		{Host: "host-a", Port: 443, TotalOpenTime: 900},
		{Host: "host-b", Port: 22, TotalOpenTime: 300},
	}
	best := LongestOpen(summaries)
	if best == nil {
		t.Fatal("expected non-nil result")
	}
	if best.Port != 443 {
		t.Errorf("expected port 443, got %d", best.Port)
	}
}

func TestLongestOpen_Empty(t *testing.T) {
	result := LongestOpen(nil)
	if result != nil {
		t.Error("expected nil for empty input")
	}
}
