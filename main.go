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

	fmt.Println("\nTop 10 Most Active Repositories:")
	fmt.Printf("%-3s | %-12s | %-14s | %-7s | %-12s | %-11s | %-11s\n",
		"Pos", "Repository", "Activity Score", "Commits", "Contributors", "Active Days", "Avg Churn")
	fmt.Println("------------------------------------------------------------------------------------------")
	for i := 0; i < 10 && i < len(ranked); i++ {
		r := ranked[i]
		avgChurn := 0.0
		if r.CommitCount > 0 {
			avgChurn = r.TotalChurn / float64(r.CommitCount)
		}
		repoName := r.Repository
		if len(repoName) > 12 {
			repoName = repoName[:9] + "..."
		}
		fmt.Printf("%3d | %-12s | %14.4f | %7d | %12d | %11d | %11.2f\n",
			i+1, repoName, r.Score, r.CommitCount, r.UniqueContributors, r.ActiveDays, avgChurn)
	}

	return nil
}
