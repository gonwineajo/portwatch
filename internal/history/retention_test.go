package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/snapshot/snapshot"
)

func makeEntry(ago time.Duration, opened, closed int) Entry {
	ports := func(n int) []int {
		out := make([]int, n)
		for i := range out {
			out[i] = 8000 + i
		}
		return out
	}
	return Entry{
		Timestamp: time.Now().Add(-ago),
		Diff: snapshot.Diff{
			Opened: ports(opened),
			Closed: ports(closed),
		},
		Host: "localhost",
	}
}

func TestRetention_MaxAge(t *testing.T) {
	p := RetentionPolicy{MaxAge: 24 * time.Hour, MaxRows: 0}
	entries := []Entry{
		makeEntry(48*time.Hour, 1, 0),
		makeEntry(12*time.Hour, 1, 0),
		makeEntry(1*time.Hour, 1, 0),
	}
	got := p.Apply(entries)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestRetention_MaxRows(t *testing.T) {
	p := RetentionPolicy{MaxAge: 0, MaxRows: 2}
	entries := []Entry{
		makeEntry(3*time.Hour, 1, 0),
		makeEntry(2*time.Hour, 1, 0),
		makeEntry(1*time.Hour, 1, 0),
	}
	got := p.Apply(entries)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got[0].Timestamp.Before(entries[1].Timestamp.Add(-time.Second)) {
		t.Error("expected most recent entries to be kept")
	}
}

func TestPrune_WritesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	for i := 0; i < 5; i++ {
		ago := time.Duration(i+1) * 24 * time.Hour
		e := makeEntry(ago, 1, 0)
		if err := Append(path, e); err != nil {
			t.Fatal(err)
		}
	}

	policy := RetentionPolicy{MaxAge: 72 * time.Hour, MaxRows: 100}
	if err := Prune(path, policy); err != nil {
		t.Fatal(err)
	}

	entries, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries after prune, got %d", len(entries))
	}
	_ = os.Remove(path)
}
