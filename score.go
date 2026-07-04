package main

import (
	"math"
	"sort"
)

type RankedRepo struct {
	Repository         string
	Score              float64
	CommitCount        int64
	UniqueContributors int
	TotalChurn         float64
	ActiveDays         int

	CommitScore        float64
	ContributorScore   float64
	ChurnScore         float64
	ConsistencyScore   float64
}

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
