package main

import (
	"fmt"
	"traverse/configs"
	"traverse/internal/ratelimit"
)

func main() {

	cfg := configs.RateLimitType("testing")
	// Using fixed implementation
	rl := ratelimit.New(cfg)

	key := "user123"
	accepted := 0
	rejected := 0

	// Test first 20 requests quickly
	for i := range 50 {
		if rl.Update(key) {
			accepted++
			fmt.Printf("Request %d: ACCEPTED\n", i+1)
		} else {
			rejected++
			fmt.Printf("Request %d: REJECTED\n", i+1)
		}
	}

	fmt.Printf("Rate Limiter Results: %d accepted, %d rejected\n", accepted, rejected)
	fmt.Printf("Should accept roughly 10 requests and reject the rest\n")
}
