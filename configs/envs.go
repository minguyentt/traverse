package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"github.com/minguyentt/traverse/internal/assert"
)

func getEnv(key string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		return fmt.Sprintf("couldn't get env variable for %s", key)
	}

	return val
}

func getEnvAsInt(key string) int {
	val, exists := os.LookupEnv(key)
	assert.Assert(exists == false, "couldn't look up env variable", key)

	i, err := strconv.Atoi(val)
	assert.Assert(err != nil, "should never fail to parse str to int", key)

	return i
}

func getEnvAsBool(key string) bool {
	val, exists := os.LookupEnv(key)
	assert.Assert(exists == false, "couldn't look up env variable", key)

	bool, err := strconv.ParseBool(val)
	assert.Assert(err != nil, "should never fail to parse str to bool", key)

	return bool
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
