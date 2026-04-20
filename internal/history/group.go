package history

import "sort"

// GroupResult holds entries grouped by a key.
type GroupResult struct {
	Key     string
	Entries []Entry
}

// GroupByHost groups history entries by their Host field.
// The returned slice is sorted alphabetically by host.
func GroupByHost(entries []Entry) []GroupResult {
	m := make(map[string][]Entry)
	for _, e := range entries {
		m[e.Host] = append(m[e.Host], e)
	}
	return sortedGroups(m)
}

// GroupByEvent groups history entries by their Event field (e.g. "opened", "closed").
// The returned slice is sorted alphabetically by event name.
func GroupByEvent(entries []Entry) []GroupResult {
	m := make(map[string][]Entry)
	for _, e := range entries {
		m[string(e.Event)] = append(m[string(e.Event)], e)
	}
	return sortedGroups(m)
}

// GroupByDate groups history entries by the calendar date (YYYY-MM-DD) of their
// Timestamp field using UTC. The returned slice is sorted chronologically.
func GroupByDate(entries []Entry) []GroupResult {
	m := make(map[string][]Entry)
	for _, e := range entries {
		day := e.Timestamp.UTC().Format("2006-01-02")
		m[day] = append(m[day], e)
	}
	return sortedGroups(m)
}

// sortedGroups converts a map of grouped entries into a sorted slice of GroupResult.
func sortedGroups(m map[string][]Entry) []GroupResult {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	result := make([]GroupResult, 0, len(keys))
	for _, k := range keys {
		result = append(result, GroupResult{Key: k, Entries: m[k]})
	}
	return result
}
