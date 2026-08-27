package app

import (
	"testing"

	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestSummarizeConsistencyScoreMatchesDetailTotal(t *testing.T) {
	s := NewService(platform.RealClock{}, platform.NewLogger())

	for _, seed := range []string{"audit-batch-2026-08", "evidence-consistency", "x"} {
		findings := s.ConsistencyFindings(seed)
		summary := s.SummarizeConsistency(seed)

		detailTotal := 0
		for _, f := range findings {
			detailTotal += f.Score
		}
		got, ok := summary["score"].(int)
		if !ok {
			t.Fatalf("score is not int: %T", summary["score"])
		}
		if got != detailTotal {
			t.Fatalf("seed=%q: summary score %d != detail total %d (off by %d)",
				seed, got, detailTotal, got-detailTotal)
		}
	}
}

func TestSummarizeConsistencyMatchesEvaluateConclusion(t *testing.T) {
	s := NewService(platform.RealClock{}, platform.NewLogger())
	seed := "audit-batch-2026-08"

	evaluated := s.EvaluateConsistency(seed)
	summary := s.SummarizeConsistency(seed)

	findings, ok := summary["findings"].([]ConsistencyFinding)
	if !ok {
		t.Fatalf("findings missing or wrong type: %T", summary["findings"])
	}

	// The evaluated finding must be the max-score finding in the same set,
	// so the audit path and the operational summary describe one batch consistently.
	maxScore := findings[0].Score
	for _, f := range findings {
		if f.Score > maxScore {
			maxScore = f.Score
		}
	}
	if evaluated.Score != maxScore {
		t.Fatalf("EvaluateConsistency score %d != max finding score %d", evaluated.Score, maxScore)
	}
}
