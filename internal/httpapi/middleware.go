package httpapi

import (
	"net/http"
	"time"

	"github.com/local/compliance-evidence-chain/internal/platform"
)

func accessLog(log *platform.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.Context().Err(); err != nil {
			return
		}
		started := time.Now()
		next.ServeHTTP(w, r)
		log.Info("http request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(started).String())
	})
}

func recoverPanic(log *platform.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Error("panic recovered", "value", recovered)
				if w.Header().Get("Content-Type") == "" {
					w.Header().Set("Content-Type", "application/json")
				}
				writeError(w, http.StatusInternalServerError, ErrPanic)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
