# Repository Activity Scorer

A high-performance, concurrent, memory-efficient Go command-line tool designed to ingest commit history from an inner-source development environment, clean data issues (e.g., duplicates), and score and rank repositories based on a multi-signal activity metric.

## Prerequisites

* **Go 1.22 or higher** is required to compile and run the project.
* **Make** (optional) for convenience targets.

## Project Structure

* `main.go` - Entrypoint that orchestrates the ingestion, processing, and output pipeline.
* `ingest.go` - High-performance streaming parser and defines the `Commit` and `IngestStats` models.
* `clean.go` - Streaming Deduplicator that filters out duplicate data entries.
* `aggregate.go` - Data structures and operations to collect stats per repository (defines `RepoStats`).
* `pipeline.go` - Concurrency orchestration layer that handles worker-pool chunking and fan-in stats merging.
* `score.go` - Multi-signal activity score scoring and ranking algorithm (defines `RankedRepo`).
* `output.go` - Output generation handler for exporting ranking results.
* `clean_test.go` - Deduplicator unit tests.
* `aggregate_test.go` - Merge and aggregate logic tests.
* `score_test.go` - Scoring logic tests.
* `ingest_test.go` - Ingestion parser tests.

## Running the Application

Ensure you have your dataset file named `commits.csv` inside the current working directory.

### Run Directly

To run the scorer immediately using Go:

```bash
go run .
```

### Build Binary

To build a compiled binary and run it:

```bash
go build -o blip-activity-scorer .
./blip-activity-scorer
```

## Running Tests

To run the unit test suite covering ingestion, deduplication, associative merging, and scoring logic:

```bash
go test -v ./...
```

## Inputs and Outputs

* **Input**: The program expects a CSV file named `commits.csv` in the root folder with columns `timestamp,username,repository,files,additions,deletions`.
* **Output (Console)**: The top 10 most active repositories printed as a formatted table.
* **Output (File)**: A full list of all ranked repositories saved to `ranking_full.csv` (excluded from Git).
