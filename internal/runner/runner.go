package runner

import (
	"fmt"
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

// Runner orchestrates scanning, snapshot diffing, and alerting.
type Runner struct {
	cfg     *config.Config
	scanner *scanner.Scanner
	alerter *alert.Alerter
}

// New creates a Runner from the given config.
func New(cfg *config.Config) *Runner {
	return &Runner{
		cfg:     cfg,
		scanner: scanner.New(cfg.Timeout),
		alerter: alert.New(nil),
	}
}

// RunOnce performs a single scan cycle for all configured hosts.
func (r *Runner) RunOnce() error {
	for _, host := range r.cfg.Hosts {
		if err := r.scanHost(host); err != nil {
			log.Printf("error scanning host %s: %v", host, err)
		}
	}
	return nil
}

// Watch runs scans in a loop at the configured interval.
func (r *Runner) Watch() error {
	log.Printf("starting portwatch (interval: %s)", r.cfg.Interval)
	for {
		if err := r.RunOnce(); err != nil {
			log.Printf("scan cycle error: %v", err)
		}
		time.Sleep(r.cfg.Interval)
	}
}

func (r *Runner) scanHost(host string) error {
	result, err := r.scanner.Scan(host, r.cfg.Ports)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	snapshotPath := snapshotFile(host)
	prev, err := snapshot.Load(snapshotPath)
	if err != nil {
		// No previous snapshot; save current and continue.
		return snapshot.Save(snapshotPath, result)
	}

	diff := snapshot.Compare(prev, result)
	if err := r.alerter.Notify(host, diff); err != nil {
		log.Printf("alert error for %s: %v", host, err)
	}

	return snapshot.Save(snapshotPath, result)
}

func snapshotFile(host string) string {
	safe := ""
	for _, c := range host {
		if c == ':' || c == '/' {
			safe += "_"
		} else {
			safe += string(c)
		}
	}
	return fmt.Sprintf(".portwatch/%s.json", safe)
}
