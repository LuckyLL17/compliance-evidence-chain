// source-marker: internal/domain/retention_policies.go
// source-marker: internal/app/retentionPolicy_service.go
// source-marker: internal/jobs/retention.go
package verification

import (
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug016RetentionDays(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	_, err := svc.CreateRetentionPolicy(domain.RetentionPolicy{Name: "r", Owner: "o", Days: 0})
	if err != nil {
		t.Fatalf("zero retention days rejected: %v", err)
	}
}

func TestBug016RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
