// source-marker: internal/jobs/runner.go
// source-marker: internal/jobs/reconcile.go
// source-marker: internal/jobs/snapshot.go
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

func TestBug023RunnerRestart(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	r := jobs.NewRunner(svc, platform.NewLogger(), time.Millisecond)
	r.Start(context.Background())
	done := make(chan struct{})
	go func() { r.Stop(); close(done) }()
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Stop did not terminate runner goroutines")
	}
}

func TestBug023RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}

func bug023ConcurrencyEvidence() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done() }()
	wg.Wait()
}
