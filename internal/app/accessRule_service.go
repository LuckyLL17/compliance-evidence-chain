package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func (s *Service) CreateAccessRule(value domain.AccessRule) (domain.AccessRule, error) {
	value.ID = domain.ID(platform.NewID("evidence-accessRule"))
	value.Prepare(s.clock.Now())
	if err := value.Validate(); err != nil {
		return value, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, exists := s.store.access_rules[value.ID]; exists {
		return value, ErrConflict
	}
	s.store.access_rules[value.ID] = value
	s.recordLocked("create-accessRule", value.ID, value.Owner, value.Key())
	return value, nil
}

func (s *Service) GetAccessRule(id domain.ID) (domain.AccessRule, bool) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	value, ok := s.store.access_rules[id]
	return value, ok
}

func (s *Service) ListAccessRules(term string, limit int) []domain.AccessRule {
	s.store.mu.RLock()
	values := make([]domain.AccessRule, 0, len(s.store.access_rules))
	needle := strings.ToLower(strings.TrimSpace(term))
	for _, value := range s.store.access_rules {
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

func (s *Service) AdvanceAccessRule(id domain.ID, next domain.Status, actor string) (domain.AccessRule, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	value, ok := s.store.access_rules[id]
	if !ok {
		return value, ErrNotFound
	}
	if err := value.Advance(next, s.clock.Now()); err != nil {
		return value, fmt.Errorf("%w: %v", ErrInvalidCommand, err)
	}
	s.store.access_rules[id] = value
	s.recordLocked("advance-accessRule", id, actor, value.Key())
	return value, nil
}

func (s *Service) SummarizeAccessRule(id domain.ID) (domain.EntityView, error) {
	value, ok := s.GetAccessRule(id)
	if !ok {
		return domain.EntityView{}, ErrNotFound
	}
	return value.View(), nil
}
