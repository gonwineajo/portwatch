package history

import (
	"testing"
	"time"
)

func baseTime() time.Time {
	return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
}

func TestSummarize_BasicCounts(t *testing.T) {
	entries := []Entry{
		{Host: "host-a", Timestamp: baseTime(), Opened: []int{80, 443}, Closed: []int{}},
		{Host: "host-a", Timestamp: baseTime().Add(time.Hour), Opened: []int{}, Closed: []int{80}},
		{Host: "host-b", Timestamp: baseTime(), Opened: []int{22}, Closed: []int{8080}},
	}

	summaries := Summarize(entries, time.Time{})

	if len(summaries) != 2 {
		t.Fatalf("expected 2 summaries, got %d", len(summaries))
	}

	sa := summaries[0] // host-a (sorted)
	if sa.Host != "host-a" || sa.Opened != 2 || sa.Closed != 1 {
		t.Errorf("host-a: unexpected summary %+v", sa)
	}
	if !sa.LastChanged.Equal(baseTime().Add(time.Hour)) {
		t.Errorf("host-a: unexpected LastChanged %v", sa.LastChanged)
	}

	sb := summaries[1] // host-b
	if sb.Host != "host-b" || sb.Opened != 1 || sb.Closed != 1 {
		t.Errorf("host-b: unexpected summary %+v", sb)
	}
}

func TestSummarize_SinceFilter(t *testing.T) {
	cutoff := baseTime().Add(30 * time.Minute)
	entries := []Entry{
		{Host: "host-a", Timestamp: baseTime(), Opened: []int{80}, Closed: []int{}},
		{Host: "host-a", Timestamp: baseTime().Add(time.Hour), Opened: []int{443}, Closed: []int{}},
	}

	summaries := Summarize(entries, cutoff)

	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(summaries))
	}
	if summaries[0].Opened != 1 {
		t.Errorf("expected 1 opened after filter, got %d", summaries[0].Opened)
	}
}

func TestSummarize_Empty(t *testing.T) {
	summaries := Summarize([]Entry{}, time.Time{})
	if len(summaries) != 0 {
		t.Errorf("expected empty summaries, got %d", len(summaries))
	}
}
