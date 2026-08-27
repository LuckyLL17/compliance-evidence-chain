package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func (s *Service) CreateControl(value domain.Control) (domain.Control, error) {
	value.ID = domain.ID(platform.NewID("evidence-control"))
	value.Prepare(s.clock.Now())
	if err := value.Validate(); err != nil {
		return value, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, exists := s.store.controls[value.ID]; exists {
		return value, ErrConflict
	}
	s.store.controls[value.ID] = value
	s.recordLocked("create-control", value.ID, value.Owner, value.Key())
	return value, nil
}

func (s *Service) GetControl(id domain.ID) (domain.Control, bool) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	value, ok := s.store.controls[id]
	return value, ok
}

func (s *Service) ListControls(term string, limit int) []domain.Control {
	s.store.mu.RLock()
	values := make([]domain.Control, 0, len(s.store.controls))
	needle := strings.ToLower(strings.TrimSpace(term))
	for _, value := range s.store.controls {
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

func (s *Service) AdvanceControl(id domain.ID, next domain.Status, actor string) (domain.Control, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	value, ok := s.store.controls[id]
	if !ok {
		return value, ErrNotFound
	}
	if err := value.Advance(next, s.clock.Now()); err != nil {
		return value, fmt.Errorf("%w: %v", ErrInvalidCommand, err)
	}
	s.store.controls[id] = value
	s.recordLocked("advance-control", id, actor, value.Key())
	return value, nil
}

func (s *Service) SummarizeControl(id domain.ID) (domain.EntityView, error) {
	value, ok := s.GetControl(id)
	if !ok {
		return domain.EntityView{}, ErrNotFound
	}
	return value.View(), nil
}
