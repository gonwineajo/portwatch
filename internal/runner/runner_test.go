package runner

import (
	"net"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func freePort(t *testing.T) (int, func()) {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	return port, func() { l.Close() }
}

func testConfig(host string, ports []int) *config.Config {
	return &config.Config{
		Hosts:    []string{host},
		Ports:    ports,
		Interval: time.Second,
		Timeout:  500 * time.Millisecond,
	}
}

func TestRunOnce_CreatesSnapshot(t *testing.T) {
	port, close := freePort(t)
	defer close()

	dir := t.TempDir()
	old, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(old)

	cfg := testConfig("127.0.0.1", []int{port})
	r := New(cfg)

	if err := r.RunOnce(); err != nil {
		t.Fatalf("RunOnce error: %v", err)
	}

	path := snapshotFile("127.0.0.1")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected snapshot file %s to exist", path)
	}
}

func TestSnapshotFile_ColonReplaced(t *testing.T) {
	got := snapshotFile("127.0.0.1:8080")
	want := ".portwatch/127.0.0.1_8080.json"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRunOnce_NoHosts(t *testing.T) {
	cfg := &config.Config{
		Hosts:    []string{},
		Ports:    []int{80},
		Interval: time.Second,
		Timeout:  time.Second,
	}
	r := New(cfg)
	if err := r.RunOnce(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
