package main

import (
	"strings"
	"testing"
)

func TestStreamCommitsMalformed(t *testing.T) {
	csvData := `timestamp,username,repository,files,additions,deletions
1610969774,user0,repo2,5,153,0
invalid_timestamp,user0,repo2,5,153,0
1610963057,user0,repo2,2,16,12
1614333792,user1,repo3,1,1
1614249997,,repo3,1,1,invalid_deletions
1614333792,user2,repo3,1,10,20
`
	reader := strings.NewReader(csvData)

	var parsed []Commit
	stats, err := StreamCommits(reader, func(c Commit) error {
		parsed = append(parsed, c)
		return nil
	})

	if err != nil {
		t.Fatalf("unexpected fatal error: %v", err)
	}

	if stats.TotalRows != 6 {
		t.Errorf("expected 6 data rows processed, got %d", stats.TotalRows)
	}

	if stats.ParsedRows != 3 {
		t.Errorf("expected 3 parsed rows, got %d", stats.ParsedRows)
	}

	if stats.SkippedRows != 3 {
		t.Errorf("expected 3 skipped rows, got %d", stats.SkippedRows)
	}

	if len(parsed) != 3 {
		t.Errorf("expected 3 parsed commits, got %d", len(parsed))
	}
}

// Aggregator accumulates repository stats for a stream of commits.
type Aggregator struct {
	Stats map[string]*RepoStats
}

// NewAggregator creates a new Aggregator.
func NewAggregator() *Aggregator {
	return &Aggregator{
		Stats: make(map[string]*RepoStats),
	}
}

// Add updates the aggregator stats for a given commit.
func (a *Aggregator) Add(c Commit) {
	s, exists := a.Stats[c.Repository]
	if !exists {
		s = NewRepoStats(c.Repository)
		a.Stats[c.Repository] = s
	}
	s.Update(c)
}
