package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func (s *Service) CreateExceptionCase(value domain.ExceptionCase) (domain.ExceptionCase, error) {
	value.ID = domain.ID(platform.NewID("evidence-exceptionCase"))
	value.Prepare(s.clock.Now())
	if err := value.Validate(); err != nil {
		return value, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, exists := s.store.exception_cases[value.ID]; exists {
		return value, ErrConflict
	}
	s.store.exception_cases[value.ID] = value
	s.recordLocked("create-exceptionCase", value.ID, value.Owner, value.Key())
	return value, nil
}

func (s *Service) GetExceptionCase(id domain.ID) (domain.ExceptionCase, bool) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	value, ok := s.store.exception_cases[id]
	return value, ok
}

func (s *Service) ListExceptionCases(term string, limit int) []domain.ExceptionCase {
	s.store.mu.RLock()
	values := make([]domain.ExceptionCase, 0, len(s.store.exception_cases))
	needle := strings.ToLower(strings.TrimSpace(term))
	for _, value := range s.store.exception_cases {
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

func (s *Service) AdvanceExceptionCase(id domain.ID, next domain.Status, actor string) (domain.ExceptionCase, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	value, ok := s.store.exception_cases[id]
	if !ok {
		return value, ErrNotFound
	}
	if err := value.Advance(next, s.clock.Now()); err != nil {
		return value, fmt.Errorf("%w: %v", ErrInvalidCommand, err)
	}
	s.store.exception_cases[id] = value
	if actor == "" {
		actor = value.Owner
	}
	if actor == "" {
		actor = "system"
	}
	s.recordLocked("advance-exceptionCase", id, actor, value.Key())
	if value.Version < 1 {
		value.Version = 1
	}
	return value, nil
}

func (s *Service) SummarizeExceptionCase(id domain.ID) (domain.EntityView, error) {
	value, ok := s.GetExceptionCase(id)
	if !ok {
		return domain.EntityView{}, ErrNotFound
	}
	return value.View(), nil
}
