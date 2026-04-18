package alert

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/snapshot"
)

func TestNotify_OpenedPorts(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)

	diff := snapshot.Diff{Opened: []int{80, 443}, Closed: []int{}}
	events := n.Notify("localhost", diff)

	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	for _, e := range events {
		if e.Level != LevelAlert {
			t.Errorf("expected level ALERT, got %s", e.Level)
		}
		if e.Host != "localhost" {
			t.Errorf("expected host localhost, got %s", e.Host)
		}
	}
	out := buf.String()
	if !strings.Contains(out, "OPEN") {
		t.Errorf("expected output to contain OPEN, got: %s", out)
	}
}

func TestNotify_ClosedPorts(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)

	diff := snapshot.Diff{Opened: []int{}, Closed: []int{22}}
	events := n.Notify("192.168.1.1", diff)

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Level != LevelWarn {
		t.Errorf("expected level WARN, got %s", events[0].Level)
	}
	if !strings.Contains(buf.String(), "CLOSED") {
		t.Errorf("expected output to contain CLOSED")
	}
}

func TestNotify_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)

	diff := snapshot.Diff{Opened: []int{}, Closed: []int{}}
	events := n.Notify("localhost", diff)

	if len(events) != 0 {
		t.Errorf("expected 0 events, got %d", len(events))
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff")
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	n := New(nil)
	if n.out == nil {
		t.Error("expected non-nil writer")
	}
}
