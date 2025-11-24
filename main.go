package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"djdees/synthetic_steller_data/config"
	"djdees/synthetic_steller_data/generator"
	"djdees/synthetic_steller_data/output"
)

func main() {
	// Parse command-line flags
	cfg := config.ParseFlags()

	// Display configuration
	fmt.Println("=== Stellargen: Synthetic Stellar Data Generator ===")
	fmt.Printf("Number of Stars: %d\n", cfg.NumStars)
	fmt.Printf("Planets per Star: %d (max)\n", cfg.PlanetsPerStar)
	fmt.Printf("Exoplanets per Star: %d (max)\n", cfg.ExoPerStar)
	fmt.Printf("Output Format: %s\n", cfg.OutputFormat)
	fmt.Printf("Output Directory: %s\n", cfg.OutputDir)
	fmt.Printf("Seed: %d\n", cfg.Seed)
	fmt.Printf("Dry Run: %t\n", cfg.DryRun)
	fmt.Println()

	// Create generator configuration
	genCfg := generator.Config{
		NumStars:       cfg.NumStars,
		PlanetsPerStar: cfg.PlanetsPerStar,
		ExoPerStar:     cfg.ExoPerStar,
		Seed:           cfg.Seed,
	}

	// Generate data
	fmt.Println("Generating stellar data...")
	startTime := time.Now()
	data := generator.GenerateAll(genCfg)
	duration := time.Since(startTime)

	fmt.Printf("Generated %d stars, %d planets, %d exoplanets in %v\n",
		len(data.Stars), len(data.Planets), len(data.Exoplanets), duration)

	if cfg.DryRun {
		fmt.Println("\nDry run mode - no output written")
		return
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Write output based on format
	fmt.Printf("\nWriting output in %s format...\n", cfg.OutputFormat)
	startTime = time.Now()

	switch cfg.OutputFormat {
	case "csv":
		if err := output.WriteCSV(data, cfg.OutputDir); err != nil {
			log.Fatalf("Failed to write CSV: %v", err)
		}
	case "json":
		if err := output.WriteJSON(data, cfg.OutputDir); err != nil {
			log.Fatalf("Failed to write JSON: %v", err)
		}
	case "parquet":
		if err := output.WriteParquet(data, cfg.OutputDir); err != nil {
			log.Fatalf("Failed to write Parquet: %v", err)
		}
	case "cassandra":
		if cfg.ConfigFile == "" {
			log.Fatal("Cassandra output requires --config flag with YAML configuration file")
		}
		if err := output.WriteToCassandra(data, cfg.ConfigFile); err != nil {
			log.Fatalf("Failed to write to Cassandra: %v", err)
		}
	default:
		log.Fatalf("Unsupported output format: %s", cfg.OutputFormat)
	}

	duration = time.Since(startTime)
	fmt.Printf("Output written successfully in %v\n", duration)
	fmt.Println("\nDone!")
}
