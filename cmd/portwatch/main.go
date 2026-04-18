package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/runner"
)

const defaultConfig = "portwatch.yaml"

func main() {
	configPath := flag.String("config", defaultConfig, "path to config file")
	once := flag.Bool("once", false, "run a single scan cycle and exit")
	flag.Parse()

	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "config file not found: %s\n", *configPath)
		os.Exit(1)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	r := runner.New(cfg)

	if *once {
		if err := r.RunOnce(); err != nil {
			log.Fatalf("scan error: %v", err)
		}
		return
	}

	if err := r.Watch(); err != nil {
		log.Fatalf("watch error: %v", err)
	}
}
