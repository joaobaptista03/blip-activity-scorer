package main

import (
	"math"
	"sort"
)

// RankedRepo represents a repository with its raw stats and normalized activity scores.
type RankedRepo struct {
	// Repository is the repository name.
	Repository string
	// Score is the final weighted composite activity score in the range [0, 1].
	Score float64
	// CommitCount is the raw number of commits.
	CommitCount int64
	// UniqueContributors is the number of distinct commit authors.
	UniqueContributors int
	// TotalChurn is the cumulative log-dampened code churn.
	TotalChurn float64
	// ActiveDays is the number of distinct UTC days with at least one commit.
	ActiveDays int

	// CommitScore is the normalized commit frequency component [0, 1].
	CommitScore float64
	// ContributorScore is the normalized contributor diversity component [0, 1].
	ContributorScore float64
	// ChurnScore is the normalized churn intensity component [0, 1].
	ChurnScore float64
	// ConsistencyScore is the normalized date consistency component [0, 1].
	ConsistencyScore float64
}

// AvgChurn returns the average log-dampened code churn per commit.
func (r RankedRepo) AvgChurn() float64 {
	if r.CommitCount == 0 {
		return 0.0
	}
	return r.TotalChurn / float64(r.CommitCount)
}

func CalculateScores(stats map[string]*RepoStats) []RankedRepo {
	if len(stats) == 0 {
		return nil
	}

	globalActiveDays := make(map[string]struct{})
	for _, s := range stats {
		for d := range s.ActiveDays {
			globalActiveDays[d] = struct{}{}
		}
	}
	totalDays := float64(len(globalActiveDays))
	if totalDays == 0 {
		totalDays = 1.0
	}

	type rawScores struct {
		commits      float64
		contributors float64
		churn        float64
		consistency  float64
	}

	rawMap := make(map[string]rawScores, len(stats))
	var maxCommits, maxContributors, maxChurn, maxConsistency float64

	for repo, s := range stats {
		rawCommits := math.Log(1.0 + float64(s.CommitCount))
		rawContributors := math.Log(1.0 + float64(len(s.UniqueContributors)))
		rawChurn := s.AvgChurn()
		rawConsistency := float64(len(s.ActiveDays)) / totalDays

		rawMap[repo] = rawScores{
			commits:      rawCommits,
			contributors: rawContributors,
			churn:        rawChurn,
			consistency:  rawConsistency,
		}

		if rawCommits > maxCommits {
			maxCommits = rawCommits
		}
		if rawContributors > maxContributors {
			maxContributors = rawContributors
		}
		if rawChurn > maxChurn {
			maxChurn = rawChurn
		}
		if rawConsistency > maxConsistency {
			maxConsistency = rawConsistency
		}
	}

	wCommits := 0.30
	wContributors := 0.20
	wChurn := 0.25
	wConsistency := 0.25

	ranked := make([]RankedRepo, 0, len(stats))
	for repo, s := range stats {
		raw := rawMap[repo]

		normCommits := 0.0
		if maxCommits > 0 {
			normCommits = raw.commits / maxCommits
		}

		normContributors := 0.0
		if maxContributors > 0 {
			normContributors = raw.contributors / maxContributors
		}

		normChurn := 0.0
		if maxChurn > 0 {
			normChurn = raw.churn / maxChurn
		}

		normConsistency := 0.0
		if maxConsistency > 0 {
			normConsistency = raw.consistency / maxConsistency
		}

		score := wCommits*normCommits +
			wContributors*normContributors +
			wChurn*normChurn +
			wConsistency*normConsistency

		ranked = append(ranked, RankedRepo{
			Repository:         repo,
			Score:              score,
			CommitCount:        s.CommitCount,
			UniqueContributors: len(s.UniqueContributors),
			TotalChurn:         s.TotalChurn,
			ActiveDays:         len(s.ActiveDays),
			CommitScore:        normCommits,
			ContributorScore:   normContributors,
			ChurnScore:         normChurn,
			ConsistencyScore:   normConsistency,
		})
	}

	sort.Slice(ranked, func(i, j int) bool {
		if math.Abs(ranked[i].Score-ranked[j].Score) < 1e-9 {
			return ranked[i].Repository < ranked[j].Repository
		}
		return ranked[i].Score > ranked[j].Score
	})

	return ranked
}
