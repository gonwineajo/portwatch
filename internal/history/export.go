package history

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ExportCSV writes history entries as CSV to the given writer.
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
