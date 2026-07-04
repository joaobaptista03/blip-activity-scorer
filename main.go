package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Repository Activity Scorer started")
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Run concurrent aggregator
	commitsFile, err := os.Open("commits.csv")
	if err != nil {
		return fmt.Errorf("failed to open commits.csv: %w", err)
	}
	defer commitsFile.Close()

	result, err := RunPipeline(commitsFile, 0)
	if err != nil {
		return fmt.Errorf("concurrent pipeline failed: %w", err)
	}

	fmt.Printf("Ingestion completed. Total: %d, Parsed: %d, Skipped: %d, Duplicates: %d\n",
		result.IngestStats.TotalRows, result.IngestStats.ParsedRows, result.IngestStats.SkippedRows, result.DuplicateCount)
	fmt.Printf("Aggregated %d distinct repositories.\n", len(result.Stats))

	// Calculate scores and rank repositories
	ranked := CalculateScores(result.Stats)

	PrintTopTable(ranked)

	// Write full ranking to CSV
	outputFile := "ranking_full.csv"
	if err := WriteCSV(outputFile, ranked); err != nil {
		return fmt.Errorf("failed to write output CSV: %w", err)
	}
	fmt.Printf("\nFull ranking written to %s\n", outputFile)

	return nil
}
