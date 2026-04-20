package history

import "time"

// Annotation attaches a human-readable note to a specific history entry.
type Annotation struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Note      string    `json:"note"`
}

// Annotate adds a note to entries matching the given host and timestamp.
// Returns the number of entries annotated.
func Annotate(entries []Entry, host string, ts time.Time, note string) ([]Entry, int) {
	count := 0
	result := make([]Entry, len(entries))
	copy(result, entries)
	for i, e := range result {
		if e.Host == host && e.Timestamp.Equal(ts) {
			result[i].Note = note
			count++
		}
	}
	return result, count
}

// Annotations returns all entries that have a non-empty note.
func Annotations(entries []Entry) []Entry {
	var out []Entry
	for _, e := range entries {
		if e.Note != "" {
			out = append(out, e)
		}
	}
	return out
}

// ClearAnnotation removes the note from entries matching the given host and
// timestamp. Returns the number of entries cleared.
func ClearAnnotation(entries []Entry, host string, ts time.Time) ([]Entry, int) {
	return Annotate(entries, host, ts, "")
}
