// Package history provides tools for recording, querying, and analysing
// port-scan history for monitored hosts.
//
// Signature analysis
//
// BuildSignatures computes a deterministic fingerprint for each host based
// on the sorted list of open ports observed in its most recent scan event.
// Two hosts that share an identical fingerprint are considered to have the
// same port profile, which can surface misconfigured or cloned machines.
//
// Example usage:
//
//	entries, _ := history.Read("portwatch-history.json")
//	sigs := history.BuildSignatures(entries)
//	matches := history.MatchSignature(sigs, "192.168.1.10")
//	for _, m := range matches {
//		fmt.Printf("%s shares the same port profile\n", m.Host)
//	}
package history
