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
	file, err := os.Open("commits.csv")
	if err != nil {
		return fmt.Errorf("failed to open commits.csv: %w", err)
	}
	defer file.Close()

	deduper := NewDeduplicator()
	var duplicateCount int64

	aggregator := NewAggregator()

	stats, err := StreamCommits(file, func(c Commit) error {
		if deduper.IsDuplicate(c) {
			duplicateCount++
			return nil
		}
		aggregator.Add(c)
		return nil
	})
	if err != nil {
		return err
	}

	fmt.Printf("Ingestion completed. Total: %d, Parsed: %d, Skipped: %d, Duplicates: %d\n",
		stats.TotalRows, stats.ParsedRows, stats.SkippedRows, duplicateCount)
	fmt.Printf("Aggregated %d distinct repositories.\n", len(aggregator.Stats))
	return nil
}
