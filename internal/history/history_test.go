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

	e1 := Entry{Timestamp: time.Now().UTC(), Host: "localhost", Opened: []uint16{80, 443}}
	e2 := Entry{Timestamp: time.Now().UTC(), Host: "localhost", Closed: []uint16{80}}

	if err := Append(path, e1); err != nil {
		t.Fatalf("Append e1: %v", err)
	}
	if err := Append(path, e2); err != nil {
		t.Fatalf("Append e2: %v", err)
	}

	entries, err := Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Host != "localhost" {
		t.Errorf("unexpected host: %s", entries[0].Host)
	}
	if len(entries[1].Closed) != 1 || entries[1].Closed[0] != 80 {
		t.Errorf("unexpected closed ports: %v", entries[1].Closed)
	}
}

func TestRead_MissingFile(t *testing.T) {
	entries, err := Read("/nonexistent/path/history.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestAppend_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "dir", "history.json")

	e := Entry{Timestamp: time.Now().UTC(), Host: "10.0.0.1", Opened: []uint16{22}}
	if err := Append(path, e); err != nil {
		t.Fatalf("Append: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
