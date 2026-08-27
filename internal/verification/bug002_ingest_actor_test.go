package verification

import (
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug002IngestActor(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	got, err := svc.Ingest(app.IngestEnvelope{Kind: "evidence", Actor: "collector-7", Payload: "object-1"})
	if err != nil {
		t.Fatal(err)
	}
	if got.Actor != "collector-7" {
		t.Fatalf("event actor = %q", got.Actor)
	}
	trail := svc.AuditTrail(1)
	if len(trail) != 1 || trail[0].Actor != "collector-7" {
		t.Fatalf("audit trail = %#v", trail)
	}
}

func TestBug002RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
