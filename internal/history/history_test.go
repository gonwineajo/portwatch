package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAppendAndRead(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	e1 := Entry{
		Timestamp: time.Now().UTC().Truncate(time.Second),
		Host:      "localhost",
		Opened:    []int{80, 443},
		Closed:    []int{},
	}
	if err := Append(path, e1); err != nil {
		t.Fatalf("Append: %v", err)
	}

	e2 := Entry{
		Timestamp: time.Now().UTC().Truncate(time.Second),
		Host:      "localhost",
		Opened:    []int{},
		Closed:    []int{80},
	}
	if err := Append(path, e2); err != nil {
		t.Fatalf("Append second: %v", err)
	}

	log, err := Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(log.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(log.Entries))
	}
	if log.Entries[0].Host != "localhost" {
		t.Errorf("unexpected host: %s", log.Entries[0].Host)
	}
	if len(log.Entries[0].Opened) != 2 {
		t.Errorf("expected 2 opened ports in first entry")
	}
}

func TestRead_MissingFile(t *testing.T) {
	log, err := Read("/tmp/portwatch_nonexistent_history.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(log.Entries) != 0 {
		t.Errorf("expected empty log")
	}
}

func TestAppend_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "history.json")
	e := Entry{Timestamp: time.Now(), Host: "h", Opened: []int{22}, Closed: []int{}}
	if err := Append(path, e); err != nil {
		t.Fatalf("Append: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
