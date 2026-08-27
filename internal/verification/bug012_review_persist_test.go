// source-marker: internal/app/reviewDecision_service.go
// source-marker: internal/domain/review_decisions.go
// source-marker: internal/httpapi/reviewDecision_handler.go
package verification

import (
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug012ReviewPersist(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	v, err := svc.CreateReviewDecision(domain.ReviewDecision{Name: "d", Owner: "o"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = svc.AdvanceReviewDecision(v.ID, domain.StatusActive, "reviewer"); err != nil {
		t.Fatal(err)
	}
	got, _ := svc.GetReviewDecision(v.ID)
	if got.Status != domain.StatusActive || got.Version != 2 {
		t.Fatalf("got=%#v", got)
	}
}

func TestBug012RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
