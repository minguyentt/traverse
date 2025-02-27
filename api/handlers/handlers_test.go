package handlers_test

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestRegistrationHandler(t *testing.T) {
	hash := "uwuraichubooboo"

	hashed := sha256.Sum256([]byte(hash))
	fmt.Println(hashed)
}
