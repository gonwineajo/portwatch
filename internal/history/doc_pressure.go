// Package history provides tools for analysing port-scan history.
//
// # Port Pressure
//
// AnalysePressure measures how "unstable" a port is on each host by counting
// the number of times it transitions between open and closed ("flips") and
// normalising that count by the observation window in hours.
//
// A high pressure score indicates a port that opens and closes frequently,
// which may warrant investigation or alerting.
//
// Usage:
//
//	entries, _ := history.Read("portwatch.history.json")
//	result := history.AnalysePressure(entries, 2) // require at least 2 flips
//	for _, r := range result.Records {
//	    fmt.Printf("%s:%d flips=%d score=%.2f/hr\n",
//	        r.Host, r.Port, r.Flips, r.Score)
//	}
package history
