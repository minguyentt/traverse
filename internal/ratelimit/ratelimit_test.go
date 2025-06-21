package ratelimit

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

//FIX: TESTING!!!!!
func TestSlidingWindowBehavior(t *testing.T) {
	fmt.Println("\n=== Testing Sliding Window Behavior ===")

	rl := New(nil)
	defer rl.ticker.Stop()

	key := "sliding_test"

	// Fill up the limit
	fmt.Println("Phase 1: Fill up the limit")
	for i := range 5 {
		accepted := rl.Update(key)
		fmt.Printf("Request %d: %v\n", i+1, accepted)
	}

	// Should be rejected now
	fmt.Println("\nPhase 2: Should be rejected")
	for i := range 3 {
		accepted := rl.Update(key)
		fmt.Printf("Overflow request %d: %v\n", i+1, accepted)
	}

	// Wait for window to slide
	fmt.Println("\nPhase 3: Waiting for window to slide...")
	time.Sleep(700 * time.Millisecond)

	// Should be accepted again
	fmt.Println("Phase 4: Should be accepted again after window slide")
	for i := range 4 {
		accepted := rl.Update(key)
		fmt.Printf("After rotation request %d: %v\n", i+1, accepted)
	}
}

// TestMultipleUsers tests behavior with different users
func TestMultipleUsers(t *testing.T) {
	fmt.Println("\n=== Testing Multiple Users ===")


	rl := New(nil)
	defer rl.ticker.Stop()

	users := []string{"alice", "bob", "charlie"}

	// Each user should be able to make their limit
	for _, user := range users {
		fmt.Printf("\nTesting user: %s\n", user)
		acceptedCount := 0

		for i := range 5 {
			if rl.Update(user) {
				acceptedCount++
				fmt.Printf("  Request %d: ACCEPTED\n", i+1)
			} else {
				fmt.Printf("  Request %d: REJECTED\n", i+1)
			}
		}

		fmt.Printf("User %s: %d requests accepted\n", user, acceptedCount)
	}
}

// TestConcurrency tests thread safety
func TestConcurrency(t *testing.T) {
	fmt.Println("\n=== Testing Concurrency ===")


	rl := New(nil)
	defer rl.ticker.Stop()

	var wg sync.WaitGroup
	var mu sync.Mutex
	totalAccepted := 0
	totalRejected := 0

	// Launch multiple goroutines
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

			fmt.Printf("Goroutine %d: %d accepted, %d rejected\n",
				goroutineID, localAccepted, localRejected)
		}(i)
	}

	wg.Wait()

	fmt.Printf("Concurrency test total: %d accepted, %d rejected out of %d requests\n",
		totalAccepted, totalRejected, numGoroutines*reqPerGoroutine)
}

// BenchmarkRateLimiter benchmarks the rate limiter performance
func BenchmarkRateLimiter(b *testing.B) {

	rl := New(nil)
	defer rl.ticker.Stop()

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
