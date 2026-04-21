package history

import (
	"testing"
	"time"
)

func riskEntry(host string, ports []int, event EventType) Entry {
	return Entry{
		Host:      host,
		Ports:     ports,
		Event:     event,
		Timestamp: time.Now(),
	}
}

func TestAssessRisk_HighRisk(t *testing.T) {
	entries := []Entry{
		riskEntry("host-a", []int{80, 443, 8080}, EventOpened),
		riskEntry("host-a", []int{22, 3306}, EventOpened),
		riskEntry("host-a", []int{80}, EventClosed),
		riskEntry("host-a", []int{9200}, EventOpened),
	}

	reports := AssessRisk(entries)
	if len(reports) == 0 {
		t.Fatal("expected reports")
	}
	if reports[0].Host != "host-a" {
		t.Errorf("expected host-a, got %s", reports[0].Host)
	}
	if reports[0].Level != RiskHigh {
		t.Errorf("expected high risk, got %s", reports[0].Level)
	}
}

func TestAssessRisk_LowRisk(t *testing.T) {
	entries := []Entry{
		riskEntry("host-b", []int{80}, EventOpened),
	}

	reports := AssessRisk(entries)
	if len(reports) == 0 {
		t.Fatal("expected reports")
	}
	if reports[0].Level != RiskLow {
		t.Errorf("expected low risk, got %s", reports[0].Level)
	}
}

func TestAssessRisk_OrderedByScore(t *testing.T) {
	entries := []Entry{
		riskEntry("host-low", []int{80}, EventOpened),
		riskEntry("host-high", []int{22, 80, 443, 8080, 3306}, EventOpened),
		riskEntry("host-high", []int{22}, EventOpened),
		riskEntry("host-high", []int{80}, EventClosed),
		riskEntry("host-high", []int{9200}, EventOpened),
	}

	reports := AssessRisk(entries)
	if len(reports) < 2 {
		t.Fatal("expected at least 2 reports")
	}
	if reports[0].Score < reports[1].Score {
		t.Errorf("reports not sorted by score descending")
	}
}

func TestAssessRisk_SkipsNoChangeEvent(t *testing.T) {
	entries := []Entry{
		riskEntry("host-c", []int{80}, EventNoChange),
	}

	reports := AssessRisk(entries)
	if len(reports) != 0 {
		t.Errorf("expected no reports for no-change events, got %d", len(reports))
	}
}

func TestAssessRisk_Empty(t *testing.T) {
	reports := AssessRisk(nil)
	if len(reports) != 0 {
		t.Errorf("expected empty result, got %d", len(reports))
	}
}
