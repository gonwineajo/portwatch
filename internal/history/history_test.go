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

	e1 := Entry{Timestamp: time.Now(), Host: "localhost", Opened: []int{80}, Closed: []int{}}
	e2 := Entry{Timestamp: time.Now(), Host: "localhost", Opened: []int{}, Closed: []int{80}}

	if err := Append(path, e1); err != nil {
		t.Fatalf("append e1: %v", err)
	}
	if err := Append(path, e2); err != nil {
		t.Fatalf("append e2: %v", err)
	}

	entries, err := Read(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Host != "localhost" {
		t.Errorf("unexpected host: %s", entries[0].Host)
	}
	if len(entries[0].Opened) != 1 || entries[0].Opened[0] != 80 {
		t.Errorf("unexpected opened ports: %v", entries[0].Opened)
	}
}

func TestRead_MissingFile(t *testing.T) {
	entries, err := Read("/nonexistent/path/history.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(entries))
	}
}

func TestAppend_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "nested", "history.json")

	e := Entry{Timestamp: time.Now(), Host: "example.com", Opened: []int{443}, Closed: []int{}}
	if err := Append(path, e); err != nil {
		t.Fatalf("append: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
