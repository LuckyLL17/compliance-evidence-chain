// source-marker: internal/app/service.go
// source-marker: internal/app/metrics.go
// source-marker: internal/jobs/metrics.go
package verification

import (
	"sync"
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug005MetricsCounter(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	svc.TouchCounter("reconcile_runs")
	m := svc.Metrics()
	if m.Counters["reconcile_runs"] != 1 {
		t.Fatalf("metrics=%#v", m)
	}
}

func TestBug005RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}

func bug005ConcurrencyEvidence() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done() }()
	wg.Wait()
}
