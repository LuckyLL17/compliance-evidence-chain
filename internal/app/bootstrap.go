package app

import (
	"fmt"

	"github.com/local/compliance-evidence-chain/internal/domain"
)

func (s *Service) Bootstrap(actor string) error {
	if actor == "" {
		return fmt.Errorf("%w: actor is required", ErrInvalidCommand)
	}
	if err := s.Require(actor, "bootstrap"); err != nil {
		return err
	}
	if _, err := s.RunWorkflow(WorkflowRequest{
		Name:   "initial-domain-bootstrap",
		Actor:  actor,
		Inputs: []domain.ID{domain.ID("bootstrap")},
		Mode:   "",
	}); err != nil {
		return fmt.Errorf("bootstrap workflow: %w", err)
	}
	return nil
}
