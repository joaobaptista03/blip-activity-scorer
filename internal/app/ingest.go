package app

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

// Commit represents a single record parsed from commits.csv.
type Commit struct {
	// Timestamp is the Unix timestamp of the commit.
	Timestamp int64
	// Username is the GitHub username of the commit author (empty if unknown).
	Username string
	// Repository is the name of the repository the commit was pushed to.
	Repository string
	// Files is the number of files changed by the commit.
	Files int
	// Additions is the number of line additions in this commit.
	Additions int
	// Deletions is the number of deletions in this commit.
	Deletions int
}

// IngestStats tracks parsing statistics during ingestion.
type IngestStats struct {
	// TotalRows is the total number of CSV data rows read (excluding header).
	TotalRows int64
	// ParsedRows is the number of rows successfully parsed into Commit structs.
	ParsedRows int64
	// SkippedRows is the number of malformed rows encountered and skipped.
	SkippedRows int64
}

// StreamCommits streams rows from r one-by-one, validates them, and invokes
// the callback function `handle` for each valid Commit.
func StreamCommits(r io.Reader, handle func(Commit) error) (IngestStats, error) {
	reader := csv.NewReader(r)
	reader.ReuseRecord = true

	// Read and validate header
	header, err := reader.Read()
	if err != nil {
		return IngestStats{}, fmt.Errorf("failed to read header: %w", err)
	}

	// Verify column names and count
	expectedHeader := []string{"timestamp", "username", "repository", "files", "additions", "deletions"}
	if len(header) != len(expectedHeader) {
		return IngestStats{}, fmt.Errorf("invalid header length: expected %d, got %d", len(expectedHeader), len(header))
	}
	for i, col := range expectedHeader {
		if header[i] != col {
			return IngestStats{}, fmt.Errorf("unexpected column at index %d: expected %q, got %q", i, col, header[i])
		}
	}

	var stats IngestStats
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		stats.TotalRows++
		if err != nil {
			stats.SkippedRows++
			continue
		}

		if len(record) != 6 {
			stats.SkippedRows++
			continue
		}

		commit, parseErr := parseRecord(record)
		if parseErr != nil {
			stats.SkippedRows++
			continue
		}

		stats.ParsedRows++
		if err := handle(commit); err != nil {
			return stats, fmt.Errorf("handler error at row %d: %w", stats.TotalRows, err)
		}
	}

	return stats, nil
}

func parseRecord(record []string) (Commit, error) {
	ts, err := strconv.ParseInt(record[0], 10, 64)
	if err != nil {
		return Commit{}, fmt.Errorf("invalid timestamp %q: %w", record[0], err)
	}

	files, err := strconv.Atoi(record[3])
	if err != nil {
		return Commit{}, fmt.Errorf("invalid files count %q: %w", record[3], err)
	}

	additions, err := strconv.Atoi(record[4])
	if err != nil {
		return Commit{}, fmt.Errorf("invalid additions %q: %w", record[4], err)
	}

	deletions, err := strconv.Atoi(record[5])
	if err != nil {
		return Commit{}, fmt.Errorf("invalid deletions %q: %w", record[5], err)
	}

	return Commit{
		Timestamp:  ts,
		Username:   record[1],
		Repository: record[2],
		Files:      files,
		Additions:  additions,
		Deletions:  deletions,
	}, nil
}
