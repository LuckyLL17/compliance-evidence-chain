package httpapi

import (
	"net/http"
	"strconv"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/domain"
)

func (r *Router) handleCollectionwindows(w http.ResponseWriter, req *http.Request) {
	parts := pathParts(req, "/api/v1/collection-windows")
	if len(parts) == 0 {
		if req.Method == http.MethodGet {
			limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
			writeJSON(w, http.StatusOK, r.service.ListCollectionWindows(req.URL.Query().Get("q"), limit))
			return
		}
		if req.Method == http.MethodPost {
			var input domain.CollectionWindow
			if err := decodeJSON(req, &input); err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			value, err := r.service.CreateCollectionWindow(input)
			if err != nil {
				writeError(w, http.StatusUnprocessableEntity, err)
				return
			}
			writeJSON(w, http.StatusCreated, value)
			return
		}
		writeError(w, http.StatusMethodNotAllowed, app.ErrInvalidCommand)
		return
	}
	id := domain.ID(parts[0])
	if len(parts) == 1 && req.Method == http.MethodGet {
		value, ok := r.service.GetCollectionWindow(id)
		if !ok {
			writeError(w, http.StatusNotFound, app.ErrNotFound)
			return
		}
		writeJSON(w, http.StatusOK, value)
		return
	}
	if len(parts) == 2 && parts[1] == "advance" && req.Method == http.MethodPost {
		var input struct {
			Status domain.Status `json:"status"`
			Actor  string        `json:"actor"`
		}
		if err := decodeJSON(req, &input); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		value, err := r.service.AdvanceCollectionWindow(id, input.Status, input.Actor)
		if err != nil {
			status := http.StatusUnprocessableEntity
			if err == app.ErrNotFound {
				status = http.StatusNotFound
			}
			writeError(w, status, err)
			return
		}
		writeJSON(w, http.StatusOK, value)
		return
	}
	writeError(w, http.StatusNotFound, app.ErrNotFound)
}
