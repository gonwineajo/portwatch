package runner

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

// Runner orchestrates scanning, diffing, alerting, and history recording.
type Runner struct {
	cfg     *config.Config
	alerter *alert.Alerter
}

// New creates a new Runner from the given config.
func New(cfg *config.Config) *Runner {
	return &Runner{
		cfg:     cfg,
		alerter: alert.New(cfg.AlertOutput),
	}
}

// RunOnce performs a single scan cycle for all configured hosts.
func (r *Runner) RunOnce() error {
	if len(r.cfg.Hosts) == 0 {
		return fmt.Errorf("no hosts configured")
	}
	for _, host := range r.cfg.Hosts {
		if err := r.scanHost(host); err != nil {
			log.Printf("error scanning %s: %v", host, err)
		}
	}
	return nil
}

func (r *Runner) scanHost(host string) error {
	s := scanner.New(host, r.cfg.Ports, r.cfg.Timeout)
	current, err := s.Scan()
	if err != nil {
		return err
	}

	file := snapshotFile(r.cfg.StateDir, host)
	prev, _ := snapshot.Load(file)
	diff := snapshot.Compare(prev, current)

	if err := snapshot.Save(file, current); err != nil {
		return fmt.Errorf("save snapshot: %w", err)
	}

	if len(diff.Opened) > 0 || len(diff.Closed) > 0 {
		r.alerter.Notify(host, diff)
		entry := history.Entry{
			Timestamp: time.Now(),
			Host:      host,
			Opened:    diff.Opened,
			Closed:    diff.Closed,
		}
		if err := history.Append(historyFile(r.cfg.StateDir, host), entry); err != nil {
			log.Printf("history append error for %s: %v", host, err)
		}
	}
	return nil
}

func snapshotFile(dir, host string) string {
	return fmt.Sprintf("%s/%s.json", dir, strings.ReplaceAll(host, ":", "_"))
}

func historyFile(dir, host string) string {
	return fmt.Sprintf("%s/%s.history.json", dir, strings.ReplaceAll(host, ":", "_"))
}
