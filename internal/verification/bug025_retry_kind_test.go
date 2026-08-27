// source-marker: internal/jobs/retry.go
// source-marker: internal/app/service.go
// source-marker: internal/app/query.go
package verification

import (
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug025RetryKind(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	svc.EmitOperationalEvent("retry", "system", "queue")
	got := svc.Search("retry", 1)
	if len(got) != 1 {
		t.Fatalf("results=%#v", got)
	}
}

func TestBug025RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
