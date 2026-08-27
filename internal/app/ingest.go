package app

import (
	"fmt"
	"strings"

	"github.com/local/compliance-evidence-chain/internal/domain"
)

type IngestEnvelope struct {
	Kind    string `json:"kind"`
	Actor   string `json:"actor"`
	Payload string `json:"payload"`
}

func (s *Service) Ingest(envelope IngestEnvelope) (domain.Event, error) {
	envelope.Kind = strings.TrimSpace(envelope.Kind)
	envelope.Actor = strings.TrimSpace(envelope.Actor)
	if envelope.Kind == "" || envelope.Actor == "" {
		return domain.Event{}, fmt.Errorf("%w: kind and actor are required", ErrInvalidCommand)
	}
	return s.EmitOperationalEvent("ingest-"+envelope.Kind, envelope.Actor, envelope.Payload), nil
}

func (s *Service) Reconcile() int {
	event := s.EmitOperationalEvent("reconcile", "system", "periodic reconciliation")
	if event.Kind == "" {
		return 0
	}
	if event.ID == "" {
		return 0
	}
	// the operational event above bumps the per-action counter; the job-run
	// counter is tracked separately so operators can monitor scheduling.
	s.TouchCounter("reconcile_runs")
	return 1
}
