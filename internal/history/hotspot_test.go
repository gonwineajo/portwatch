package history

import (
	"testing"
	"time"
)

func hotEntry(host, event string, ports []int, offset time.Duration) Entry {
	return Entry{
		Host:      host,
		Event:     event,
		Ports:     ports,
		Timestamp: time.Now().Add(offset),
	}
}

func TestDetectHotspots_BasicFlip(t *testing.T) {
	entries := []Entry{
		hotEntry("host-a", EventOpened, []int{80}, 0),
		hotEntry("host-a", EventClosed, []int{80}, time.Minute),
		hotEntry("host-a", EventOpened, []int{80}, 2*time.Minute),
		hotEntry("host-a", EventClosed, []int{80}, 3*time.Minute),
	}

	result := DetectHotspots(entries, 1)
	if len(result) != 1 {
		t.Fatalf("expected 1 hotspot, got %d", len(result))
	}
	if result[0].Host != "host-a" || result[0].Port != 80 {
		t.Errorf("unexpected hotspot: %+v", result[0])
	}
	if result[0].Flips != 3 {
		t.Errorf("expected 3 flips, got %d", result[0].Flips)
	}
}

func TestDetectHotspots_BelowMinFlips(t *testing.T) {
	entries := []Entry{
		hotEntry("host-b", EventOpened, []int{443}, 0),
		hotEntry("host-b", EventClosed, []int{443}, time.Minute),
	}

	result := DetectHotspots(entries, 3)
	if len(result) != 0 {
		t.Errorf("expected no hotspots above minFlips=3, got %d", len(result))
	}
}

func TestDetectHotspots_SkipsNoChange(t *testing.T) {
	entries := []Entry{
		hotEntry("host-c", EventNoChange, []int{22}, 0),
		hotEntry("host-c", EventNoChange, []int{22}, time.Minute),
	}

	result := DetectHotspots(entries, 1)
	if len(result) != 0 {
		t.Errorf("expected no hotspots for NoChange events, got %d", len(result))
	}
}

func TestDetectHotspots_MultipleHosts(t *testing.T) {
	entries := []Entry{
		hotEntry("alpha", EventOpened, []int{8080}, 0),
		hotEntry("alpha", EventClosed, []int{8080}, time.Minute),
		hotEntry("alpha", EventOpened, []int{8080}, 2*time.Minute),
		hotEntry("beta", EventOpened, []int{9090}, 0),
		hotEntry("beta", EventClosed, []int{9090}, time.Minute),
		hotEntry("beta", EventOpened, []int{9090}, 2*time.Minute),
		hotEntry("beta", EventClosed, []int{9090}, 3*time.Minute),
	}

	result := DetectHotspots(entries, 1)
	if len(result) != 2 {
		t.Fatalf("expected 2 hotspots, got %d", len(result))
	}
	// beta has more flips, should rank first
	if result[0].Host != "beta" {
		t.Errorf("expected beta first, got %s", result[0].Host)
	}
}

func TestDetectHotspots_Empty(t *testing.T) {
	result := DetectHotspots(nil, 1)
	if len(result) != 0 {
		t.Errorf("expected empty result for nil input")
	}
}
