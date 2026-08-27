package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func (s *Service) CreateFramework(value domain.Framework) (domain.Framework, error) {
	value.ID = domain.ID(platform.NewID("evidence-framework"))
	value.Prepare(s.clock.Now())
	if err := value.Validate(); err != nil {
		return value, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, exists := s.store.frameworks[value.ID]; exists {
		return value, ErrConflict
	}
	s.store.frameworks[value.ID] = value
	s.recordLocked("create-framework", value.ID, value.Owner, value.Key())
	return value, nil
}

func (s *Service) GetFramework(id domain.ID) (domain.Framework, bool) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	value, ok := s.store.frameworks[id]
	return value, ok
}

func (s *Service) ListFrameworks(term string, limit int) []domain.Framework {
	s.store.mu.RLock()
	values := make([]domain.Framework, 0, len(s.store.frameworks))
	needle := strings.ToLower(strings.TrimSpace(term))
	for _, value := range s.store.frameworks {
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

func (s *Service) AdvanceFramework(id domain.ID, next domain.Status, actor string) (domain.Framework, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	value, ok := s.store.frameworks[id]
	if !ok {
		return value, ErrNotFound
	}
	if err := value.Advance(next, s.clock.Now()); err != nil {
		return value, fmt.Errorf("%w: %v", ErrInvalidCommand, err)
	}
	s.store.frameworks[id] = value
	s.recordLocked("advance-framework", id, actor, value.Key())
	return value, nil
}

func (s *Service) SummarizeFramework(id domain.ID) (domain.EntityView, error) {
	value, ok := s.GetFramework(id)
	if !ok {
		return domain.EntityView{}, ErrNotFound
	}
	return value.View(), nil
}
