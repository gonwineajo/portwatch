package snapshot

import (
	"encoding/json"
	"os"
	"time"
)

// PortSnapshot represents the state of open ports on a host at a point in time.
type PortSnapshot struct {
	Host      string    `json:"host"`
	Ports     []int     `json:"ports"`
	ScannedAt time.Time `json:"scanned_at"`
}

// Diff holds the changes between two snapshots.
type Diff struct {
	Host    string
	Opened  []int
	Closed  []int
}

// Save writes a snapshot to the given file path as JSON.
func Save(path string, snap PortSnapshot) error {
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Load reads a snapshot from the given file path.
func Load(path string) (PortSnapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return PortSnapshot{}, err
	}
	var snap PortSnapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return PortSnapshot{}, err
	}
	return snap, nil
}

// Compare returns the diff between a previous and current snapshot.
func Compare(prev, curr PortSnapshot) Diff {
	prevSet := toSet(prev.Ports)
	currSet := toSet(curr.Ports)

	diff := Diff{Host: curr.Host}

	for p := range currSet {
		if !prevSet[p] {
			diff.Opened = append(diff.Opened, p)
		}
	}
	for p := range prevSet {
		if !currSet[p] {
			diff.Closed = append(diff.Closed, p)
		}
	}
	return diff
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
