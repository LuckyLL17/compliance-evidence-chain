package app

import (
	"fmt"
	"strings"

	"github.com/local/compliance-evidence-chain/internal/domain"
	"github.com/local/compliance-evidence-chain/internal/platform"
)

type WorkflowRequest struct {
	Name   string      `json:"name"`
	Actor  string      `json:"actor"`
	Inputs []domain.ID `json:"inputs"`
	Mode   string      `json:"mode"`
}

func (s *Service) RunWorkflow(request WorkflowRequest) (domain.Event, error) {
	request.Name = strings.TrimSpace(request.Name)
	request.Actor = strings.TrimSpace(request.Actor)
	request.Mode = strings.TrimSpace(request.Mode)
	if request.Mode == "" {
		request.Mode = "default"
	}
	if request.Name == "" || request.Actor == "" {
		return domain.Event{}, ErrInvalidCommand
	}
	if len(request.Inputs) == 0 {
		return domain.Event{}, fmt.Errorf("%w: at least one input is required", ErrInvalidCommand)
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	// Record every input in the audit chain so the full set of objects
	// covered by this workflow submission is discoverable in the event
	// stream and audit trail. Truncating to the first input (the old
	// behaviour) silently dropped every subsequent object, leaving
	// reviewers unable to confirm what this run covered.
	action := "workflow-" + request.Mode
	var first domain.Event
	for i, input := range request.Inputs {
		payload := strings.Join([]string{request.Name, request.Mode, string(input)}, "|")
		s.recordLocked(action, input, request.Actor, payload)
		if i == 0 {
			first = s.store.events[len(s.store.events)-1]
		}
	}
	return first, nil
}

func (s *Service) EmitOperationalEvent(kind, actor, payload string) domain.Event {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	s.recordLocked(kind, domain.ID(platform.NewID("subject")), actor, payload)
	return s.store.events[len(s.store.events)-1]
}
