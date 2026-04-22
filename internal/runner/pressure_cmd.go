package runner

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/user/portwatch/internal/history"
)

// PrintPressure loads history from path and prints a pressure report to stdout.
// minFlips controls the minimum number of open/close transitions required for
// a (host, port) pair to appear in the report.
func PrintPressure(histPath string, minFlips int) error {
	return printPressureTo(os.Stdout, histPath, minFlips)
}

func printPressureTo(w io.Writer, histPath string, minFlips int) error {
	entries, err := history.Read(histPath)
	if err != nil {
		return fmt.Errorf("pressure: read history: %w", err)
	}

	result := history.AnalysePressure(entries, minFlips)

	if len(result.Records) == 0 {
		fmt.Fprintln(w, "no port pressure detected")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "HOST\tPORT\tFLIPS\tSCORE (/hr)")
	for _, r := range result.Records {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%.4f\n",
			r.Host, r.Port, r.Flips, r.Score)
	}
	return tw.Flush()
}
