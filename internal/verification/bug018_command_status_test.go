// source-marker: internal/httpapi/command.go
// source-marker: internal/app/commands.go
// source-marker: internal/app/audit.go
package verification

import (
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug018CommandStatus(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	e, err := svc.Apply(app.Command{Action: "rotate", Actor: "admin", Subject: "s"})
	if err != nil {
		t.Fatal(err)
	}
	if e.Kind != "command-rotate" || len(svc.EventStream(0)) != 1 {
		t.Fatalf("event=%#v", e)
	}
}

func TestBug018RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
