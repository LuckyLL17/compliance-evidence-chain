// source-marker: internal/app/control_service.go
// source-marker: internal/httpapi/control_handler.go
// source-marker: internal/app/audit.go
package verification

import (
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug008ControlAuditActor(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	v, err := svc.CreateControl(domain.Control{Name: "c", Owner: "owner"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = svc.AdvanceControl(v.ID, domain.StatusActive, "operator")
	if err != nil {
		t.Fatal(err)
	}
	a := svc.AuditTrail(1)
	if len(a) != 1 || a[0].Actor != "operator" {
		t.Fatalf("audit=%#v", a)
	}
}

func TestBug008RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
