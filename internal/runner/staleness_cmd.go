package runner

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/history"
)

// PrintStaleness reads the history file and prints staleness information
// for each host. Hosts whose last scan exceeds the threshold are flagged.
func PrintStaleness(histFile string, threshold time.Duration) error {
	return printStalenessTo(os.Stdout, histFile, threshold)
}

func printStalenessTo(w io.Writer, histFile string, threshold time.Duration) error {
	entries, err := history.Read(histFile)
	if err != nil {
		return fmt.Errorf("read history: %w", err)
	}

	results := history.AnalyseStaleness(entries, threshold, time.Now())
	if len(results) == 0 {
		fmt.Fprintln(w, "no scan history found")
		return nil
	}

	fmt.Fprintf(w, "%-20s  %-30s  %-14s  %s\n", "HOST", "LAST SEEN", "STALENESS", "STATUS")
	fmt.Fprintf(w, "%-20s  %-30s  %-14s  %s\n",
		"--------------------",
		"------------------------------",
	"--------------",
		"------")

	for _, r := range results {
		status := "ok"
		if r.IsStale {
			status = "STALE"
		}
		fmt.Fprintf(w, "%-20s  %-30s  %-14s  %s\n",
			r.Host,
			r.LastSeen.Format(time.RFC3339),
			formatDuration(r.Staleness),
			status,
		)
	}
	return nil
}
