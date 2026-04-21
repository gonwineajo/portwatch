package history

import (
	"testing"
	"time"
)

var anomalyBase = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func anomalyEntry(host string, opened, closed []int, offset int) Entry {
	e := Entry{
		Host:        host,
		ScannedAt:   anomalyBase.Add(time.Duration(offset) * time.Hour),
		OpenedPorts: opened,
		ClosedPorts: closed,
	}
	switch {
	case len(opened) > 0:
		e.Event = EventOpened
	case len(closed) > 0:
		e.Event = EventClosed
	default:
		e.Event = EventNoChange
	}
	return e
}

func TestDetectAnomalies_RareOpened(t *testing.T) {
	entries := []Entry{
		anomalyEntry("host-a", []int{9999}, nil, 0), // only once → anomaly
		anomalyEntry("host-a", []int{80}, nil, 1),
		anomalyEntry("host-a", []int{80}, nil, 2), // twice → not anomalous
	}
	reports := DetectAnomalies(entries, 2)
	if len(reports) != 1 {
		t.Fatalf("expected 1 anomaly, got %d", len(reports))
	}
	if reports[0].Port != 9999 {
		t.Errorf("expected port 9999, got %d", reports[0].Port)
	}
	if reports[0].Event != EventOpened {
		t.Errorf("expected event opened, got %s", reports[0].Event)
	}
}

func TestDetectAnomalies_SkipsNoChange(t *testing.T) {
	entries := []Entry{
		anomalyEntry("host-a", nil, nil, 0),
		anomalyEntry("host-a", nil, nil, 1),
	}
	reports := DetectAnomalies(entries, 2)
	if len(reports) != 0 {
		t.Errorf("expected no anomalies for no-change events, got %d", len(reports))
	}
}

func TestDetectAnomalies_Empty(t *testing.T) {
	reports := DetectAnomalies(nil, 2)
	if reports != nil && len(reports) != 0 {
		t.Errorf("expected empty result, got %v", reports)
	}
}

func TestAnomaliesByHost_Groups(t *testing.T) {
	reports := []AnomalyReport{
		{Host: "host-a", Port: 9999, Event: EventOpened, OccurredAt: anomalyBase},
		{Host: "host-b", Port: 8080, Event: EventClosed, OccurredAt: anomalyBase},
		{Host: "host-a", Port: 1234, Event: EventOpened, OccurredAt: anomalyBase},
	}
	byHost := AnomaliesByHost(reports)
	if len(byHost["host-a"]) != 2 {
		t.Errorf("expected 2 anomalies for host-a, got %d", len(byHost["host-a"]))
	}
	if len(byHost["host-b"]) != 1 {
		t.Errorf("expected 1 anomaly for host-b, got %d", len(byHost["host-b"]))
	}
}

func TestDetectAnomalies_DefaultMinCount(t *testing.T) {
	entries := []Entry{
		anomalyEntry("host-a", []int{443}, nil, 0),
	}
	// minCount=0 should default to 2, so single occurrence is anomalous
	reports := DetectAnomalies(entries, 0)
	if len(reports) != 1 {
		t.Errorf("expected 1 anomaly with default minCount, got %d", len(reports))
	}
}
