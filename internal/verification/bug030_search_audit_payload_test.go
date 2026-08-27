// source-marker: internal/app/query.go
// source-marker: internal/httpapi/search.go
// source-marker: internal/app/audit.go
package verification

import (
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug030SearchAuditPayload(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	svc.EmitOperationalEvent("op", "a", "needle-payload")
	got := svc.Search("needle-payload", 10)
	if len(got) != 1 {
		t.Fatalf("results=%#v", got)
	}
}

func TestBug030RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
