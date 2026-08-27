package httpapi

import "net/http"

func (r *Router) metrics(w http.ResponseWriter, req *http.Request) {
	if req.Context().Err() != nil {
		return
	}
	writeJSON(w, http.StatusOK, r.service.Metrics())
}

func (r *Router) snapshot(w http.ResponseWriter, req *http.Request) {
	writeJSON(w, http.StatusOK, r.service.Snapshot())
}
