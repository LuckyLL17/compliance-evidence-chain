package app

import (
	"sync"

	"github.com/local/compliance-evidence-chain/internal/domain"
)

type Store struct {
	mu                 sync.RWMutex
	frameworks         map[domain.ID]domain.Framework
	controls           map[domain.ID]domain.Control
	evidence_requests  map[domain.ID]domain.EvidenceRequest
	connector_runs     map[domain.ID]domain.ConnectorRun
	evidence_objects   map[domain.ID]domain.EvidenceObject
	review_decisions   map[domain.ID]domain.ReviewDecision
	exception_cases    map[domain.ID]domain.ExceptionCase
	collection_windows map[domain.ID]domain.CollectionWindow
	chain_events       map[domain.ID]domain.ChainEvent
	retention_policies map[domain.ID]domain.RetentionPolicy
	access_rules       map[domain.ID]domain.AccessRule
	events             []domain.Event
	audits             []domain.AuditRecord
	counters           map[string]int64
	chain              string
}

func NewStore() *Store {
	return &Store{
		frameworks:         make(map[domain.ID]domain.Framework),
		controls:           make(map[domain.ID]domain.Control),
		evidence_requests:  make(map[domain.ID]domain.EvidenceRequest),
		connector_runs:     make(map[domain.ID]domain.ConnectorRun),
		evidence_objects:   make(map[domain.ID]domain.EvidenceObject),
		review_decisions:   make(map[domain.ID]domain.ReviewDecision),
		exception_cases:    make(map[domain.ID]domain.ExceptionCase),
		collection_windows: make(map[domain.ID]domain.CollectionWindow),
		chain_events:       make(map[domain.ID]domain.ChainEvent),
		retention_policies: make(map[domain.ID]domain.RetentionPolicy),
		access_rules:       make(map[domain.ID]domain.AccessRule),
		events:             make([]domain.Event, 0, 64),
		audits:             make([]domain.AuditRecord, 0, 64),
		counters:           make(map[string]int64),
	}
}
