package history

import (
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Diff computes the change between two snapshots and returns an Entry.
func Diff(host string, prev, curr snapshot.Snapshot, ts time.Time) Entry {
	diff := snapshot.Compare(prev, curr)

	entry := Entry{
		Timestamp: ts,
		Host:      host,
		Opened:    diff.Opened,
		Closed:    diff.Closed,
	}

	if len(diff.Opened) > 0 {
		entry.Event = "opened"
	} else if len(diff.Closed) > 0 {
		entry.Event = "closed"
	} else {
		entry.Event = "unchanged"
	}

	return entry
}

// HasChanges returns true if the entry contains any opened or closed ports.
func HasChanges(e Entry) bool {
	return len(e.Opened) > 0 || len(e.Closed) > 0
}
