package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/config"
)

const version = "0.1.0"

func main() {
	// CLI flags
	configPath := flag.String("config", "portwatch.yaml", "Path to configuration file")
	interval := flag.Int("interval", 30, "Polling interval in seconds")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	showVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("portwatch v%s\n", version)
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		// Config file is optional; use defaults if not found
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
			os.Exit(1)
		}
		cfg = config.Default()
	}

	// Override config with CLI flags if explicitly set
	if *interval != 30 {
		cfg.Interval = *interval
	}
	if *verbose {
		cfg.Verbose = true
	}

	fmt.Printf("portwatch v%s starting — polling every %ds\n", version, cfg.Interval)

	// Set up the port monitor
	m := monitor.New(cfg)

	// Handle graceful shutdown on SIGINT / SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Printf("\nreceived %s, shutting down...\n", sig)
		m.Stop()
	}()

	// Start monitoring — blocks until Stop() is called
	if err := m.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "monitor error: %v\n", err)
		os.Exit(1)
	}
}
