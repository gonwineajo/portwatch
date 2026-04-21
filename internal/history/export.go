package history

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ExportCSV writes history entries as CSV to the given writer.
// The CSV includes a header row with columns: timestamp, host, opened, closed.
// Port lists in the opened and closed columns are semicolon-separated.
func ExportCSV(entries []Entry, w io.Writer) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"timestamp", "host", "opened", "closed"}); err != nil {
		return fmt.Errorf("write header: %w", err)
	}
	for _, e := range entries {
		row := []string{
			e.Timestamp.UTC().Format("2006-01-02T15:04:05Z"),
			e.Host,
			joinInts(e.Opened),
			joinInts(e.Closed),
		}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}
	cw.Flush()
	return cw.Error()
}

// ExportCSVFiltered writes only entries matching the given host to the writer.
// If host is empty, all entries are written (equivalent to ExportCSV).
func ExportCSVFiltered(entries []Entry, host string, w io.Writer) error {
	if host == "" {
		return ExportCSV(entries, w)
	}
	filtered := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if e.Host == host {
			filtered = append(filtered, e)
		}
	}
	return ExportCSV(filtered, w)
}

// joinInts converts a slice of ints to a semicolon-separated string.
func joinInts(vals []int) string {
	if len(vals) == 0 {
		return ""
	}
	parts := make([]string, len(vals))
	for i, v := range vals {
		parts[i] = strconv.Itoa(v)
	}
	return strings.Join(parts, ";")
}
