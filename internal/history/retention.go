package history

import (
	"time"
)

// RetentionPolicy defines how long history entries are kept.
type RetentionPolicy struct {
	MaxAge  time.Duration
	MaxRows int
}

// DefaultRetention keeps entries for 30 days, max 1000 rows.
var DefaultRetention = RetentionPolicy{
	MaxAge:  30 * 24 * time.Hour,
	MaxRows: 1000,
}

// Apply filters entries according to the retention policy.
// It trims by age first, then by row count (keeping most recent).
func (p RetentionPolicy) Apply(entries []Entry) []Entry {
	now := time.Now()
	filtered := entries[:0]
	for _, e := range entries {
		if p.MaxAge > 0 && now.Sub(e.Timestamp) > p.MaxAge {
			continue
		}
		filtered = append(filtered, e)
	}
	if p.MaxRows > 0 && len(filtered) > p.MaxRows {
		filtered = filtered[len(filtered)-p.MaxRows:]
	}
	return filtered
}

// Prune loads the history file, applies the retention policy, and rewrites it.
func Prune(path string, policy RetentionPolicy) error {
	entries, err := load(path)
	if err != nil {
		return err
	}
	pruned := policy.Apply(entries)
	if len(pruned) == len(entries) {
		return nil
	}
	return write(path, pruned)
}
