package runner

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/user/portwatch/internal/history"
)

// PrintSnapshotDiffs loads the history file for the given host and writes a
// human-readable diff table to w. If w is nil, os.Stdout is used.
func PrintSnapshotDiffs(histFile string, w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}

	entries, err := history.Read(histFile)
	if err != nil {
		return fmt.Errorf("read history: %w", err)
	}

	diffs := history.SnapshotDiffs(entries)
	if len(diffs) == 0 {
		fmt.Fprintln(w, "no snapshot diffs found")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "HOST\tTIMESTAMP\tOPENED\tCLOSED\tSTABLE")
	for _, d := range diffs {
		fmt.Fprintf(tw, "%s\t%s\t%v\t%v\t%v\n",
			d.Host,
			d.At.Format("2006-01-02 15:04:05"),
			formatPorts(d.Opened),
			formatPorts(d.Closed),
			formatPorts(d.Stable),
		)
	}
	return tw.Flush()
}

func formatPorts(ports []int) string {
	if len(ports) == 0 {
		return "-"
	}
	out := ""
	for i, p := range ports {
		if i > 0 {
			out += ","
		}
		out += fmt.Sprintf("%d", p)
	}
	return out
}
