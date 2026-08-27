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
	// Derive the workers' context from the parent so that cancellation of the
	// parent (signal) context propagates to every loop immediately. Previously
	// this was rooted at context.Background(), which detached the workers from
	// the caller's lifecycle and left retryLoop emitting events after shutdown.
	workCtx, cancel := context.WithCancel(ctx)
	r.cancel = cancel
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
