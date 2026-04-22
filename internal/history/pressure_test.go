package history

import (
	"testing"
	"time"
)

var pressureBase = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func pressureEntry(host string, event EventType, ports []int, offsetHours int) Entry {
	return Entry{
		Host:      host,
		EventType: event,
		Ports:     ports,
		Timestamp: pressureBase.Add(time.Duration(offsetHours) * time.Hour),
	}
}

func TestAnalysePressure_BasicFlips(t *testing.T) {
	entries := []Entry{
		pressureEntry("host-a", EventOpened, []int{80}, 0),
		pressureEntry("host-a", EventClosed, []int{80}, 2),
		pressureEntry("host-a", EventOpened, []int{80}, 4),
		pressureEntry("host-a", EventClosed, []int{80}, 6),
	}
	res := AnalysePressure(entries, 1)
	if len(res.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(res.Records))
	}
	r := res.Records[0]
	if r.Flips != 4 {
		t.Errorf("expected 4 flips, got %d", r.Flips)
	}
	if r.Port != 80 {
		t.Errorf("expected port 80, got %d", r.Port)
	}
}

func TestAnalysePressure_SkipsNoChange(t *testing.T) {
	entries := []Entry{
		pressureEntry("host-a", EventNoChange, []int{443}, 0),
		pressureEntry("host-a", EventNoChange, []int{443}, 1),
	}
	res := AnalysePressure(entries, 1)
	if len(res.Records) != 0 {
		t.Errorf("expected 0 records, got %d", len(res.Records))
	}
}

func TestAnalysePressure_MinFlipsFilter(t *testing.T) {
	entries := []Entry{
		pressureEntry("host-a", EventOpened, []int{22}, 0),
		pressureEntry("host-a", EventOpened, []int{80}, 1),
		pressureEntry("host-a", EventOpened, []int{80}, 2),
		pressureEntry("host-a", EventOpened, []int{80}, 3),
	}
	res := AnalysePressure(entries, 2)
	if len(res.Records) != 1 {
		t.Fatalf("expected 1 record after minFlips filter, got %d", len(res.Records))
	}
	if res.Records[0].Port != 80 {
		t.Errorf("expected port 80, got %d", res.Records[0].Port)
	}
}

func TestAnalysePressure_ScoreOrder(t *testing.T) {
	entries := []Entry{
		// port 22: 2 flips over 10 hours => 0.2/hr
		pressureEntry("host-a", EventOpened, []int{22}, 0),
		pressureEntry("host-a", EventClosed, []int{22}, 10),
		// port 80: 4 flips over 4 hours => 1.0/hr
		pressureEntry("host-a", EventOpened, []int{80}, 0),
		pressureEntry("host-a", EventClosed, []int{80}, 1),
		pressureEntry("host-a", EventOpened, []int{80}, 2),
		pressureEntry("host-a", EventClosed, []int{80}, 4),
	}
	res := AnalysePressure(entries, 1)
	if len(res.Records) < 2 {
		t.Fatalf("expected at least 2 records, got %d", len(res.Records))
	}
	if res.Records[0].Score <= res.Records[1].Score {
		t.Errorf("expected descending score order, got %.4f then %.4f",
			res.Records[0].Score, res.Records[1].Score)
	}
	if res.Records[0].Port != 80 {
		t.Errorf("expected port 80 first, got %d", res.Records[0].Port)
	}
}

func TestAnalysePressure_Empty(t *testing.T) {
	res := AnalysePressure([]Entry{}, 1)
	if len(res.Records) != 0 {
		t.Errorf("expected 0 records for empty input")
	}
}
