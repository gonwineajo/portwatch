// Package notifier provides pluggable notification channels for portwatch.
//
// Supported channels:
//   - stdout  — prints human-readable messages to standard output (default)
//   - webhook — HTTP POST JSON payload to a configured URL (e.g. Slack, Teams)
//   - script  — executes an external script with host and message arguments
//
// Usage:
//
//	n := notifier.New(notifier.Config{
//		Channel: notifier.ChannelWebhook,
//		Target:  "https://hooks.example.com/portwatch",
//	})
//	n.Notify("192.168.1.1", diff)
package notifier
