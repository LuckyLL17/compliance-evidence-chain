// source-marker: internal/app/workflow.go
// source-marker: internal/app/commands.go
// source-marker: internal/app/audit.go
package verification

import (
	"strings"
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug001WorkflowInputs(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	got, err := svc.RunWorkflow(app.WorkflowRequest{Name: "collect", Actor: "alice", Mode: "evidence", Inputs: []domain.ID{"one", "two", "three"}})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"one", "two", "three"} {
		if !strings.Contains(got.Payload, want) {
			t.Fatalf("payload %q misses %s", got.Payload, want)
		}
	}
}

func TestBug001RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
