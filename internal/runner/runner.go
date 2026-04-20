package runner

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/notifier"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

// Runner orchestrates a single scan cycle.
type Runner struct {
	cfg      *config.Config
	scanner  *scanner.Scanner
	alerter  *alert.Alerter
	notifier *notifier.Notifier
}

// New creates a Runner from the provided config.
func New(cfg *config.Config) *Runner {
	return &Runner{
		cfg:      cfg,
		scanner:  scanner.New(cfg.Timeout),
		alerter:  alert.New(nil),
		notifier: notifier.New(cfg),
	}
}

// RunOnce performs a scan for every configured host and processes diffs.
func (r *Runner) RunOnce() error {
	if len(r.cfg.Hosts) == 0 {
		return fmt.Errorf("no hosts configured")
	}
	ports, err := r.cfg.ParsedPorts()
	if err != nil {
		return fmt.Errorf("invalid port range: %w", err)
	}
	for _, host := range r.cfg.Hosts {
		if err := r.scanHost(host, ports); err != nil {
			log.Printf("error scanning %s: %v", host, err)
		}
	}
	return nil
}

func (r *Runner) scanHost(host string, ports []int) error {
	result := r.scanner.Scan(host, ports)
	current := snapshot.Snapshot{Host: host, Ports: result.Open, ScannedAt: time.Now()}

	file := snapshotFile(host, r.cfg.SnapshotDir)
	prev, _ := snapshot.Load(file)
	diff := snapshot.Compare(prev, current)

	if err := snapshot.Save(file, current); err != nil {
		return fmt.Errorf("save snapshot: %w", err)
	}

	r.alerter.Notify(diff)
	r.notifier.Notify(diff)

	entry := history.Entry{
		Host:      host,
		Timestamp: time.Now(),
		Ports:     result.Open,
		Event:     history.EventScan,
	}
	if history.HasChanges(diff) {
		if len(diff.Opened) > 0 {
			entry.Event = history.EventOpened
		} else {
			entry.Event = history.EventClosed
		}
	}
	if err := history.Append(historyFile(host, r.cfg.HistoryDir), entry); err != nil {
		return fmt.Errorf("append history: %w", err)
	}
	return nil
}

func snapshotFile(host, dir string) string {
	safe := strings.ReplaceAll(host, ":", "_")
	return fmt.Sprintf("%s/%s.json", dir, safe)
}

func historyFile(host, dir string) string {
	safe := strings.ReplaceAll(host, ":", "_")
	return fmt.Sprintf("%s/%s.json", dir, safe)
}
