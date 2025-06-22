package main

import (
	"fmt"
	"sync"
	"testing"

	"github.com/minguyentt/traverse/configs"
	"github.com/minguyentt/traverse/internal/ratelimit"
	"github.com/minguyentt/traverse/pkg/tests"
)

func TestCountMinSketchRateLimiter(t *testing.T) {
	cfg := configs.Env.RATELIMIT.Test
	rl := ratelimit.New(cfg)
	defer rl.StopTicker()

	key1 := "userA"
	key2 := "userB"

	for i := range cfg.Limit {
		accepted := rl.Update(key1)
		if !accepted {
			t.Errorf("expected request %d from keyA to be accepted", i+1)
		}
	}

	// Simulate key2 with potential hash collision
	accepted := rl.Update(key2)

	if !accepted {
		t.Log("possible false positive: keyB rejected due to hash collision with keyA")
	}
}

func TestMultipleUsers(t *testing.T) {
	cfg := configs.Env.RATELIMIT.Test
	rl := ratelimit.New(cfg)
	defer rl.StopTicker()

	users := []string{"alice", "bob", "charlie"}

	for _, user := range users {
		t.Run(fmt.Sprintf("user: %s", user), func(t *testing.T) {
			accepted := 0
			for range 5 {
				if rl.Update(user) {
					accepted++
				}
			}
			if accepted != 5 {
				t.Errorf("expected 5 accepted requests for %s, got %d", user, accepted)
			}
		})
	}
}

func TestConcurrency(t *testing.T) {
	cfg := configs.Env.RATELIMIT.Test
	rl := ratelimit.New(cfg)
	defer rl.StopTicker()

	var wg sync.WaitGroup
	var mu sync.Mutex
	totalAccepted := 0
	totalRejected := 0

	numGoroutines := 10
	reqPerGoroutine := 10

	for i := range numGoroutines {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			key := fmt.Sprintf("user_%d", goroutineID)
			localAccepted := 0
			localRejected := 0

			for range reqPerGoroutine {
				if rl.Update(key) {
					localAccepted++
				} else {
					localRejected++
				}
			}

			mu.Lock()
			totalAccepted += localAccepted
			totalRejected += localRejected
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	expectedTotal := numGoroutines * reqPerGoroutine
	tests.AssertEqual(t, totalAccepted+totalRejected, expectedTotal, "request mismatch")
}

// BenchmarkRateLimiter benchmarks the rate limiter performance
func BenchmarkRateLimiter(b *testing.B) {
	cfg := configs.Env.RATELIMIT.Test
	rl := ratelimit.New(cfg)
	defer rl.StopTicker()

	keys := make([]string, 100)
	for i := range keys {
		keys[i] = fmt.Sprintf("user_%d", i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := keys[i%len(keys)]
			rl.Update(key)
			i++
		}
	})
}
