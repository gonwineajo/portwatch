package runner

import (
	"fmt"
	"io"
	"os"

	"github.com/user/portwatch/internal/history"
)

// PrintDominance loads history from path and prints a dominance report to
// stdout. minHosts filters ports seen on fewer than that many distinct hosts.
func PrintDominance(histPath string, minHosts int) error {
	return printDominanceTo(os.Stdout, histPath, minHosts)
}

func printDominanceTo(w io.Writer, histPath string, minHosts int) error {
	entries, err := history.Read(histPath)
	if err != nil {
		return fmt.Errorf("dominance: read history: %w", err)
	}

	results := history.AnalyseDominance(entries, minHosts)
	if len(results) == 0 {
		fmt.Fprintln(w, "no dominant ports found")
		return nil
	}

	fmt.Fprintf(w, "%-8s  %-10s  %-10s  %s\n", "PORT", "HOSTS", "OPENS", "SCORE")
	fmt.Fprintf(w, "%-8s  %-10s  %-10s  %s\n", "----", "-----", "-----", "-----")
	for _, r := range results {
		fmt.Fprintf(w, "%-8d  %-10d  %-10d  %.1f\n",
			r.Port, r.HostCount, r.TotalOpen, r.Score)
	}
	return nil
}
