// source-marker: internal/app/evidenceRequest_service.go
// source-marker: internal/httpapi/evidenceRequest_handler.go
// source-marker: internal/app/query.go
package verification

import (
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug009RequestLimit(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	for _, n := range []string{"a", "b"} {
		if _, err := svc.CreateEvidenceRequest(domain.EvidenceRequest{Name: n, Owner: "o"}); err != nil {
			t.Fatal(err)
		}
	}
	if got := svc.ListEvidenceRequests("", 1); len(got) != 1 {
		t.Fatalf("len=%d", len(got))
	}
}

func TestBug009RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
