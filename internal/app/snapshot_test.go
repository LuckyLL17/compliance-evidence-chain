package app

import (
	"testing"

	"github.com/local/compliance-evidence-chain/internal/platform"
)

func newTestService() *Service {
	return NewService(platform.RealClock{}, platform.NewLogger())
}

// TestSnapshotDigestDetectsContentRewrite reproduces the bug where the audit
// snapshot digest only covered event/audit counts: with the same count but
// rewritten content the digest must change, otherwise the integrity check
// cannot detect the tampering.
func TestSnapshotDigestDetectsContentRewrite(t *testing.T) {
	s := newTestService()

	// Emit two operational events so the store has content to cover.
	s.EmitOperationalEvent("ingest-create", "alice", "first payload")
	s.EmitOperationalEvent("ingest-update", "bob", "second payload")

	first := s.Snapshot()

	// Rewrite the content of an existing event while keeping the count
	// unchanged. This is the exact tamper the bug report describes.
	s.store.mu.Lock()
	for i := range s.store.events {
		if s.store.events[i].Payload == "second payload" {
			s.store.events[i].Payload = "tampered payload"
			s.store.events[i].Actor = "mallory"
			break
		}
	}
	s.store.mu.Unlock()

	second := s.Snapshot()

	if first.Events != second.Events {
		t.Fatalf("count should be unchanged: first=%d second=%d", first.Events, second.Events)
	}
	if first.Audits != second.Audits {
		t.Fatalf("audit count should be unchanged: first=%d second=%d", first.Audits, second.Audits)
	}
	if first.Digest == second.Digest {
		t.Fatalf("snapshot digest did not change after event content rewrite: %s", first.Digest)
	}
}

// TestSnapshotDigestStableForIdenticalContent ensures the digest is
// deterministic for unchanged content, so legitimate recomputation does not
// produce spurious integrity findings.
func TestSnapshotDigestStableForIdenticalContent(t *testing.T) {
	s := newTestService()
	s.EmitOperationalEvent("ingest-create", "alice", "stable payload")
	s.EmitOperationalEvent("ingest-update", "bob", "stable payload two")

	first := s.Snapshot()
	second := s.Snapshot()

	if first.Digest != second.Digest {
		t.Fatalf("snapshot digest should be stable for identical content: first=%s second=%s", first.Digest, second.Digest)
	}
}

// TestSnapshotDigestDetectsAuditContentRewrite extends the coverage to audit
// records: rewriting an audit's content while leaving the count unchanged
// must also change the digest.
func TestSnapshotDigestDetectsAuditContentRewrite(t *testing.T) {
	s := newTestService()
	s.EmitOperationalEvent("ingest-create", "alice", "first payload")
	s.EmitOperationalEvent("ingest-update", "bob", "second payload")

	first := s.Snapshot()

	s.store.mu.Lock()
	for i := range s.store.audits {
		if s.store.audits[i].Actor == "bob" {
			s.store.audits[i].Actor = "mallory"
			s.store.audits[i].Action = "ingest-delete"
			break
		}
	}
	s.store.mu.Unlock()

	second := s.Snapshot()

	if first.Audits != second.Audits {
		t.Fatalf("audit count should be unchanged: first=%d second=%d", first.Audits, second.Audits)
	}
	if first.Digest == second.Digest {
		t.Fatalf("snapshot digest did not change after audit content rewrite: %s", first.Digest)
	}
}
