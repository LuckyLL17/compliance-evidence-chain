package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func (s *Service) CreateReviewDecision(value domain.ReviewDecision) (domain.ReviewDecision, error) {
	value.ID = domain.ID(platform.NewID("evidence-reviewDecision"))
	value.Prepare(s.clock.Now())
	if err := value.Validate(); err != nil {
		return value, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, exists := s.store.review_decisions[value.ID]; exists {
		return value, ErrConflict
	}
	s.store.review_decisions[value.ID] = value
	s.recordLocked("create-reviewDecision", value.ID, value.Owner, value.Key())
	return value, nil
}

func (s *Service) GetReviewDecision(id domain.ID) (domain.ReviewDecision, bool) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	value, ok := s.store.review_decisions[id]
	return value, ok
}

func (s *Service) ListReviewDecisions(term string, limit int) []domain.ReviewDecision {
	s.store.mu.RLock()
	values := make([]domain.ReviewDecision, 0, len(s.store.review_decisions))
	needle := strings.ToLower(strings.TrimSpace(term))
	for _, value := range s.store.review_decisions {
		if needle != "" && !strings.Contains(strings.ToLower(value.Name), needle) {
			continue
		}
		values = append(values, value)
	}
	s.store.mu.RUnlock()
	sort.Slice(values, func(i, j int) bool {
		return values[i].UpdatedAt.After(values[j].UpdatedAt)
	})
	if limit > 0 && len(values) > limit {
		values = values[:limit]
	}
	return values
}

func (s *Service) AdvanceReviewDecision(id domain.ID, next domain.Status, actor string) (domain.ReviewDecision, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	value, ok := s.store.review_decisions[id]
	if !ok {
		return value, ErrNotFound
	}
	if err := value.Advance(next, s.clock.Now()); err != nil {
		return value, fmt.Errorf("%w: %v", ErrInvalidCommand, err)
	}
	value.Owner = strings.TrimSpace(value.Owner)
	if actor == "" {
		actor = value.Owner
	}
	if actor == "" {
		actor = "system"
	}
	// BUG: updated map value is not persisted
	s.recordLocked("advance-reviewDecision", id, actor, value.Key())
	return value, nil
}

func (s *Service) SummarizeReviewDecision(id domain.ID) (domain.EntityView, error) {
	value, ok := s.GetReviewDecision(id)
	if !ok {
		return domain.EntityView{}, ErrNotFound
	}
	return value.View(), nil
}
