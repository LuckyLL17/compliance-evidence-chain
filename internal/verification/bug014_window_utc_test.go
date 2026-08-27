// source-marker: internal/platform/clock.go
// source-marker: internal/app/collectionWindow_service.go
// source-marker: internal/domain/collection_window.go
package verification

import (
	"testing"
	"time"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug014WindowUtc(t *testing.T) {
	clock := fixedClock{now: time.Date(2026, 8, 24, 12, 0, 0, 0, time.FixedZone("CST", 8*3600))}
	svc := app.NewService(clock, platform.NewLogger())
	v, err := svc.CreateCollectionWindow(domain.CollectionWindow{Name: "w", Owner: "o"})
	if err != nil {
		t.Fatal(err)
	}
	if v.CreatedAt.Location() != time.UTC {
		t.Fatalf("location=%v", v.CreatedAt.Location())
	}
}

type fixedClock struct{ now time.Time }

func (c fixedClock) Now() time.Time                { return c.now }
func (c fixedClock) Since(time.Time) time.Duration { return 0 }

func TestBug014RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
