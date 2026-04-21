package history

import (
	"testing"
	"time"
)

var outlierBase = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func outlierEntry(host string, ports []int, event EventType, offset int) Entry {
	return Entry{
		Host:      host,
		Ports:     ports,
		Event:     event,
		Timestamp: outlierBase.Add(time.Duration(offset) * time.Minute),
	}
}

func TestDetectOutliers_FlagsHighPortHost(t *testing.T) {
	entries := []Entry{
		outlierEntry("host-a", []int{80, 443}, EventScan, 0),
		outlierEntry("host-b", []int{80, 443}, EventScan, 0),
		outlierEntry("host-c", []int{80, 443, 8080, 8443, 9000, 9001, 9002, 9003}, EventScan, 0),
	}

	results := DetectOutliers(entries, 1.0)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	// host-c should be the top outlier.
	if results[0].Host != "host-c" {
		t.Errorf("expected host-c as top outlier, got %s", results[0].Host)
	}
	if !results[0].IsOutlier {
		t.Errorf("expected host-c to be flagged as outlier")
	}
}

func TestDetectOutliers_NoOutliersWhenUniform(t *testing.T) {
	entries := []Entry{
		outlierEntry("host-a", []int{80, 443}, EventScan, 0),
		outlierEntry("host-b", []int{80, 443}, EventScan, 0),
		outlierEntry("host-c", []int{80, 443}, EventScan, 0),
	}

	results := DetectOutliers(entries, 2.0)
	for _, r := range results {
		if r.IsOutlier {
			t.Errorf("expected no outliers for uniform data, got %s", r.Host)
		}
	}
}

func TestDetectOutliers_SkipsNonScanEvents(t *testing.T) {
	entries := []Entry{
		outlierEntry("host-a", []int{80}, EventScan, 0),
		outlierEntry("host-a", []int{80, 443, 8080}, EventOpened, 1),
	}

	results := DetectOutliers(entries, 1.0)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].OpenPorts != 1 {
		t.Errorf("expected 1 open port from scan only, got %d", results[0].OpenPorts)
	}
}

func TestDetectOutliers_Empty(t *testing.T) {
	results := DetectOutliers(nil, 2.0)
	if results != nil {
		t.Errorf("expected nil for empty input, got %v", results)
	}
}

func TestDetectOutliers_DefaultThreshold(t *testing.T) {
	entries := []Entry{
		outlierEntry("host-a", []int{80}, EventScan, 0),
		outlierEntry("host-b", []int{80, 443, 8080, 8443, 9000}, EventScan, 0),
	}

	// threshold=0 should default to 2.0 (no panic)
	results := DetectOutliers(entries, 0)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}
