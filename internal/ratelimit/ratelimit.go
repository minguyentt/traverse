package ratelimit

import (
	"log/slog"
	"sync"
	"time"
	"traverse/configs"
)

type RateLimiter struct {
	sketches []*countMinSketch
	buckets  uint
	depth    uint

	limit  int
	window time.Duration
	numWin int

	ticker *time.Ticker
	logger *slog.Logger
	mu     sync.RWMutex

	currentIdx int
}

func New(opts *configs.RateLimitConfig, logger *slog.Logger) *RateLimiter {
	if opts == nil {
		return runTester()
	}

	r := &RateLimiter{
		ticker:     time.NewTicker(opts.Window),
		buckets:    opts.Buckets,
		depth:      opts.Depth,
		limit:      opts.Limit,
		window:     opts.Window,
		numWin:     opts.NumWin,
		logger:     logger,
		currentIdx: 0,
	}

	r.sketches = make([]*countMinSketch, opts.NumWin)
	for i := range opts.NumWin {
		cms, _ := NewCMS(opts.Buckets, opts.Depth)
		r.sketches[i] = cms
	}

	go r.Rotate()

	return r
}

func runTester() *RateLimiter {
	rl := &RateLimiter{
		ticker:     time.NewTicker(time.Second),
		buckets:    1000,
		depth:      3,
		limit:      10,
		window:     time.Second,
		numWin:     3,
		currentIdx: 0,
	}

	rl.sketches = make([]*countMinSketch, rl.numWin)
	for i := range rl.numWin {
		cms, _ := NewCMS(rl.buckets, rl.depth)
		rl.sketches[i] = cms
	}

	go rl.Rotate()

	return rl
}

func (r *RateLimiter) Rotate() {
	for range r.ticker.C {
		r.mu.Lock()
		r.currentIdx = (r.currentIdx + 1) % len(r.sketches)

		// Create new sketch for the new current window
		cms, _ := NewCMS(r.buckets, r.depth)
		r.sketches[r.currentIdx] = cms
		r.mu.Unlock()
		r.logger.Warn("sketch rotation", "window", r.currentIdx)
	}
}

func (r *RateLimiter) Update(key string) bool {
	r.mu.Lock()
	// Only update current window
	r.sketches[r.currentIdx].Update([]byte(key), 1)
	r.mu.Unlock()

	r.mu.RLock()
	// Sum estimates from all windows
	count := uint64(0)
	for _, s := range r.sketches {
		count += s.Estimate([]byte(key))
	}
	r.mu.RUnlock()
	r.logger.Warn("estimated count update", "key", key, "count", count)
	return count <= uint64(r.limit)
}
