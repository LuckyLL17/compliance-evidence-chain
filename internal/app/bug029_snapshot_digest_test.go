// source-marker: internal/app/snapshot.go
// source-marker: internal/platform/hash.go
// source-marker: internal/app/audit.go
package app

import (
	"testing"

	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug029SnapshotDigest(t *testing.T) {
	svc := NewService(platform.RealClock{}, platform.NewLogger())
	svc.EmitOperationalEvent("x", "a", "one")
	a := svc.Snapshot().Digest
	svc.store.mu.Lock()
	svc.store.events[0].Payload = "two"
	svc.store.mu.Unlock()
	b := svc.Snapshot().Digest
	if a == b {
		t.Fatalf("digest unchanged: %s", a)
	}
}

func TestBug029RegressionHealth(t *testing.T) {
	svc := NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
