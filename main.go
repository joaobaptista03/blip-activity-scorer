package main

import (
	"fmt"
	"os"

	"blip-activity-scorer/internal/app"
)

func main() {
	// Load configuration from config.yaml (falls back to defaults if missing)
	cfg, err := app.LoadConfig("config.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load configuration: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Loaded scoring weights: commits=%0.2f, contributors=%0.2f, churn=%0.2f, consistency=%0.2f\n",
		cfg.Weights.Commits, cfg.Weights.Contributors, cfg.Weights.Churn, cfg.Weights.Consistency)

	// Run concurrent aggregator
	commitsFile, err := os.Open("commits.csv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to open commits.csv: %v\n", err)
		os.Exit(1)
	}
	defer commitsFile.Close()

	result, err := app.RunPipeline(commitsFile, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: concurrent pipeline failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Ingestion completed. Total: %d, Parsed: %d, Skipped: %d, Duplicates: %d\n",
		result.IngestStats.TotalRows, result.IngestStats.ParsedRows, result.IngestStats.SkippedRows, result.DuplicateCount)
	fmt.Printf("Aggregated %d distinct repositories.\n", len(result.Stats))

	// Calculate scores and rank repositories
	ranked := app.CalculateScores(result.Stats, cfg.Weights)

	app.PrintTopTable(ranked)

	// Write full ranking to CSV
	outputFile := "ranking_full.csv"
	if err := app.WriteCSV(outputFile, ranked); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to write output CSV: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nFull ranking written to %s\n", outputFile)
}
