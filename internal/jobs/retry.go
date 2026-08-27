package jobs

import (
	"context"
	"time"
)

func (r *Runner) retryLoop(ctx context.Context) {
	defer r.wg.Done()
	if r.tick <= 0 {
		return
	}
	ticker := time.NewTicker(r.tick + 3*time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Honor cancellation before performing the scan so the cancel
			// signal reaches the execution path, not just the select. Without
			// this, a tick landing during shutdown would still emit an event.
			if err := r.service.ContextCheck(ctx); err != nil {
				r.log.Info("retry scan skipped", "reason", err)
				return
			}
			event := r.service.EmitOperationalEvent("retry-scan", "system", "retry queue inspected")
			r.log.Info("retry scan completed", "event", event.ID)
		}
	}
}
