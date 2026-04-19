package history

import (
	"testing"
	"time"
)

var annotateTime = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func annotateEntries() []Entry {
	return []Entry{
		{Host: "host-a", Timestamp: annotateTime, Event: "opened", Ports: []int{80}},
		{Host: "host-b", Timestamp: annotateTime, Event: "closed", Ports: []int{443}},
		{Host: "host-a", Timestamp: annotateTime.Add(time.Hour), Event: "opened", Ports: []int{8080}},
	}
}

func TestAnnotate_MatchingEntry(t *testing.T) {
	entries := annotateEntries()
	result, count := Annotate(entries, "host-a", annotateTime, "suspicious port")
	if count != 1 {
		t.Fatalf("expected 1 annotation, got %d", count)
	}
	if result[0].Note != "suspicious port" {
		t.Errorf("expected note on first entry, got %q", result[0].Note)
	}
	if result[1].Note != "" {
		t.Errorf("expected no note on second entry")
	}
}

func TestAnnotate_NoMatch(t *testing.T) {
	entries := annotateEntries()
	_, count := Annotate(entries, "host-z", annotateTime, "note")
	if count != 0 {
		t.Errorf("expected 0 matches, got %d", count)
	}
}

func TestAnnotate_DoesNotMutateOriginal(t *testing.T) {
	entries := annotateEntries()
	_, _ = Annotate(entries, "host-a", annotateTime, "note")
	if entries[0].Note != "" {
		t.Errorf("original entries should not be mutated")
	}
}

func TestAnnotations_ReturnsOnlyNoted(t *testing.T) {
	entries := annotateEntries()
	result, _ := Annotate(entries, "host-b", annotateTime, "known maintenance")
	noted := Annotations(result)
	if len(noted) != 1 {
		t.Fatalf("expected 1 annotated entry, got %d", len(noted))
	}
	if noted[0].Host != "host-b" {
		t.Errorf("unexpected host %q", noted[0].Host)
	}
}

func TestAnnotations_Empty(t *testing.T) {
	noted := Annotations(annotateEntries())
	if len(noted) != 0 {
		t.Errorf("expected no annotations, got %d", len(noted))
	}
}
