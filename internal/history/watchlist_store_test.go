package history

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadWatchlist(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "watchlist.json")

	wl := NewWatchlist()
	wl.Add("host-a", 22, "SSH")
	wl.Add("host-b", 443, "HTTPS")

	if err := SaveWatchlist(path, wl); err != nil {
		t.Fatalf("SaveWatchlist: %v", err)
	}

	loaded, err := LoadWatchlist(path)
	if err != nil {
		t.Fatalf("LoadWatchlist: %v", err)
	}

	entries := loaded.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Host != "host-a" || entries[0].Port != 22 {
		t.Errorf("unexpected first entry: %+v", entries[0])
	}
	if entries[1].Label != "HTTPS" {
		t.Errorf("expected label HTTPS, got %s", entries[1].Label)
	}
}

func TestLoadWatchlist_MissingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent.json")

	wl, err := LoadWatchlist(path)
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(wl.Entries()) != 0 {
		t.Errorf("expected empty watchlist for missing file")
	}
}

func TestSaveWatchlist_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "dir", "watchlist.json")

	wl := NewWatchlist()
	wl.Add("host-c", 8080, "HTTP-alt")

	if err := SaveWatchlist(path, wl); err != nil {
		t.Fatalf("SaveWatchlist with nested dir: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
