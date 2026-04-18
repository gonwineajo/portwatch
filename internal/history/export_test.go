package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestExportCSV_Headers(t *testing.T) {
	var buf bytes.Buffer
	if err := ExportCSV(nil, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := strings.SplitN(buf.String(), "\n", 2)[0]
	if line != "timestamp,host,opened,closed" {
		t.Errorf("unexpected header: %q", line)
	}
}

func TestExportCSV_Rows(t *testing.T) {
	ts := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	entries := []Entry{
		{Timestamp: ts, Host: "localhost", Opened: []int{80, 443}, Closed: []int{22}},
		{Timestamp: ts, Host: "10.0.0.1", Opened: nil, Closed: nil},
	}
	var buf bytes.Buffer
	if err := ExportCSV(entries, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[1], "localhost") {
		t.Errorf("expected localhost in row: %s", lines[1])
	}
	if !strings.Contains(lines[1], "80;443") {
		t.Errorf("expected opened ports in row: %s", lines[1])
	}
	if !strings.Contains(lines[1], "22") {
		t.Errorf("expected closed port in row: %s", lines[1])
	}
	if !strings.HasSuffix(lines[2], ",10.0.0.1,,") {
		t.Errorf("expected empty ports for second row: %s", lines[2])
	}
}
