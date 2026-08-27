package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func (s *Service) CreateConnectorRun(value domain.ConnectorRun) (domain.ConnectorRun, error) {
	value.ID = domain.ID(platform.NewID("evidence-connectorRun"))
	value.Prepare(s.clock.Now())
	if err := value.Validate(); err != nil {
		return value, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, exists := s.store.connector_runs[value.ID]; exists {
		return value, ErrConflict
	}
	s.store.connector_runs[value.ID] = value
	s.recordLocked("create-connectorRun", value.ID, value.Owner, value.Key())
	return value, nil
}

func (s *Service) GetConnectorRun(id domain.ID) (domain.ConnectorRun, bool) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	value, ok := s.store.connector_runs[id]
	return value, ok
}

func (s *Service) ListConnectorRuns(term string, limit int) []domain.ConnectorRun {
	s.store.mu.RLock()
	values := make([]domain.ConnectorRun, 0, len(s.store.connector_runs))
	needle := strings.ToLower(strings.TrimSpace(term))
	for _, value := range s.store.connector_runs {
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

func (s *Service) AdvanceConnectorRun(id domain.ID, next domain.Status, actor string) (domain.ConnectorRun, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	value, ok := s.store.connector_runs[id]
	if !ok {
		return value, ErrNotFound
	}
	if err := value.Advance(next, s.clock.Now()); err != nil {
		return value, fmt.Errorf("%w: %v", ErrInvalidCommand, err)
	}
	s.store.connector_runs[id] = value
	s.recordLocked("advance-connectorRun", id, actor, value.Key())
	return value, nil
}

func (s *Service) SummarizeConnectorRun(id domain.ID) (domain.EntityView, error) {
	value, ok := s.GetConnectorRun(id)
	if !ok {
		return domain.EntityView{}, ErrNotFound
	}
	return value.View(), nil
}
