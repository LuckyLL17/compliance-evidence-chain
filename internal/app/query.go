package app

import (
	"sort"
	"strings"

	"github.com/local/compliance-evidence-chain/internal/domain"
)

func (s *Service) Search(term string, limit int) []domain.EntityView {
	needle := strings.ToLower(strings.TrimSpace(term))
	results := make([]domain.EntityView, 0)
	s.store.mu.RLock()
	for _, event := range s.store.events {
		if needle != "" && !strings.Contains(strings.ToLower(event.Payload), needle) && !strings.Contains(strings.ToLower(event.Kind), needle) {
			continue
		}
		results = append(results, domain.EntityView{
			ID:        event.SubjectID,
			Kind:      event.Kind,
			Name:      event.Payload,
			Status:    domain.StatusActive,
			Owner:     event.Actor,
			Version:   1,
			UpdatedAt: event.CreatedAt,
		})
	}
	s.store.mu.RUnlock()
	sort.Slice(results, func(i, j int) bool {
		return results[i].UpdatedAt.After(results[j].UpdatedAt) || (results[i].UpdatedAt.Equal(results[j].UpdatedAt) && results[i].ID > results[j].ID)
	})
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}
	return results
}

func (s *Service) Describe() map[string]any {
	return map[string]any{
		"domain":       "control versions, evidence collection windows, cryptographic fingerprints, review decisions, and exceptions",
		"primary_kind": s.PrimaryKind(),
		"primary_type": s.PrimaryType(),
		"modules":      []string{"domain", "store", "service", "http", "jobs", "audit"},
	}
}
