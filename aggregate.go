package main

import (
	"math"
	"time"
)

// RepoStats tracks running statistics for a single repository.
type RepoStats struct {
	// Repository is the name of the repository.
	Repository string
	// CommitCount is the running count of commits.
	CommitCount int64
	// UniqueContributors tracks unique authors (including an "<unknown>" sentinel).
	UniqueContributors map[string]struct{}
	// TotalChurn accumulates the log-dampened churn values.
	TotalChurn float64
	// ActiveDays tracks unique dates of activity in UTC format ("2006-01-02").
	ActiveDays map[string]struct{}
}

// NewRepoStats initializes and returns a new RepoStats.
func NewRepoStats(repo string) *RepoStats {
	return &RepoStats{
		Repository:         repo,
		UniqueContributors: make(map[string]struct{}),
		ActiveDays:         make(map[string]struct{}),
	}
}

// Update updates the repository's stats with a commit record.
func (s *RepoStats) Update(c Commit) {
	s.CommitCount++

	// 1. Contributor diversity: treat blank username as a single "<unknown>" author
	author := c.Username
	if author == "" {
		author = "<unknown>"
	}
	s.UniqueContributors[author] = struct{}{}

	// 2. Log-dampened churn: ln(1 + additions + deletions)
	churn := float64(c.Additions + c.Deletions)
	s.TotalChurn += math.Log(1.0 + churn)

	// 3. Active days tracking: format timestamp as calendar day (UTC)
	day := time.Unix(c.Timestamp, 0).UTC().Format("2006-01-02")
	s.ActiveDays[day] = struct{}{}
}

// Merge merges another RepoStats into this one. This operation is associative and commutative.
func (s *RepoStats) Merge(other *RepoStats) {
	s.CommitCount += other.CommitCount
	s.TotalChurn += other.TotalChurn

	for contributor := range other.UniqueContributors {
		s.UniqueContributors[contributor] = struct{}{}
	}
	for day := range other.ActiveDays {
		s.ActiveDays[day] = struct{}{}
	}
}

// AvgChurn returns the average log-dampened code churn per commit.
func (s *RepoStats) AvgChurn() float64 {
	if s.CommitCount == 0 {
		return 0.0
	}
	return s.TotalChurn / float64(s.CommitCount)
}
