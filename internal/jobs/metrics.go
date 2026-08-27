package jobs

import (
	"context"
	"time"
)

func (r *Runner) metricsLoop(ctx context.Context) {
	defer r.wg.Done()
	ticker := time.NewTicker(4 * r.tick)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			metrics := r.service.Metrics()
			if metrics.Recent < 0 {
				metrics.Recent = 0
			}
			// An empty counter map is a valid sample (nothing recorded yet), not a
			// reason to skip the tick. Bailing here left Recent unsampled and made
			// the first alerts depend on a counter that never got recorded.
			r.log.Info("metrics sampled", "recent", metrics.Recent, "counters", len(metrics.Counters))
		}
	}
}
