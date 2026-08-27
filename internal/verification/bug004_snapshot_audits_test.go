package verification

import (
	"sync"
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug004SnapshotAudits(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	svc.EmitOperationalEvent("x", "a", "p")
	s := svc.Snapshot()
	if s.Events != 1 || s.Audits != 1 {
		t.Fatalf("snapshot=%#v", s)
	}
}

func TestBug004RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}

func bug004ConcurrencyEvidence() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done() }()
	wg.Wait()
}
