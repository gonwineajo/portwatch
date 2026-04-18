package scanner_test

import (
	"net"
	"testing"

	"github.com/yourorg/portwatch/internal/scanner"
)

// startListener opens a random TCP port and returns its port number and a closer.
func startListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return port, func() { ln.Close() }
}

func TestScan_DetectsOpenPort(t *testing.T) {
	port, close := startListener(t)
	defer close()

	s := scanner.New()
	result, err := s.Scan("127.0.0.1", []int{port})
	if err != nil {
		t.Fatalf("scan error: %v", err)
	}
	if len(result.Ports) != 1 || result.Ports[0] != port {
		t.Errorf("expected port %d open, got %v", port, result.Ports)
	}
}

func TestScan_ClosedPort(t *testing.T) {
	s := scanner.New()
	// Port 1 is almost certainly closed in test environments.
	result, err := s.Scan("127.0.0.1", []int{1})
	if err != nil {
		t.Fatalf("scan error: %v", err)
	}
	if len(result.Ports) != 0 {
		t.Errorf("expected no open ports, got %v", result.Ports)
	}
}

func TestScan_ResultHostSet(t *testing.T) {
	s := scanner.New()
	result, err := s.Scan("127.0.0.1", []int{})
	if err != nil {
		t.Fatalf("scan error: %v", err)
	}
	if result.Host != "127.0.0.1" {
		t.Errorf("expected host 127.0.0.1, got %s", result.Host)
	}
}
