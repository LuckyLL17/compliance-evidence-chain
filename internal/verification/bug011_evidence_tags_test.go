// source-marker: internal/domain/evidence_objects.go
// source-marker: internal/app/evidenceObject_service.go
// source-marker: internal/httpapi/evidenceObject_handler.go
package verification

import (
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug011EvidenceTags(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	v, err := svc.CreateEvidenceObject(domain.EvidenceObject{Name: "e", Owner: "o", Tags: []string{"pci", "fresh"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(v.Tags) != 2 || v.Tags[0] != "pci" {
		t.Fatalf("tags=%#v", v.Tags)
	}
}

func TestBug011RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
