// source-marker: internal/app/accessRule_service.go
// source-marker: internal/httpapi/accessRule_handler.go
// source-marker: internal/app/audit.go
package verification

import (
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug017AccessActor(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	v, err := svc.CreateAccessRule(domain.AccessRule{Name: "a", Owner: "owner"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = svc.AdvanceAccessRule(v.ID, domain.StatusActive, "admin")
	if err != nil {
		t.Fatal(err)
	}
	a := svc.AuditTrail(1)
	if a[0].Actor != "admin" {
		t.Fatalf("audit=%#v", a)
	}
}

func TestBug017RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
