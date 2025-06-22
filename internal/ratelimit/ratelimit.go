package ratelimit

import (
	"log/slog"
	"sync"
	"time"

	"github.com/minguyentt/traverse/configs"
)

type RateLimiter struct {
	sketches []*countMinSketch
	buckets  uint
	depth    uint

	limit  int
	window time.Duration
	numWin int

	ticker *time.Ticker
	mu     sync.RWMutex
	logger *slog.Logger

	currentIdx int
}

func New(cfg *configs.RateLimitOpts) *RateLimiter {
	l := slog.Default()
	r := &RateLimiter{
		ticker:     time.NewTicker(cfg.Window),
		buckets:    cfg.Buckets,
		depth:      cfg.Depth,
		limit:      cfg.Limit,
		window:     cfg.Window,
		numWin:     cfg.NumWin,
		logger: l,
		currentIdx: 0,
	}

	r.sketches = make([]*countMinSketch, cfg.NumWin)
	for i := range cfg.NumWin {
		cms, _ := NewCMS(cfg.Buckets, cfg.Depth)
		r.sketches[i] = cms
	}

	go r.Rotate()

	return r
}

func (r *RateLimiter) Rotate() {
	for range r.ticker.C {
		r.mu.Lock()
		r.currentIdx = (r.currentIdx + 1) % len(r.sketches)

		// Create new sketch for the new current window
		// ISSUE: this causes memory overhead
		// cms, _ := NewCMS(r.buckets, r.depth)
		// r.sketches[r.currentIdx] = cms

		r.sketches[r.currentIdx].Reset()
		r.mu.Unlock()
		r.logger.Warn("sketch rotation", "window", r.currentIdx)
	}
}

func (c *countMinSketch) Reset() {
	for i := range c.counter {
		for j := range c.counter[i] {
			c.counter[i][j] = 0
		}
	}
}

func (r *RateLimiter) Update(key string) bool {
	r.mu.Lock()
	// Only update current window
	r.sketches[r.currentIdx].Update([]byte(key), 1)
	r.mu.Unlock()

	// Sum estimates from all windows
	count := uint64(0)
	for _, s := range r.sketches {
		count += s.Estimate([]byte(key))
	}

	return count <= uint64(r.limit)
}

func (r *RateLimiter) StopTicker() {
	r.ticker.Stop()
}
