package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "portwatch.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoad_Valid(t *testing.T) {
	path := writeTemp(t, `
hosts:
  - localhost
  - 192.168.1.1
ports: "22,80,443"
interval: 30
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Hosts) != 2 {
		t.Errorf("expected 2 hosts, got %d", len(cfg.Hosts))
	}
	if cfg.Ports != "22,80,443" {
		t.Errorf("unexpected ports: %s", cfg.Ports)
	}
	if cfg.Interval != 30 {
		t.Errorf("expected interval 30, got %d", cfg.Interval)
	}
}

func TestLoad_Defaults(t *testing.T) {
	path := writeTemp(t, `
hosts:
  - localhost
ports: "1-1024"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != 60 {
		t.Errorf("expected default interval 60, got %d", cfg.Interval)
	}
	if cfg.SnapshotDir != ".portwatch" {
		t.Errorf("expected default snapshot_dir, got %q", cfg.SnapshotDir)
	}
	if cfg.Alert.Type != "stdout" {
		t.Errorf("expected default alert type stdout, got %q", cfg.Alert.Type)
	}
}

func TestLoad_MissingHosts(t *testing.T) {
	path := writeTemp(t, `ports: "80"`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing hosts")
	}
}

func TestLoad_MissingPorts(t *testing.T) {
	path := writeTemp(t, `hosts:\n  - localhost`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing ports")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/portwatch.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
