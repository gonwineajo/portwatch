package history

import (
	"testing"
	"time"
)

var sigBase = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func sigEntry(host string, event EventType, ports []int, offset int) Entry {
	return Entry{
		Host:      host,
		Event:     event,
		Ports:     ports,
		Timestamp: sigBase.Add(time.Duration(offset) * time.Hour),
	}
}

func TestBuildSignatures_BasicHosts(t *testing.T) {
	entries := []Entry{
		sigEntry("host-a", EventScan, []int{80, 443}, 0),
		sigEntry("host-b", EventScan, []int{80, 443}, 1),
	}
	sigs := BuildSignatures(entries)
	if len(sigs) != 2 {
		t.Fatalf("expected 2 signatures, got %d", len(sigs))
	}
	if sigs[0].Host != "host-a" {
		t.Errorf("expected host-a first, got %s", sigs[0].Host)
	}
}

func TestBuildSignatures_UsesLatestScan(t *testing.T) {
	entries := []Entry{
		sigEntry("host-a", EventScan, []int{80}, 0),
		sigEntry("host-a", EventScan, []int{80, 443}, 2),
	}
	sigs := BuildSignatures(entries)
	if len(sigs) != 1 {
		t.Fatalf("expected 1 signature, got %d", len(sigs))
	}
	if len(sigs[0].Ports) != 2 {
		t.Errorf("expected latest scan ports, got %v", sigs[0].Ports)
	}
}

func TestBuildSignatures_SkipsNonScan(t *testing.T) {
	entries := []Entry{
		sigEntry("host-a", EventOpened, []int{8080}, 0),
		sigEntry("host-b", EventScan, []int{22}, 0),
	}
	sigs := BuildSignatures(entries)
	if len(sigs) != 1 {
		t.Fatalf("expected 1 signature, got %d", len(sigs))
	}
	if sigs[0].Host != "host-b" {
		t.Errorf("expected host-b, got %s", sigs[0].Host)
	}
}

func TestBuildSignatures_Fingerprint(t *testing.T) {
	entries := []Entry{
		sigEntry("host-a", EventScan, []int{443, 80}, 0),
	}
	sigs := BuildSignatures(entries)
	if sigs[0].Fingerprint != "80,443" {
		t.Errorf("expected sorted fingerprint '80,443', got %s", sigs[0].Fingerprint)
	}
}

func TestMatchSignature_FindsMatch(t *testing.T) {
	entries := []Entry{
		sigEntry("host-a", EventScan, []int{80, 443}, 0),
		sigEntry("host-b", EventScan, []int{80, 443}, 0),
		sigEntry("host-c", EventScan, []int{22}, 0),
	}
	sigs := BuildSignatures(entries)
	matches := MatchSignature(sigs, "host-a")
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].Host != "host-b" {
		t.Errorf("expected host-b, got %s", matches[0].Host)
	}
}

func TestMatchSignature_NoMatch(t *testing.T) {
	entries := []Entry{
		sigEntry("host-a", EventScan, []int{80}, 0),
		sigEntry("host-b", EventScan, []int{443}, 0),
	}
	sigs := BuildSignatures(entries)
	matches := MatchSignature(sigs, "host-a")
	if len(matches) != 0 {
		t.Errorf("expected no matches, got %d", len(matches))
	}
}

func TestMatchSignature_UnknownHost(t *testing.T) {
	entries := []Entry{
		sigEntry("host-a", EventScan, []int{80}, 0),
	}
	sigs := BuildSignatures(entries)
	matches := MatchSignature(sigs, "ghost")
	if matches != nil {
		t.Errorf("expected nil for unknown host, got %v", matches)
	}
}
