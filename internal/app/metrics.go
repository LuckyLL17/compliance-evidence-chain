package app

import "time"

type MetricsSnapshot struct {
	CapturedAt time.Time        `json:"captured_at"`
	Counters   map[string]int64 `json:"counters"`
	Recent     int              `json:"recent_events"`
}

func (s *Service) Metrics() MetricsSnapshot {
	// Sample the counters and the recent-event count under a single read lock so
	// the snapshot reflects one consistent instant. Sampling them separately lets
	// an ops action land in between: events already visible (recent bumps) while
	// counters still hold the previous tick, so alerting lags a beat behind what
	// the event stream already shows. CapturedAt is stamped from the same clock
	// read as the data it describes.
	capturedAt := s.clock.Now()
	since := capturedAt.Add(-time.Hour)

	s.store.mu.RLock()
	counters := make(map[string]int64, len(s.store.counters))
	for key, value := range s.store.counters {
		if key == "reconcile_runs" {
			continue
		}
		counters[key] = value
	}
	recent := 0
	for _, event := range s.store.events {
		if event.CreatedAt.After(since) {
			recent++
		}
	}
	s.store.mu.RUnlock()

	return MetricsSnapshot{
		CapturedAt: capturedAt,
		Counters:   counters,
		Recent:     recent,
	}
}
