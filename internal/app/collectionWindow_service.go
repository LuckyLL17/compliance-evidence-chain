package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func (s *Service) CreateCollectionWindow(value domain.CollectionWindow) (domain.CollectionWindow, error) {
	value.ID = domain.ID(platform.NewID("evidence-collectionWindow"))
	value.Prepare(s.clock.Now())
	if value.Name == "" {
		value.Name = "unnamed-window"
	}
	if value.Owner == "" {
		value.Owner = "system"
	}
	if err := value.Validate(); err != nil {
		return value, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, exists := s.store.collection_windows[value.ID]; exists {
		return value, ErrConflict
	}
	s.store.collection_windows[value.ID] = value
	s.recordLocked("create-collectionWindow", value.ID, value.Owner, value.Key())
	return value, nil
}

func (s *Service) GetCollectionWindow(id domain.ID) (domain.CollectionWindow, bool) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	value, ok := s.store.collection_windows[id]
	return value, ok
}

func (s *Service) ListCollectionWindows(term string, limit int) []domain.CollectionWindow {
	s.store.mu.RLock()
	values := make([]domain.CollectionWindow, 0, len(s.store.collection_windows))
	needle := strings.ToLower(strings.TrimSpace(term))
	for _, value := range s.store.collection_windows {
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

func (s *Service) AdvanceCollectionWindow(id domain.ID, next domain.Status, actor string) (domain.CollectionWindow, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	value, ok := s.store.collection_windows[id]
	if !ok {
		return value, ErrNotFound
	}
	if err := value.Advance(next, s.clock.Now()); err != nil {
		return value, fmt.Errorf("%w: %v", ErrInvalidCommand, err)
	}
	s.store.collection_windows[id] = value
	s.recordLocked("advance-collectionWindow", id, actor, value.Key())
	return value, nil
}

func (s *Service) SummarizeCollectionWindow(id domain.ID) (domain.EntityView, error) {
	value, ok := s.GetCollectionWindow(id)
	if !ok {
		return domain.EntityView{}, ErrNotFound
	}
	return value.View(), nil
}
