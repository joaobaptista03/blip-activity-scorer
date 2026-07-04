package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func WriteCSV(filePath string, ranked []RankedRepo) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"rank",
		"repository",
		"score",
		"commits",
		"contributors",
		"active_days",
		"avg_churn_dampened",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	for i, r := range ranked {
		row := []string{
			strconv.Itoa(i + 1),
			r.Repository,
			strconv.FormatFloat(r.Score, 'f', 6, 64),
			strconv.FormatInt(r.CommitCount, 10),
			strconv.Itoa(r.UniqueContributors),
			strconv.Itoa(r.ActiveDays),
			strconv.FormatFloat(r.AvgChurn(), 'f', 4, 64),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row %d: %w", i+1, err)
		}
	}

	return nil
}

func PrintTopTable(ranked []RankedRepo) {
	fmt.Println("\nTop 10 Most Active Repositories:")
	fmt.Printf("%-3s | %-12s | %-14s | %-7s | %-12s | %-11s | %-11s\n",
		"Pos", "Repository", "Activity Score", "Commits", "Contributors", "Active Days", "Avg Churn")
	fmt.Println("------------------------------------------------------------------------------------------")
	for i := 0; i < 10 && i < len(ranked); i++ {
		r := ranked[i]
		repoName := r.Repository
		if len(repoName) > 12 {
			// Truncate repository name if it exceeds 12 characters to maintain table alignment
			repoName = repoName[:9] + "..."
		}
		fmt.Printf("%3d | %-12s | %14.4f | %7d | %12d | %11d | %11.2f\n",
			i+1, repoName, r.Score, r.CommitCount, r.UniqueContributors, r.ActiveDays, r.AvgChurn())
	}
}
