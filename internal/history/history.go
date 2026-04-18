package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single history record of port changes.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Opened    []int     `json:"opened"`
	Closed    []int     `json:"closed"`
}

// Append adds a new entry to the history file at path.
func Append(path string, entry Entry) error {
	entries, err := load(path)
	if err != nil {
		return err
	}
	entries = append(entries, entry)
	return write(path, entries)
}

// Read returns all history entries from the file at path.
func Read(path string) ([]Entry, error) {
	return load(path)
}

func load(path string) ([]Entry, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []Entry{}, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func write(path string, entries []Entry) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
