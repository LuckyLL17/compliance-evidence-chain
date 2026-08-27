// source-marker: internal/httpapi/middleware.go
// source-marker: internal/httpapi/response.go
// source-marker: internal/httpapi/router.go
package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

func TestBug019PanicResponse(t *testing.T) {
	r := recoverPanic(platform.NewLogger(), http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(200); panic("boom") }))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	if rr.Code != 500 {
		t.Fatalf("code=%d body=%s", rr.Code, rr.Body.String())
	}
}

func TestBug019RegressionHealth(t *testing.T) {
	svc := app.NewService(platform.RealClock{}, platform.NewLogger())
	if got := svc.Health()["status"]; got != "ok" {
		t.Fatalf("health status = %v", got)
	}
}
