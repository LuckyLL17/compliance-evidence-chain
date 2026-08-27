package app

import (
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func (s *Service) Snapshot() domain.Snapshot {
	s.store.mu.RLock()
	counts := make(map[string]int64, len(s.store.counters))
	for key, value := range s.store.counters {
		counts[key] = value
	}
	events := append([]domain.Event(nil), s.store.events...)
	audits := append([]domain.AuditRecord(nil), s.store.audits...)
	chain := s.store.chain
	s.store.mu.RUnlock()

	// Order records deterministically so identical content yields an identical
	// digest regardless of insertion order. ID breaks CreatedAt ties.
	sort.SliceStable(events, func(i, j int) bool {
		if events[i].CreatedAt.Equal(events[j].CreatedAt) {
			return events[i].ID < events[j].ID
		}
		return events[i].CreatedAt.Before(events[j].CreatedAt)
	})
	sort.SliceStable(audits, func(i, j int) bool {
		if audits[i].CreatedAt.Equal(audits[j].CreatedAt) {
			return audits[i].ID < audits[j].ID
		}
		return audits[i].CreatedAt.Before(audits[j].CreatedAt)
	})

	// Fold the full content of every record into the digest, not just the
	// counts. The chain head anchors the summary to the tamper-evident chain;
	// re-hashing each record's actual fields guarantees that replacing an
	// event/audit's content while leaving the counts unchanged still changes
	// the digest, so the integrity check can detect the rewrite.
	root := chain
	for _, event := range events {
		root = platform.ChainHash(root, eventContent(event))
	}
	for _, audit := range audits {
		root = platform.ChainHash(root, auditContent(audit))
	}

	payload, _ := json.Marshal(struct {
		Counts  map[string]int64
		Events  int
		Audits  int
		Content string
	}{counts, len(events), len(audits), root})
	digest := platform.Hash(strings.TrimSpace(string(payload)))
	return domain.Snapshot{
		CreatedAt: s.clock.Now(),
		Counts:    counts,
		Digest:    digest,
		Events:    len(events),
		Audits:    len(audits),
	}
}

// eventContent renders an event's content into a stable string so it can be
// folded into the snapshot digest. Every field that can be tampered with is
// represented, so any rewrite is reflected in the digest.
func eventContent(event domain.Event) string {
	return strings.Join([]string{
		"event",
		string(event.ID),
		event.Kind,
		string(event.SubjectID),
		event.Actor,
		event.Payload,
		event.CreatedAt.Format(time.RFC3339Nano),
	}, "|")
}

// auditContent renders an audit record's content into a stable string so it
// can be folded into the snapshot digest. The record's Digest (its chain link)
// is included so the per-record tamper evidence is part of the summary.
func auditContent(audit domain.AuditRecord) string {
	return strings.Join([]string{
		"audit",
		string(audit.ID),
		audit.Action,
		string(audit.SubjectID),
		audit.Actor,
		audit.Digest,
		audit.CreatedAt.Format(time.RFC3339Nano),
	}, "|")
}
