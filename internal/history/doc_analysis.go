// Package history provides tools for recording, querying, and analysing
// port-scan history for portwatch.
//
// Analysis sub-system overview
//
// The analysis layer sits on top of the raw []Entry slice and exposes a
// collection of focused, composable functions. Each function is pure: it
// accepts entries (and optional parameters) and returns a result without
// mutating the input or performing I/O.
//
// Quick-reference
//
//	Function              | What it answers
//	----------------------|------------------------------------------------
//	DetectAnomalies       | Ports that appear far less often than average
//	AssessRisk            | Per-host risk score based on open/closed events
//	Forecast              | Predicted next-open time for recurring ports
//	DetectPatterns        | Ports that open on a recurring schedule
//	CorrelateOpenPorts    | Port pairs that frequently open together
//	Velocity              | Rate of port-change events over time
//	Trend                 | Open/closed counts bucketed into time windows
//	RollingWindow         | Aggregate stats within a sliding time window
//	Heatmap               | Activity counts by hour-of-day and day-of-week
//	ScoreHosts            | Composite activity score per host
//	DetectOutliers        | Hosts whose open-port count is a statistical outlier
//	ClusterByPorts        | Group hosts that share an identical port fingerprint
//	Pivot                 | Transpose entries into a host×port matrix
//	BuildChains           | Link opened→closed events into lifecycle chains
//	Lifecycle             | Duration a specific port stayed open per host
//	PortDeltas            | Net change in open-port count between snapshots
//	SnapshotDiffs         | Opened/closed sets between consecutive scans
//	BuildSignatures       | Stable fingerprint for each host's port profile
//	MatchSignature        | Find hosts whose fingerprint matches a query
//	ComputeSimilarity     | Jaccard similarity between two hosts' port sets
//	AllPairSimilarity     | Similarity matrix for every host pair
//	Baseline / Deviate    | Compare current scan against a stored baseline
//	Annotate              | Attach free-text notes to history entries
//	Compact               | Remove redundant no-change scan entries
//	Normalize             | Sort ports and deduplicate entries for clean storage
//
// All functions are safe to call concurrently on independent entry slices.
package history
