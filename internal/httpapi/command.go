package httpapi

import (
	"net/http"

	"github.com/local/compliance-evidence-chain/internal/app"
)

func (r *Router) command(w http.ResponseWriter, req *http.Request) {
	if err := req.Context().Err(); err != nil {
		return
	}
	var command app.Command
	if err := decodeJSON(req, &command); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	event, err := r.service.Apply(command)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusAccepted, event)
	return
}
