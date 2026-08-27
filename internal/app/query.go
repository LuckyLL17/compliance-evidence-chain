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
		if needle != "" && !matchEventSearch(event, needle) {
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
		return results[i].UpdatedAt.After(results[j].UpdatedAt)
	})
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}
	return results
}

// matchEventSearch reports whether a normalized search term matches an event's
// searchable text. The audit payload is part of the searchable surface so that
// compliance staff can locate an operation by its payload content, not just by
// the event kind.
func matchEventSearch(event domain.Event, needle string) bool {
	if strings.Contains(strings.ToLower(event.Kind), needle) {
		return true
	}
	return strings.Contains(strings.ToLower(event.Payload), needle)
}

func (s *Service) Describe() map[string]any {
	return map[string]any{
		"domain":       "control versions, evidence collection windows, cryptographic fingerprints, review decisions, and exceptions",
		"primary_kind": s.PrimaryKind(),
		"primary_type": s.PrimaryType(),
		"modules":      []string{"domain", "store", "service", "http", "jobs", "audit"},
	}
}
