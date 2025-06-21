package ratelimit

import (
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type RateLimiter struct {
	sketches   []*countMinSketch
	currentIdx int
	limit      int
	window     time.Duration
	ticker     *time.Ticker
	opts       *SketchOpts
	logger     *slog.Logger
	mu         sync.RWMutex
}

func New(sketchOpts *SketchOpts, limit int, window time.Duration, numWin int) *RateLimiter {
	rate := &RateLimiter{
		window:     window,
		ticker:     time.NewTicker(window),
		opts:       sketchOpts,
		limit:      limit,
		currentIdx: 0,
	}

	rate.sketches = make([]*countMinSketch, numWin)
	for i := range numWin {
		cms, _ := NewCMS(sketchOpts)
		rate.sketches[i] = cms
	}

	go rate.Rotate()
	return rate
}

func (r *RateLimiter) Rotate() {
	for range r.ticker.C {
		r.mu.Lock()
		r.currentIdx = (r.currentIdx + 1) % len(r.sketches)

		// Create new sketch for the new current window
		cms, _ := NewCMS(r.opts)
		r.sketches[r.currentIdx] = cms
		r.mu.Unlock()
		fmt.Printf("sketch rotated to window %d\n", r.currentIdx)
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

	return count <= uint64(r.limit)
}
