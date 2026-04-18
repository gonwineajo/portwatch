package history

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"
	"time"
)

func TestExportCSV_Headers(t *testing.T) {
	var buf bytes.Buffer
	if err := ExportCSV(nil, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line (header only), got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "timestamp,host") {
		t.Errorf("unexpected header: %s", lines[0])
	}
}

func TestExportCSV_Rows(t *testing.T) {
	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	entries := []Entry{
		{Timestamp: now, Host: "192.168.1.1", OpenedPorts: []int{80, 443}, ClosedPorts: []int{22}},
		{Timestamp: now, Host: "10.0.0.1", OpenedPorts: nil, ClosedPorts: []int{8080}},
	}

	var buf bytes.Buffer
	if err := ExportCSV(entries, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := csv.NewReader(&buf)
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("csv parse error: %v", err)
	}
	if len(records) != 3 {
		t.Fatalf("expected 3 records (header+2), got %d", len(records))
	}

	row := records[1]
	if row[1] != "192.168.1.1" {
		t.Errorf("host mismatch: %s", row[1])
	}
	if row[2] != "80;443" {
		t.Errorf("opened_ports mismatch: %s", row[2])
	}
	if row[3] != "22" {
		t.Errorf("closed_ports mismatch: %s", row[3])
	}

	row2 := records[2]
	if row2[2] != "" {
		t.Errorf("expected empty opened_ports, got %s", row2[2])
	}
}
