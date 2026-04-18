package scanner_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/scanner"
)

func TestParsePortRange_Single(t *testing.T) {
	ports, err := scanner.ParsePortRange("80")
	if err != nil || len(ports) != 1 || ports[0] != 80 {
		t.Errorf("unexpected result: %v %v", ports, err)
	}
}

func TestParsePortRange_List(t *testing.T) {
	ports, err := scanner.ParsePortRange("22,80,443")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []int{22, 80, 443}
	for i, p := range expected {
		if ports[i] != p {
			t.Errorf("pos %d: want %d got %d", i, p, ports[i])
		}
	}
}

func TestParsePortRange_Range(t *testing.T) {
	ports, err := scanner.ParsePortRange("8000-8003")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 4 {
		t.Errorf("expected 4 ports, got %d", len(ports))
	}
}

func TestParsePortRange_Mixed(t *testing.T) {
	ports, err := scanner.ParsePortRange("22,8000-8001,443")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 4 {
		t.Errorf("expected 4 ports, got %d", len(ports))
	}
}

func TestParsePortRange_InvalidRange(t *testing.T) {
	_, err := scanner.ParsePortRange("9000-8000")
	if err == nil {
		t.Error("expected error for inverted range")
	}
}

func TestParsePortRange_InvalidToken(t *testing.T) {
	_, err := scanner.ParsePortRange("abc")
	if err == nil {
		t.Error("expected error for non-numeric token")
	}
}
