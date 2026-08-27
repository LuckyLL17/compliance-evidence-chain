// source-marker: internal/app/connectorRun_service.go
// source-marker: internal/domain/common.go
// source-marker: internal/app/audit.go
package verification

import (
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug010ConnectorInvalidAudit(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	v, err := svc.CreateConnectorRun(domain.ConnectorRun{Name: "r", Owner: "o"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = svc.AdvanceConnectorRun(v.ID, domain.StatusCompleted, "op"); err == nil {
		t.Fatal("expected invalid transition")
	}
	if got := len(svc.AuditTrail(0)); got != 1 {
		t.Fatalf("audit count=%d", got)
	}
}

func TestBug010RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
