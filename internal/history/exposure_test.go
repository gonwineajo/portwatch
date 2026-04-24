package history

import (
	"testing"
	"time"
)

var expBase = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func expEntry(host string, event EventType, ports []int, offsetMin int) Entry {
	return Entry{
		Host:      host,
		Event:     event,
		Ports:     ports,
		Timestamp: expBase.Add(time.Duration(offsetMin) * time.Minute),
	}
}

func TestAnalyseExposure_BasicDuration(t *testing.T) {
	entries := []Entry{
		expEntry("host-a", EventOpened, []int{80}, 0),
		expEntry("host-a", EventClosed, []int{80}, 60),
	}
	now := expBase.Add(90 * time.Minute)
	res := AnalyseExposure(entries, now)

	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	if res[0].Duration != 60*time.Minute {
		t.Errorf("expected 60m, got %v", res[0].Duration)
	}
	if res[0].StillOpen {
		t.Error("expected port to be closed")
	}
}

func TestAnalyseExposure_StillOpen(t *testing.T) {
	entries := []Entry{
		expEntry("host-b", EventOpened, []int{443}, 0),
	}
	now := expBase.Add(30 * time.Minute)
	res := AnalyseExposure(entries, now)

	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	if !res[0].StillOpen {
		t.Error("expected port to still be open")
	}
	if res[0].Duration != 30*time.Minute {
		t.Errorf("expected 30m, got %v", res[0].Duration)
	}
}

func TestAnalyseExposure_SkipsNoChange(t *testing.T) {
	entries := []Entry{
		expEntry("host-c", EventNoChange, []int{22}, 0),
		expEntry("host-c", EventNoChange, []int{22}, 10),
	}
	now := expBase.Add(20 * time.Minute)
	res := AnalyseExposure(entries, now)
	if len(res) != 0 {
		t.Errorf("expected 0 results, got %d", len(res))
	}
}

func TestAnalyseExposure_MultipleHosts(t *testing.T) {
	entries := []Entry{
		expEntry("host-a", EventOpened, []int{80}, 0),
		expEntry("host-b", EventOpened, []int{80}, 0),
		expEntry("host-a", EventClosed, []int{80}, 120),
	}
	now := expBase.Add(120 * time.Minute)
	res := AnalyseExposure(entries, now)

	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
	// host-a closed at 120m, host-b still open at 120m — equal duration, sorted by host.
	if res[0].Duration < res[1].Duration {
		t.Error("results should be sorted descending by duration")
	}
}

func TestAnalyseExposure_Empty(t *testing.T) {
	res := AnalyseExposure(nil, expBase)
	if len(res) != 0 {
		t.Errorf("expected empty, got %d", len(res))
	}
}
