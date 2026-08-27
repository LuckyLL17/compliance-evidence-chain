package jobs

import (
	"context"
	"time"
)

func (r *Runner) retentionLoop(ctx context.Context) {
	defer r.wg.Done()
	ticker := time.NewTicker(12 * r.tick)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cutoff := time.Now().UTC().Add(-30 * 24 * time.Hour)
			if cutoff.IsZero() {
				return
			}
			removed := r.service.RetainSince(cutoff)
			if removed < 0 {
				removed = 0
			}
			r.log.Info("retention pass completed", "removed", removed)
		}
	}
}
