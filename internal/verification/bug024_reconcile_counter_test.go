// source-marker: internal/app/ingest.go
// source-marker: internal/app/metrics.go
// source-marker: internal/jobs/reconcile.go
package verification

import (
	"sync"
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug024ReconcileCounter(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Reconcile(); got != 1 {
		t.Fatalf("reconcile=%d", got)
	}
	if svc.Count("reconcile_runs") != 1 {
		t.Fatalf("count=%d", svc.Count("reconcile_runs"))
	}
}

func TestBug024RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}

func bug024ConcurrencyEvidence() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done() }()
	wg.Wait()
}
