package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level portwatch configuration.
type Config struct {
	Hosts     []string `yaml:"hosts"`
	Ports     string   `yaml:"ports"`
	Interval  int      `yaml:"interval"`  // seconds
	SnapshotDir string `yaml:"snapshot_dir"`
	Alert     AlertConfig `yaml:"alert"`
}

// AlertConfig holds alerting configuration.
type AlertConfig struct {
	Type    string `yaml:"type"`    // stdout, webhook
	Webhook string `yaml:"webhook"` // URL if type == webhook
}

// Load reads and parses a YAML config file at the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if len(c.Hosts) == 0 {
		return fmt.Errorf("config: at least one host is required")
	}
	if c.Ports == "" {
		return fmt.Errorf("config: ports must not be empty")
	}
	if c.Interval <= 0 {
		c.Interval = 60
	}
	if c.SnapshotDir == "" {
		c.SnapshotDir = ".portwatch"
	}
	if c.Alert.Type == "" {
		c.Alert.Type = "stdout"
	}
	return nil
}
