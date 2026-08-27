// source-marker: internal/app/replay_module.go
// source-marker: internal/httpapi/router.go
// source-marker: internal/app/query.go
package verification

import (
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug028ReplayOrder(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	fs := svc.ReplayFindings("seed")
	if len(fs) < 2 {
		t.Fatal("not enough findings")
	}
	got := svc.EvaluateReplay("seed")
	if got.CreatedAt.Before(fs[0].CreatedAt) {
		t.Fatalf("got=%#v first=%#v", got, fs[0])
	}
}

func TestBug028RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
