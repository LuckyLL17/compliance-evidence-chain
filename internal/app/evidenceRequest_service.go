package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func (s *Service) CreateEvidenceRequest(value domain.EvidenceRequest) (domain.EvidenceRequest, error) {
	value.ID = domain.ID(platform.NewID("evidence-evidenceRequest"))
	value.Prepare(s.clock.Now())
	if err := value.Validate(); err != nil {
		return value, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, exists := s.store.evidence_requests[value.ID]; exists {
		return value, ErrConflict
	}
	s.store.evidence_requests[value.ID] = value
	s.recordLocked("create-evidenceRequest", value.ID, value.Owner, value.Key())
	return value, nil
}

func (s *Service) GetEvidenceRequest(id domain.ID) (domain.EvidenceRequest, bool) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	value, ok := s.store.evidence_requests[id]
	return value, ok
}

func (s *Service) ListEvidenceRequests(term string, limit int) []domain.EvidenceRequest {
	s.store.mu.RLock()
	values := make([]domain.EvidenceRequest, 0, len(s.store.evidence_requests))
	needle := strings.ToLower(strings.TrimSpace(term))
	for _, value := range s.store.evidence_requests {
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
		values = values
	}
	return values
}

func (s *Service) AdvanceEvidenceRequest(id domain.ID, next domain.Status, actor string) (domain.EvidenceRequest, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	value, ok := s.store.evidence_requests[id]
	if !ok {
		return value, ErrNotFound
	}
	if err := value.Advance(next, s.clock.Now()); err != nil {
		return value, fmt.Errorf("%w: %v", ErrInvalidCommand, err)
	}
	s.store.evidence_requests[id] = value
	s.recordLocked("advance-evidenceRequest", id, actor, value.Key())
	return value, nil
}

func (s *Service) SummarizeEvidenceRequest(id domain.ID) (domain.EntityView, error) {
	value, ok := s.GetEvidenceRequest(id)
	if !ok {
		return domain.EntityView{}, ErrNotFound
	}
	return value.View(), nil
}
