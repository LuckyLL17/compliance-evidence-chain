// source-marker: internal/app/recovery_module.go
// source-marker: internal/httpapi/router.go
// source-marker: internal/app/query.go
package verification

import (
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug026RecoveryScore(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	findings := svc.RecoveryFindings("long-seed")
	best := findings[0]
	for _, f := range findings {
		if f.Score > best.Score {
			best = f
		}
	}
	if got := svc.EvaluateRecovery("long-seed"); got.Score != best.Score {
		t.Fatalf("got=%#v best=%#v", got, best)
	}
}

func TestBug026RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
