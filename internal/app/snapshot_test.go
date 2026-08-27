package app

import (
	"testing"
	"time"

	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

// fakeClock gives deterministic, advancing timestamps so snapshots are reproducible.
type fakeClock struct{ t time.Time }

func (f *fakeClock) Now() time.Time                  { return f.t }
func (f *fakeClock) Since(v time.Time) time.Duration { return f.t.Sub(v) }

// TestSnapshotEventsMatchAudits reproduces the run-check failure
// snapshot=domain.Snapshot{Events:1, Audits:0}: after processing one batch of
// events, the snapshot's event and audit counts must stay in lockstep, since
// recordLocked appends one audit and one event per call.
func TestSnapshotEventsMatchAudits(t *testing.T) {
	clk := &fakeClock{t: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
	svc := NewService(clk, platform.NewLogger())

	// Process one batch of events through real service operations.
	frameworks := []domain.Framework{
		{Name: "ISO 27001", Owner: "ops", Jurisdiction: "EU"},
		{Name: "SOC 2", Owner: "ops", Jurisdiction: "US"},
	}
	for i, fw := range frameworks {
		clk.t = clk.t.Add(time.Minute)
		created, err := svc.CreateFramework(fw)
		if err != nil {
			t.Fatalf("CreateFramework[%d]: %v", i, err)
		}
		clk.t = clk.t.Add(time.Minute)
		if _, err := svc.AdvanceFramework(created.ID, "active", "ops"); err != nil {
			t.Fatalf("AdvanceFramework[%d]: %v", i, err)
		}
		// After each processed event the snapshot must report identical counts.
		snap := svc.Snapshot()
		if snap.Events != snap.Audits {
			t.Fatalf("snapshot counts diverged after event %d: snapshot=%+v", i, snap)
		}
	}

	// Final batch must still be consistent.
	snap := svc.Snapshot()
	if want := len(frameworks) * 2; snap.Events != want {
		t.Fatalf("snapshot.Events = %d, want %d", snap.Events, want)
	}
	if snap.Audits != snap.Events {
		t.Fatalf("snapshot counts diverged: snapshot=%+v", snap)
	}
}
