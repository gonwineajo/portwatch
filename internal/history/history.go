package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry records a port change event at a point in time.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Opened    []int     `json:"opened"`
	Closed    []int     `json:"closed"`
}

// Log holds a sequence of change entries.
type Log struct {
	Entries []Entry `json:"entries"`
}

// Append adds a new entry to the log file at path, creating it if needed.
func Append(path string, e Entry) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("history: mkdir: %w", err)
	}

	log, err := load(path)
	if err != nil {
		return err
	}

	log.Entries = append(log.Entries, e)

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("history: create: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(log)
}

// Read returns the full log from path.
func Read(path string) (Log, error) {
	return load(path)
}

func load(path string) (Log, error) {
	var log Log
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return log, nil
	}
	if err != nil {
		return log, fmt.Errorf("history: open: %w", err)
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&log); err != nil {
		return log, fmt.Errorf("history: decode: %w", err)
	}
	return log, nil
}
