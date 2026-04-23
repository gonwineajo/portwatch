package runner

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/user/portwatch/internal/history"
)

// PrintChurn reads the history file at path and prints a churn
// analysis table to stdout. If path is empty the default history
// file derived from the first configured host is used.
func PrintChurn(histPath string) error {
	return printChurnTo(os.Stdout, histPath)
}

func printChurnTo(w io.Writer, histPath string) error {
	if histPath == "" {
		return fmt.Errorf("history path must not be empty")
	}

	entries, err := history.Read(histPath)
	if err != nil {
		return fmt.Errorf("reading history: %w", err)
	}

	results := history.AnalyseChurn(entries)
	if len(results) == 0 {
		fmt.Fprintln(w, "no churn data available")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "HOST\tFLIPS\tUNIQUE PORTS\tSCORE")
	for _, r := range results {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%.2f\n",
			r.Host, r.TotalFlips, r.UniquePorts, r.Score)
	}
	return tw.Flush()
}
