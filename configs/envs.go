package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func getEnv(key string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		return fmt.Sprintf("couldn't get env variable for %s", key)
	}

	return val
}

func getEnvAsInt32(key string, fb int32) int32 {
	if val, exists := os.LookupEnv(key); exists {
		parsed, err := strconv.Atoi(val)
		if err != nil {
			return fb
		}

		return int32(parsed)
	}

	return fb
}

func getEnvAsTime(key string, fb time.Duration) time.Duration {
	val, exists := os.LookupEnv(key)
	if exists {
		parsed, err := time.ParseDuration(val)
		if err != nil {
			return fb
		}

		return parsed
	}

	return fb
}
