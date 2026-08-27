// source-marker: internal/app/chainEvent_service.go
// source-marker: internal/domain/chain_events.go
// source-marker: internal/app/audit.go
package verification

import (
	"strings"
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug015ChainAuditKey(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	v, err := svc.CreateChainEvent(domain.ChainEvent{Name: "e", Owner: "o"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = svc.AdvanceChainEvent(v.ID, domain.StatusActive, "op"); err != nil {
		t.Fatal(err)
	}
	e := svc.EventStream(1)
	if len(e) != 1 || !strings.Contains(e[0].Payload, "active") {
		t.Fatalf("event=%#v", e)
	}
}

func TestBug015RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
