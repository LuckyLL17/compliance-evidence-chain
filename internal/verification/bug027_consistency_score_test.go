// source-marker: internal/app/consistency_module.go
// source-marker: internal/httpapi/router.go
// source-marker: internal/app/metrics.go
package verification

import (
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug027ConsistencyScore(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	fs := svc.ConsistencyFindings("seed")
	want := 0
	for _, f := range fs {
		want += f.Score
	}
	got := svc.SummarizeConsistency("seed")["score"].(int)
	if got != want {
		t.Fatalf("got=%d want=%d", got, want)
	}
}

func TestBug027RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
