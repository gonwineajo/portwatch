package history

import (
	"testing"
	"time"
)

func seqEntry(host, event string, ports []int, offset time.Duration) Entry {
	return Entry{
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Add(offset),
		Host:      host,
		Event:     event,
		Ports:     ports,
	}
}

func TestBuildSequences_BasicOrder(t *testing.T) {
	entries := []Entry{
		seqEntry("host-a", EventOpened, []int{80}, 0),
		seqEntry("host-a", EventOpened, []int{443}, time.Minute),
		seqEntry("host-a", EventOpened, []int{8080}, 2*time.Minute),
	}
	seqs := BuildSequences(entries)
	if len(seqs) != 1 {
		t.Fatalf("expected 1 sequence, got %d", len(seqs))
	}
	want := []int{80, 443, 8080}
	for i, p := range want {
		if seqs[0].Ports[i] != p {
			t.Errorf("pos %d: want %d, got %d", i, p, seqs[0].Ports[i])
		}
	}
}

func TestBuildSequences_SkipsNonOpened(t *testing.T) {
	entries := []Entry{
		seqEntry("host-a", EventOpened, []int{80}, 0),
		seqEntry("host-a", EventClosed, []int{443}, time.Minute),
		seqEntry("host-a", EventNoChange, []int{8080}, 2*time.Minute),
	}
	seqs := BuildSequences(entries)
	if len(seqs[0].Ports) != 1 || seqs[0].Ports[0] != 80 {
		t.Errorf("expected only port 80, got %v", seqs[0].Ports)
	}
}

func TestBuildSequences_DeduplicatesPorts(t *testing.T) {
	entries := []Entry{
		seqEntry("host-a", EventOpened, []int{80, 443}, 0),
		seqEntry("host-a", EventOpened, []int{80, 8080}, time.Minute),
	}
	seqs := BuildSequences(entries)
	if len(seqs[0].Ports) != 3 {
		t.Errorf("expected 3 unique ports, got %d: %v", len(seqs[0].Ports), seqs[0].Ports)
	}
}

func TestBuildSequences_MultipleHosts(t *testing.T) {
	entries := []Entry{
		seqEntry("host-b", EventOpened, []int{22}, 0),
		seqEntry("host-a", EventOpened, []int{80}, 0),
	}
	seqs := BuildSequences(entries)
	if len(seqs) != 2 {
		t.Fatalf("expected 2 sequences, got %d", len(seqs))
	}
	// sorted by host name
	if seqs[0].Host != "host-a" || seqs[1].Host != "host-b" {
		t.Errorf("unexpected host order: %s, %s", seqs[0].Host, seqs[1].Host)
	}
}

func TestSequenceForHost_Found(t *testing.T) {
	entries := []Entry{
		seqEntry("host-a", EventOpened, []int{80}, 0),
	}
	s := SequenceForHost(entries, "host-a")
	if s == nil {
		t.Fatal("expected non-nil sequence")
	}
	if len(s.Ports) != 1 || s.Ports[0] != 80 {
		t.Errorf("unexpected ports: %v", s.Ports)
	}
}

func TestSequenceForHost_Missing(t *testing.T) {
	s := SequenceForHost([]Entry{}, "host-z")
	if s != nil {
		t.Errorf("expected nil for missing host, got %+v", s)
	}
}
