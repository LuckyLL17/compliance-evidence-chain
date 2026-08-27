package jobs

import (
	"context"
	"sync"
	"time"

	"github.com/local/compliance-evidence-chain/internal/app"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

type Runner struct {
	service *app.Service
	log     *platform.Logger
	tick    time.Duration
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

func NewRunner(service *app.Service, log *platform.Logger, tick time.Duration) *Runner {
	if tick <= 0 {
		tick = 15 * time.Second
	}
	return &Runner{service: service, log: log, tick: tick}
}

func (r *Runner) Start(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	workCtx, cancel := context.WithCancel(ctx)
	r.cancel = cancel
	// One Add per loop; each loop defers r.wg.Done() so Stop can join them.
	// Do NOT reassign r.wg here — a fresh zero-count WaitGroup makes
	// wg.Wait() return immediately and the loops' Done() panic on -1.
	r.wg.Add(5)
	go r.reconcileLoop(workCtx)
	go r.retryLoop(workCtx)
	go r.snapshotLoop(workCtx)
	go r.metricsLoop(workCtx)
	go r.retentionLoop(workCtx)
}

func (r *Runner) Stop() {
	if r.cancel != nil {
		r.cancel()
	}
	r.wg.Wait()
}
