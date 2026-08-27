// source-marker: internal/app/query.go
// source-marker: internal/httpapi/search.go
// source-marker: internal/app/audit.go
package verification

import (
	"fmt"
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug020SearchLimit(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	for i := 0; i < 3; i++ {
		svc.EmitOperationalEvent("searchable", "a", fmt.Sprint(i))
	}
	if got := svc.Search("searchable", 1); len(got) != 1 {
		t.Fatalf("len=%d", len(got))
	}
}

func TestBug020RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
