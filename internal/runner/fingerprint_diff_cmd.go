package runner

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/portwatch/internal/history"
)

// PrintFingerprintDiff loads history from path, builds signatures from the
// oldest and newest halves of entries, then prints any fingerprint changes.
func PrintFingerprintDiff(histPath string) error {
	return printFingerprintDiffTo(os.Stdout, histPath)
}

func printFingerprintDiffTo(w io.Writer, histPath string) error {
	entries, err := history.Read(histPath)
	if err != nil {
		return fmt.Errorf("fingerprint-diff: read history: %w", err)
	}
	if len(entries) == 0 {
		fmt.Fprintln(w, "no history entries found")
		return nil
	}

	mid := len(entries) / 2
	if mid == 0 {
		mid = 1
	}
	baseline := history.BuildSignatures(entries[:mid])
	current := history.BuildSignatures(entries[mid:])

	diffs := history.DiffSignatures(baseline, current)
	if len(diffs) == 0 {
		fmt.Fprintln(w, "no fingerprint changes detected")
		return nil
	}

	fmt.Fprintln(w, "FINGERPRINT CHANGES")
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, d := range diffs {
		fmt.Fprintf(w, "host: %s\n", d.Host)
		if len(d.Added) > 0 {
			fmt.Fprintf(w, "  added  : %s\n", formatPorts(d.Added))
		}
		if len(d.Removed) > 0 {
			fmt.Fprintf(w, "  removed: %s\n", formatPorts(d.Removed))
		}
	}
	return nil
}
