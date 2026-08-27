// source-marker: internal/app/audit.go
// source-marker: internal/app/workflow.go
// source-marker: internal/httpapi/audit.go
package verification

import (
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug006AuditActor(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	svc.EmitOperationalEvent("op", "alice", "payload")
	e := svc.EventStream(1)
	a := svc.AuditTrail(1)
	if e[0].Actor != "alice" || a[0].Actor != "alice" {
		t.Fatalf("event=%#v audit=%#v", e, a)
	}
}

func TestBug006RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
