package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event represents a single alert event for a port change.
type Event struct {
	Timestamp time.Time
	Host      string
	Level     Level
	Message   string
}

// Notifier writes alert events to an output destination.
type Notifier struct {
	out io.Writer
}

// New creates a Notifier that writes to the given writer.
// Pass nil to default to os.Stdout.
func New(out io.Writer) *Notifier {
	if out == nil {
		out = os.Stdout
	}
	return &Notifier{out: out}
}

// Notify formats and writes alert events derived from a snapshot diff.
func (n *Notifier) Notify(host string, diff snapshot.Diff) []Event {
	var events []Event

	for _, port := range diff.Opened {
		e := Event{
			Timestamp: time.Now(),
			Host:      host,
			Level:     LevelAlert,
			Message:   fmt.Sprintf("port %d newly OPEN", port),
		}
		events = append(events, e)
		fmt.Fprintf(n.out, "[%s] %s %s: %s\n", e.Timestamp.Format(time.RFC3339), e.Level, e.Host, e.Message)
	}

	for _, port := range diff.Closed {
		e := Event{
			Timestamp: time.Now(),
			Host:      host,
			Level:     LevelWarn,
			Message:   fmt.Sprintf("port %d newly CLOSED", port),
		}
		events = append(events, e)
		fmt.Fprintf(n.out, "[%s] %s %s: %s\n", e.Timestamp.Format(time.RFC3339), e.Level, e.Host, e.Message)
	}

	return events
}
