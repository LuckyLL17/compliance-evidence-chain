// source-marker: internal/jobs/runner.go
// source-marker: internal/jobs/retry.go
// source-marker: internal/app/service.go
package verification

import (
	"context"
	"testing"
	"time"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/jobs"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug022RunnerContext(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	r := jobs.NewRunner(svc, platform.NewLogger(), time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r.Start(ctx)
	time.Sleep(8 * time.Millisecond)
	if got := len(svc.EventStream(0)); got != 0 {
		r.Stop()
		t.Fatalf("events=%d after canceled parent", got)
	}
	r.Stop()
}

func TestBug022RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
