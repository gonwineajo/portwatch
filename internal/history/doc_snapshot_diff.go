// Package history provides utilities for analysing port-scan history.
//
// # SnapshotDiff
//
// SnapshotDiffs computes the port-level difference between consecutive scan
// entries for each host found in a history slice.
//
// Each [SnapshotDiffEntry] records:
//   - Opened  – ports present in the later scan but not the earlier one.
//   - Closed  – ports present in the earlier scan but not the later one.
//   - Stable  – ports present in both scans.
//
// Only entries whose Event field equals "scan" are considered; all other event
// types ("opened", "closed", "no_change", …) are ignored so that the diff
// reflects actual scanner observations rather than derived events.
//
// Example:
//
//	diffs := history.SnapshotDiffs(entries)
//	for _, d := range diffs {
//		fmt.Printf("%s @ %s: +%v -%v\n", d.Host, d.At.Format(time.RFC3339), d.Opened, d.Closed)
//	}
package history
