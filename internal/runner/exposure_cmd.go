package runner

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/history"
)

// PrintExposure loads history from path and prints port exposure durations.
func PrintExposure(histPath string, limit int) error {
	return printExposureTo(os.Stdout, histPath, limit)
}

func printExposureTo(w io.Writer, histPath string, limit int) error {
	entries, err := history.Read(histPath)
	if err != nil {
		return fmt.Errorf("read history: %w", err)
	}

	results := history.AnalyseExposure(entries, time.Now())

	if len(results) == 0 {
		fmt.Fprintln(w, "no exposure data found")
		return nil
	}

	if limit > 0 && limit < len(results) {
		results = results[:limit]
	}

	fmt.Fprintf(w, "%-20s %-8s %-12s %s\n", "HOST", "PORT", "DURATION", "STATUS")
	fmt.Fprintf(w, "%-20s %-8s %-12s %s\n",
		"--------------------", "--------", "------------", "------")

	for _, r := range results {
		status := "closed"
		if r.StillOpen {
			status = "open"
		}
		fmt.Fprintf(w, "%-20s %-8d %-12s %s\n",
			r.Host, r.Port, formatDuration(r.Duration), status)
	}
	return nil
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.1fh", d.Hours())
	}
	return fmt.Sprintf("%.1fd", d.Hours()/24)
}
