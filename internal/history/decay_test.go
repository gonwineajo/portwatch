package history

import (
	"testing"
	"time"
)

var decayBase = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func decayEntry(host string, event EventType, hoursAgo float64) Entry {
	return Entry{
		Host:      host,
		Event:     event,
		Timestamp: decayBase.Add(-time.Duration(hoursAgo * float64(time.Hour))),
		Ports:     []int{80},
	}
}

func TestDecayScores_RecentHigherThanOld(t *testing.T) {
	entries := []Entry{
		decayEntry("recent", EventOpened, 1),
		decayEntry("old", EventOpened, 48),
	}
	res := DecayScores(entries, decayBase, 24*time.Hour)
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
	if res[0].Host != "recent" {
		t.Errorf("expected 'recent' first, got %q", res[0].Host)
	}
	if res[0].Score <= res[1].Score {
		t.Errorf("recent score %f should exceed old score %f", res[0].Score, res[1].Score)
	}
}

func TestDecayScores_SkipsNoChange(t *testing.T) {
	entries := []Entry{
		decayEntry("a", EventNoChange, 1),
		decayEntry("b", EventOpened, 1),
	}
	res := DecayScores(entries, decayBase, 24*time.Hour)
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	if res[0].Host != "b" {
		t.Errorf("expected host 'b', got %q", res[0].Host)
	}
}

func TestDecayScores_Empty(t *testing.T) {
	res := DecayScores(nil, decayBase, 24*time.Hour)
	if len(res) != 0 {
		t.Errorf("expected empty result, got %d", len(res))
	}
}

func TestDecayScores_MultipleEventsAccumulate(t *testing.T) {
	entries := []Entry{
		decayEntry("busy", EventOpened, 2),
		decayEntry("busy", EventClosed, 3),
		decayEntry("quiet", EventOpened, 2),
	}
	res := DecayScores(entries, decayBase, 24*time.Hour)
	if res[0].Host != "busy" {
		t.Errorf("expected 'busy' first, got %q", res[0].Host)
	}
	if res[0].EventCount != 2 {
		t.Errorf("expected EventCount=2 for busy, got %d", res[0].EventCount)
	}
}

func TestDecayScores_DefaultHalfLife(t *testing.T) {
	entries := []Entry{
		decayEntry("x", EventOpened, 1),
	}
	// zero halfLife should use default (24h) without panic
	res := DecayScores(entries, decayBase, 0)
	if len(res) != 1 {
		t.Fatalf("expected 1 result")
	}
	if res[0].Score <= 0 {
		t.Errorf("expected positive score, got %f", res[0].Score)
	}
}
