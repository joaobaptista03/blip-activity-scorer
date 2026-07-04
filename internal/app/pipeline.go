package app

import (
	"io"
	"runtime"
	"sync"
)

// RunPipeline runs the concurrent ingestion, deduplication, and aggregation pipeline.
// It streams records from r, deduplicates them using a main-thread Deduplicator,
// streams them in batches to worker goroutines, and merges their partial maps.
// numWorkers specifies the worker count (defaults to runtime.NumCPU() if <= 0).

// PipelineResult encapsulates the results of the concurrent pipeline.
type PipelineResult struct {
	Stats          map[string]*RepoStats
	IngestStats    IngestStats
	DuplicateCount int64
}

func RunPipeline(r io.Reader, numWorkers int) (PipelineResult, error) {
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}

	deduper := NewDeduplicator()
	var duplicateCount int64

	const batchSize = 1000
	commitChan := make(chan []Commit, numWorkers*2)
	resultChan := make(chan map[string]*RepoStats, numWorkers)

	var wg sync.WaitGroup

	// Spawn worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			localStats := make(map[string]*RepoStats)
			for batch := range commitChan {
				for _, c := range batch {
					s, exists := localStats[c.Repository]
					if !exists {
						s = NewRepoStats(c.Repository)
						localStats[c.Repository] = s
					}
					s.Update(c)
				}
			}
			resultChan <- localStats
		}()
	}

	// Stream and dispatch records in batches
	var currentBatch []Commit
	stats, err := StreamCommits(r, func(c Commit) error {
		if deduper.IsDuplicate(c) {
			duplicateCount++
			return nil
		}

		currentBatch = append(currentBatch, c)
		if len(currentBatch) >= batchSize {
			commitChan <- currentBatch
			currentBatch = make([]Commit, 0, batchSize)
		}
		return nil
	})

	if len(currentBatch) > 0 {
		commitChan <- currentBatch
	}
	close(commitChan)

	// Wait for workers in background and close results channel
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Merge partial results from all workers
	finalStats := make(map[string]*RepoStats)
	for workerStats := range resultChan {
		for repo, partial := range workerStats {
			s, exists := finalStats[repo]
			if !exists {
				finalStats[repo] = partial
			} else {
				s.Merge(partial)
			}
		}
	}

	return PipelineResult{
		Stats:          finalStats,
		IngestStats:    stats,
		DuplicateCount: duplicateCount,
	}, err
}
