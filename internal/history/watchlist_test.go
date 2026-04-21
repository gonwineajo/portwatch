package history

import (
	"testing"
	"time"
)

func TestWatchlist_AddAndEntries(t *testing.T) {
	wl := NewWatchlist()
	wl.Add("host-a", 80, "HTTP")
	wl.Add("host-b", 443, "HTTPS")

	entries := wl.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Label != "HTTP" {
		t.Errorf("expected label HTTP, got %s", entries[0].Label)
	}
}

func TestWatchlist_Remove(t *testing.T) {
	wl := NewWatchlist()
	wl.Add("host-a", 80, "HTTP")
	wl.Add("host-a", 443, "HTTPS")
	wl.Remove("host-a", 80)

	entries := wl.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry after remove, got %d", len(entries))
	}
	if entries[0].Port != 443 {
		t.Errorf("expected port 443 to remain")
	}
}

func TestWatchlist_Remove_NoMatch(t *testing.T) {
	wl := NewWatchlist()
	wl.Add("host-a", 80, "HTTP")
	wl.Remove("host-b", 80) // different host, no-op

	if len(wl.Entries()) != 1 {
		t.Errorf("expected entry to remain when host does not match")
	}
}

func TestWatchlist_Match(t *testing.T) {
	wl := NewWatchlist()
	wl.Add("host-a", 22, "SSH")

	now := time.Now()
	historyEntries := []Entry{
		{Host: "host-a", Ports: []int{22, 80}, Event: "scan", Timestamp: now},
		{Host: "host-b", Ports: []int{443}, Event: "scan", Timestamp: now},
	}

	matched := wl.Match(historyEntries)
	if len(matched) != 1 {
		t.Fatalf("expected 1 matched entry, got %d", len(matched))
	}
	if matched[0].Host != "host-a" {
		t.Errorf("expected host-a, got %s", matched[0].Host)
	}
}

func TestWatchlist_Match_Empty(t *testing.T) {
	wl := NewWatchlist()
	now := time.Now()
	entries := []Entry{
		{Host: "host-a", Ports: []int{80}, Event: "scan", Timestamp: now},
	}
	matched := wl.Match(entries)
	if len(matched) != 0 {
		t.Errorf("expected no matches for empty watchlist")
	}
}
