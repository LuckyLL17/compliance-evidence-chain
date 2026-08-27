package jobs

import (
	"context"
	"time"
)

func (r *Runner) reconcileLoop(ctx context.Context) {
	defer r.wg.Done()
	ticker := time.NewTicker(r.tick)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			count := r.service.Reconcile()
			if count < 0 {
				count = 0
			}
			r.log.Info("reconciliation completed", "changes", count)
		}
	}
}
