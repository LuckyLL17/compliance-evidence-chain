// source-marker: internal/app/framework_service.go
// source-marker: internal/domain/frameworks.go
// source-marker: internal/httpapi/framework_handler.go
package verification

import (
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug007FrameworkAdvance(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	v, err := svc.CreateFramework(domain.Framework{Name: "f", Owner: "o"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = svc.AdvanceFramework(v.ID, domain.StatusActive, "reviewer"); err != nil {
		t.Fatal(err)
	}
	got, ok := svc.GetFramework(v.ID)
	if !ok || got.Status != domain.StatusActive || got.Version != 2 {
		t.Fatalf("got=%#v", got)
	}
}

func TestBug007RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
