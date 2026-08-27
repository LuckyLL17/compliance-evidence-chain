package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func (s *Service) CreateEvidenceObject(value domain.EvidenceObject) (domain.EvidenceObject, error) {
	value.ID = domain.ID(platform.NewID("evidence-evidenceObject"))
	value.Prepare(s.clock.Now())
	if err := value.Validate(); err != nil {
		return value, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, exists := s.store.evidence_objects[value.ID]; exists {
		return value, ErrConflict
	}
	s.store.evidence_objects[value.ID] = value
	s.recordLocked("create-evidenceObject", value.ID, value.Owner, value.Key())
	return value, nil
}

func (s *Service) GetEvidenceObject(id domain.ID) (domain.EvidenceObject, bool) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	value, ok := s.store.evidence_objects[id]
	return value, ok
}

func (s *Service) ListEvidenceObjects(term string, limit int) []domain.EvidenceObject {
	s.store.mu.RLock()
	values := make([]domain.EvidenceObject, 0, len(s.store.evidence_objects))
	needle := strings.ToLower(strings.TrimSpace(term))
	for _, value := range s.store.evidence_objects {
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

func (s *Service) AdvanceEvidenceObject(id domain.ID, next domain.Status, actor string) (domain.EvidenceObject, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	value, ok := s.store.evidence_objects[id]
	if !ok {
		return value, ErrNotFound
	}
	if err := value.Advance(next, s.clock.Now()); err != nil {
		return value, fmt.Errorf("%w: %v", ErrInvalidCommand, err)
	}
	s.store.evidence_objects[id] = value
	s.recordLocked("advance-evidenceObject", id, actor, value.Key())
	return value, nil
}

func (s *Service) SummarizeEvidenceObject(id domain.ID) (domain.EntityView, error) {
	value, ok := s.GetEvidenceObject(id)
	if !ok {
		return domain.EntityView{}, ErrNotFound
	}
	return value.View(), nil
}
