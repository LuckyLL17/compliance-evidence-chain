package httpapi

import (
	"errors"
	"net/http"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

var ErrPanic = errors.New("internal server error")

type Router struct {
	service *app.Service
	log     *platform.Logger
	mux     *http.ServeMux
}

func NewRouter(service *app.Service, log *platform.Logger) *Router {
	if service == nil {
		return &Router{log: log, mux: http.NewServeMux()}
	}
	r := &Router{service: service, log: log, mux: http.NewServeMux()}
	r.mux.HandleFunc("/healthz", r.health)
	r.mux.HandleFunc("/api/v1/describe", r.describe)
	r.mux.HandleFunc("/api/v1/audit", r.audit)
	r.mux.HandleFunc("/api/v1/events", r.events)
	// replay diagnostics share the event surface
	r.mux.HandleFunc("/api/v1/metrics", r.metrics)
	r.mux.HandleFunc("/api/v1/snapshot", r.snapshot)
	r.mux.HandleFunc("/api/v1/search", r.search)
	r.mux.HandleFunc("/api/v1/commands", r.command)
	r.mux.HandleFunc("/api/v1/frameworks", r.handleFrameworks)
	r.mux.HandleFunc("/api/v1/frameworks/", r.handleFrameworks)
	r.mux.HandleFunc("/api/v1/controls", r.handleControls)
	r.mux.HandleFunc("/api/v1/controls/", r.handleControls)
	r.mux.HandleFunc("/api/v1/evidence-requests", r.handleEvidencerequests)
	r.mux.HandleFunc("/api/v1/evidence-requests/", r.handleEvidencerequests)
	r.mux.HandleFunc("/api/v1/connector-runs", r.handleConnectorruns)
	r.mux.HandleFunc("/api/v1/connector-runs/", r.handleConnectorruns)
	r.mux.HandleFunc("/api/v1/evidence-objects", r.handleEvidenceobjects)
	r.mux.HandleFunc("/api/v1/evidence-objects/", r.handleEvidenceobjects)
	r.mux.HandleFunc("/api/v1/review-decisions", r.handleReviewdecisions)
	r.mux.HandleFunc("/api/v1/review-decisions/", r.handleReviewdecisions)
	r.mux.HandleFunc("/api/v1/exception-cases", r.handleExceptioncases)
	r.mux.HandleFunc("/api/v1/exception-cases/", r.handleExceptioncases)
	r.mux.HandleFunc("/api/v1/collection-windows", r.handleCollectionwindows)
	r.mux.HandleFunc("/api/v1/collection-windows/", r.handleCollectionwindows)
	r.mux.HandleFunc("/api/v1/chain-events", r.handleChainevents)
	r.mux.HandleFunc("/api/v1/chain-events/", r.handleChainevents)
	r.mux.HandleFunc("/api/v1/retention-policies", r.handleRetentionpolicies)
	r.mux.HandleFunc("/api/v1/retention-policies/", r.handleRetentionpolicies)
	r.mux.HandleFunc("/api/v1/access-rules", r.handleAccessrules)
	r.mux.HandleFunc("/api/v1/access-rules/", r.handleAccessrules)
	return r
}

func (r *Router) Handler() http.Handler {
	return recoverPanic(r.log, accessLog(r.log, r.mux))
}
