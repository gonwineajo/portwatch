package history

import (
	"time"

	"github.com/portwatch/internal/snapshot"
)

// Entry records a single scan diff event for a host.
type Entry struct {
	Timestamp time.Time     `json:"timestamp"`
	Host      string        `json:"host"`
	Diff      snapshot.Diff `json:"diff"`
}

// HasChanges returns true when the entry contains opened or closed ports.
func (e Entry) HasChanges() bool {
	return len(e.Diff.Opened) > 0 || len(e.Diff.Closed) > 0
}

// Summary returns a short human-readable description of the entry.
func (e Entry) Summary() string {
	if !e.HasChanges() {
		return e.Host + ": no changes"
	}
	msg := e.Host + ":"
	if len(e.Diff.Opened) > 0 {
		msg += " opened ports detected"
	}
	if len(e.Diff.Closed) > 0 {
		msg += " closed ports detected"
	}
	return msg
}
