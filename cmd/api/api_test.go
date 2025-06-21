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
}

// func TestSlidingWindowBehavior(t *testing.T) {
// 	cfg := configs.Env.RATELIMIT.Test
// 	rl := ratelimit.New(cfg, nil)
// 	defer rl.StopTicker()
//
// 	key := "sliding_window"
//
// 	t.Run("Should fill up the limit", func(t *testing.T) {
// 		for i := range 5 {
// 			if accepted := rl.Update(key); !accepted {
// 				t.Errorf("expected request %d to be accepted", i+1)
// 			}
// 		}
// 	})
//
// 	t.Run("should be rejected", func(t *testing.T) {
// 		for i := range 3 {
// 			if accepted := rl.Update(key); accepted {
// 				t.Errorf("expected overflow request %d to be rejected", i+1)
// 			}
// 		}
// 	})
//
// 	t.Run("Should accept again after window slide", func(t *testing.T) {
// 		for i := range 4 {
// 			if accepted := rl.Update(key); !accepted {
// 				t.Errorf("expected request %d after sliding window to be accepted", i+1)
// 			}
// 		}
// 	})
// }

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
