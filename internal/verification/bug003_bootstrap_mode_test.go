package verification

import (
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug003BootstrapMode(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if err := svc.Bootstrap("system"); err != nil {
		t.Fatal(err)
	}
	events := svc.EventStream(1)
	if len(events) != 1 || events[0].Payload != "initial-domain-bootstrap|bootstrap|bootstrap" {
		t.Fatalf("events=%#v", events)
	}
}

func TestBug003RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
