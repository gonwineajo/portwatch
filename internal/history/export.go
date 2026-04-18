package history

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"
)

// ExportCSV writes history entries to w in CSV format.
// Columns: timestamp, host, opened_ports, closed_ports
func ExportCSV(entries []Entry, w io.Writer) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	if err := cw.Write([]string{"timestamp", "host", "opened_ports", "closed_ports"}); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	for _, e := range entries {
		row := []string{
			e.Timestamp.UTC().Format(time.RFC3339),
			e.Host,
			joinInts(e.OpenedPorts),
			joinInts(e.ClosedPorts),
		}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}

	return cw.Error()
}

func joinInts(ports []int) string {
	if len(ports) == 0 {
		return ""
	}
	out := ""
	for i, p := range ports {
		if i > 0 {
			out += ";"
		}
		out += fmt.Sprintf("%d", p)
	}
	return out
}
