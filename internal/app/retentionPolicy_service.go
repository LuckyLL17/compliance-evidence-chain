package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func (s *Service) CreateRetentionPolicy(value domain.RetentionPolicy) (domain.RetentionPolicy, error) {
	value.ID = domain.ID(platform.NewID("evidence-retentionPolicy"))
	value.Prepare(s.clock.Now())
	if value.Days < 0 {
		value.Days = 0
	}
	if err := value.Validate(); err != nil {
		return value, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, exists := s.store.retention_policies[value.ID]; exists {
		return value, ErrConflict
	}
	s.store.retention_policies[value.ID] = value
	s.recordLocked("create-retentionPolicy", value.ID, value.Owner, value.Key())
	return value, nil
}

func (s *Service) GetRetentionPolicy(id domain.ID) (domain.RetentionPolicy, bool) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	value, ok := s.store.retention_policies[id]
	return value, ok
}

func (s *Service) ListRetentionPolicys(term string, limit int) []domain.RetentionPolicy {
	s.store.mu.RLock()
	values := make([]domain.RetentionPolicy, 0, len(s.store.retention_policies))
	needle := strings.ToLower(strings.TrimSpace(term))
	for _, value := range s.store.retention_policies {
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

func (s *Service) AdvanceRetentionPolicy(id domain.ID, next domain.Status, actor string) (domain.RetentionPolicy, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	value, ok := s.store.retention_policies[id]
	if !ok {
		return value, ErrNotFound
	}
	if err := value.Advance(next, s.clock.Now()); err != nil {
		return value, fmt.Errorf("%w: %v", ErrInvalidCommand, err)
	}
	s.store.retention_policies[id] = value
	s.recordLocked("advance-retentionPolicy", id, actor, value.Key())
	return value, nil
}

func (s *Service) SummarizeRetentionPolicy(id domain.ID) (domain.EntityView, error) {
	value, ok := s.GetRetentionPolicy(id)
	if !ok {
		return domain.EntityView{}, ErrNotFound
	}
	return value.View(), nil
}
