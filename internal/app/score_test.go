package app

import (
	"testing"
)

func TestScoring(t *testing.T) {
	stats := make(map[string]*RepoStats)

	r1 := NewRepoStats("repo1")
	r1.CommitCount = 10
	r1.UniqueContributors["user1"] = struct{}{}
	r1.UniqueContributors["user2"] = struct{}{}
	r1.UniqueContributors["user3"] = struct{}{}
	r1.TotalChurn = 25.0
	r1.ActiveDays["2021-01-01"] = struct{}{}
	r1.ActiveDays["2021-01-02"] = struct{}{}
	r1.ActiveDays["2021-01-03"] = struct{}{}
	stats["repo1"] = r1

	r2 := NewRepoStats("repo2")
	r2.CommitCount = 5
	r2.UniqueContributors["user1"] = struct{}{}
	r2.TotalChurn = 10.0
	r2.ActiveDays["2021-01-01"] = struct{}{}
	r2.ActiveDays["2021-01-02"] = struct{}{}
	stats["repo2"] = r2

	r3 := NewRepoStats("repo3")
	r3.CommitCount = 1
	r3.UniqueContributors["user2"] = struct{}{}
	r3.TotalChurn = 1.0
	r3.ActiveDays["2021-01-01"] = struct{}{}
	stats["repo3"] = r3

	ranked := CalculateScores(stats, DefaultConfig().Weights)

	if len(ranked) != 3 {
		t.Fatalf("expected 3 ranked repos, got %d", len(ranked))
	}

	if ranked[0].Repository != "repo1" {
		t.Errorf("expected rank 1 to be repo1, got %s", ranked[0].Repository)
	}

	if ranked[2].Repository != "repo3" {
		t.Errorf("expected rank 3 to be repo3, got %s", ranked[2].Repository)
	}

	for _, r := range ranked {
		if r.Score < 0.0 || r.Score > 1.0 {
			t.Errorf("repo %s has score out of bounds: %f", r.Repository, r.Score)
		}
	}
}
