// source-marker: internal/domain/exception_cases.go
// source-marker: internal/app/exceptionCase_service.go
// source-marker: internal/httpapi/exceptionCase_handler.go
package verification

import (
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug013ExceptionVersion(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	v, err := svc.CreateExceptionCase(domain.ExceptionCase{Name: "x", Owner: "o"})
	if err != nil {
		t.Fatal(err)
	}
	got, err := svc.AdvanceExceptionCase(v.ID, domain.StatusActive, "op")
	if err != nil {
		t.Fatal(err)
	}
	if got.Version != 2 || got.Status != domain.StatusActive {
		t.Fatalf("got=%#v", got)
	}
}

func TestBug013RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
