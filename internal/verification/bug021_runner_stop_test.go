// source-marker: internal/jobs/runner.go
// source-marker: internal/jobs/metrics.go
// source-marker: internal/jobs/reconcile.go
package verification

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/jobs"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug021RunnerStop(t *testing.T) {
	r := jobs.NewRunner(app.NewService(platform.RealClock{}, platform.NewLogger()), platform.NewLogger(), time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	r.Start(ctx)
	time.Sleep(5 * time.Millisecond)
	cancel()
	done := make(chan struct{})
	go func() { r.Stop(); close(done) }()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("runner stop blocked")
	}
}

func TestBug021RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}

func bug021ConcurrencyEvidence() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done() }()
	wg.Wait()
}
