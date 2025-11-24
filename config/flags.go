package config

import (
	"flag"
	"time"
)

// AppConfig holds the application configuration from command-line flags
type AppConfig struct {
	NumStars       int
	PlanetsPerStar int
	ExoPerStar     int
	OutputFormat   string
	OutputDir      string
	ConfigFile     string
	Seed           int64
	DryRun         bool
}

// ParseFlags parses command-line flags and returns an AppConfig
func ParseFlags() *AppConfig {
	cfg := &AppConfig{}

	flag.IntVar(&cfg.NumStars, "num-stars", 100, "Number of stars to generate")
	flag.IntVar(&cfg.PlanetsPerStar, "planets-per-star", 8, "Maximum planets per star")
	flag.IntVar(&cfg.ExoPerStar, "exo-per-star", 5, "Maximum exoplanets per star")
	flag.StringVar(&cfg.OutputFormat, "output-format", "csv", "Output format: csv, json, parquet, cassandra")
	flag.StringVar(&cfg.OutputDir, "output-dir", "output", "Output directory")
	flag.StringVar(&cfg.ConfigFile, "config", "", "YAML config file for Cassandra")
	flag.Int64Var(&cfg.Seed, "seed", 0, "Random seed (0 for time-based)")
	flag.BoolVar(&cfg.DryRun, "dry-run", false, "Dry run mode (no output)")

	flag.Parse()

	// Use time-based seed if seed is 0
	if cfg.Seed == 0 {
		cfg.Seed = time.Now().UnixNano()
	}

	// Validate configuration
	if cfg.NumStars < 1 {
		cfg.NumStars = 1
	}
	if cfg.PlanetsPerStar > 15 {
		cfg.PlanetsPerStar = 15
	}
	if cfg.ExoPerStar > 8 {
		cfg.ExoPerStar = 8
	}

	return cfg
}
