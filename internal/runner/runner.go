package runner

import (
	"fmt"
	"log"
	"path/filepath"
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
	alert   *alert.Alerter
	histDir string
}

// New creates a Runner from cfg, writing history under histDir.
func New(cfg *config.Config, histDir string) *Runner {
	return &Runner{
		cfg:     cfg,
		alert:   alert.New(cfg),
		histDir: histDir,
	}
}

// RunOnce performs a single scan cycle for all configured hosts.
func (r *Runner) RunOnce() error {
	for _, host := range r.cfg.Hosts {
		if err := r.runHost(host); err != nil {
			log.Printf("error scanning %s: %v", host, err)
		}
	}
	return nil
}

func (r *Runner) runHost(host string) error {
	sc := scanner.New(host, r.cfg.Ports, r.cfg.Timeout)
	result, err := sc.Scan()
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	snapFile := snapshotFile(r.cfg.SnapshotDir, host)
	old, _ := snapshot.Load(snapFile)
	diff := snapshot.Compare(old, result)

	if len(diff.Opened) > 0 || len(diff.Closed) > 0 {
		r.alert.Notify(host, diff)

		entry := history.Entry{
			Timestamp: time.Now().UTC(),
			Host:      host,
			Opened:    diff.Opened,
			Closed:    diff.Closed,
		}
		hPath := filepath.Join(r.histDir, snapshotFile("", host)+".history.json")
		if herr := history.Append(hPath, entry); herr != nil {
			log.Printf("history append: %v", herr)
		}
	}

	return snapshot.Save(snapFile, result)
}

func snapshotFile(dir, host string) string {
	safe := strings.ReplaceAll(host, ":", "_")
	if dir == "" {
		return safe
	}
	return filepath.Join(dir, safe+".json")
}
