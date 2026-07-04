# Repository Activity Scorer

A high-performance, concurrent, memory-efficient Go command-line tool designed to ingest commit history from an inner-source development environment, clean data issues (e.g., duplicates), and score and rank repositories based on a multi-signal activity metric.

## 1. Algorithm & Scoring Design

To measure repository activity, the scorer aggregates statistics on each repository and evaluates a multi-signal scoring formula:

$$\text{Score} = w_1 \cdot \text{CommitScore} + w_2 \cdot \text{ContributorScore} + w_3 \cdot \text{ChurnScore} + w_4 \cdot \text{ConsistencyScore}$$

Where the sub-metrics and weights are:
* **$w_1$ (Commit Frequency) = 0.30**: Dampened using $\ln(1 + \text{CommitCount})$ to prevent minor commit spamming from dominating.
* **$w_2$ (Contributor Diversity) = 0.20**: Dampened using $\ln(1 + \text{UniqueContributors})$ to value team collaboration.
* **$w_3$ (Code Churn Intensity) = 0.25**: Average log-churn per commit, $\frac{\sum \ln(1 + \text{additions} + \text{deletions})}{\text{CommitCount}}$, dampening massive single commits.
* **$w_4$ (Consistency) = 0.25**: Active consistency over the time period, computed as $\frac{\text{ActiveDays}}{\text{TotalDays}}$ in UTC.

Each raw metric is normalized relative to the maximum observed value across all repositories to the range $[0, 1]$.

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

## Makefile

A `Makefile` is provided for convenience:

| Target | Description |
| ------ | ----------- |
| `make build` | Compile the binary |
| `make test` | Run unit tests with verbose output |
| `make run` | Build and run the scorer |
| `make bench` | Run Go benchmarks with memory stats |
| `make lint` | Run `go vet` and `gofmt -l` |
| `make clean` | Remove compiled binary and generated output |

The codebase is clean of static analysis warnings and complies fully with standard formatting rules. Run verification with `make lint`.

## Inputs and Outputs

* **Input**: The program expects a CSV file named `commits.csv` in the root folder with columns `timestamp,username,repository,files,additions,deletions`.
* **Output (Console)**: The top 10 most active repositories printed as a formatted table.
* **Output (File)**: A full list of all ranked repositories saved to `ranking_full.csv` (excluded from Git).
