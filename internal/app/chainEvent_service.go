package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func (s *Service) CreateChainEvent(value domain.ChainEvent) (domain.ChainEvent, error) {
	value.ID = domain.ID(platform.NewID("evidence-chainEvent"))
	value.Prepare(s.clock.Now())
	if err := value.Validate(); err != nil {
		return value, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, exists := s.store.chain_events[value.ID]; exists {
		return value, ErrConflict
	}
	s.store.chain_events[value.ID] = value
	s.recordLocked("create-chainEvent", value.ID, value.Owner, value.Key())
	return value, nil
}

func (s *Service) GetChainEvent(id domain.ID) (domain.ChainEvent, bool) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	value, ok := s.store.chain_events[id]
	return value, ok
}

func (s *Service) ListChainEvents(term string, limit int) []domain.ChainEvent {
	s.store.mu.RLock()
	values := make([]domain.ChainEvent, 0, len(s.store.chain_events))
	needle := strings.ToLower(strings.TrimSpace(term))
	for _, value := range s.store.chain_events {
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

func (s *Service) AdvanceChainEvent(id domain.ID, next domain.Status, actor string) (domain.ChainEvent, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	value, ok := s.store.chain_events[id]
	if !ok {
		return value, ErrNotFound
	}
	if err := value.Advance(next, s.clock.Now()); err != nil {
		return value, fmt.Errorf("%w: %v", ErrInvalidCommand, err)
	}
	s.store.chain_events[id] = value
	s.recordLocked("advance-chainEvent", id, actor, value.Key())
	return value, nil
}

func (s *Service) SummarizeChainEvent(id domain.ID) (domain.EntityView, error) {
	value, ok := s.GetChainEvent(id)
	if !ok {
		return domain.EntityView{}, ErrNotFound
	}
	return value.View(), nil
}
