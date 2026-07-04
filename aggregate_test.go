package main

import (
	"math"
	"testing"
)

func TestMerge(t *testing.T) {
	commits1 := []Commit{
		{Timestamp: 1600000000, Username: "user1", Repository: "repoA", Files: 1, Additions: 10, Deletions: 5},
	}
	commits2 := []Commit{
		{Timestamp: 1600086400, Username: "user2", Repository: "repoA", Files: 2, Additions: 20, Deletions: 10},
		{Timestamp: 1600086400, Username: "", Repository: "repoA", Files: 1, Additions: 0, Deletions: 0},
	}
	commits3 := []Commit{
		{Timestamp: 1600172800, Username: "user1", Repository: "repoA", Files: 5, Additions: 50, Deletions: 25},
	}

	buildStats := func(commits []Commit) *RepoStats {
		s := NewRepoStats("repoA")
		for _, c := range commits {
			s.Update(c)
		}
		return s
	}

	statsA1 := buildStats(commits1)
	statsA2 := buildStats(commits2)
	statsA3 := buildStats(commits3)
	statsA1.Merge(statsA2)
	statsA1.Merge(statsA3)

	statsB1 := buildStats(commits1)
	statsB2 := buildStats(commits2)
	statsB3 := buildStats(commits3)
	statsB2.Merge(statsB3)
	statsB1.Merge(statsB2)

	statsC1 := buildStats(commits1)
	statsC2 := buildStats(commits2)
	statsC3 := buildStats(commits3)
	statsC3.Merge(statsC2)
	statsC3.Merge(statsC1)

	compareStats := func(s1, s2 *RepoStats) bool {
		if s1.Repository != s2.Repository {
			return false
		}
		if s1.CommitCount != s2.CommitCount {
			return false
		}
		if math.Abs(s1.TotalChurn-s2.TotalChurn) > 1e-9 {
			return false
		}
		if len(s1.UniqueContributors) != len(s2.UniqueContributors) {
			return false
		}
		for k := range s1.UniqueContributors {
			if _, ok := s2.UniqueContributors[k]; !ok {
				return false
			}
		}
		if len(s1.ActiveDays) != len(s2.ActiveDays) {
			return false
		}
		for k := range s1.ActiveDays {
			if _, ok := s2.ActiveDays[k]; !ok {
				return false
			}
		}
		return true
	}

	if !compareStats(statsA1, statsB1) {
		t.Error("Merge associativity check failed: (S1 + S2) + S3 != S1 + (S2 + S3)")
	}
	if !compareStats(statsA1, statsC3) {
		t.Error("Merge commutativity check failed: S1 + S2 + S3 != S3 + S2 + S1")
	}
}
