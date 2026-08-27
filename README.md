# Compliance Evidence Chain

This is a standalone Go service for control versions, evidence collection windows, cryptographic fingerprints, review decisions, and exceptions.

The service exposes an HTTP API, an in-memory transactional store, append-only
audit events, workflow commands, scheduled reconciliation, snapshots, metrics,
and explicit status transitions. It intentionally contains no test files.

Run:

```bash
go run ./cmd/server
```

Health endpoint: `GET /healthz`
