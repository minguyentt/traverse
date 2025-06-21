package tests

import "testing"

func AssertExpect(t *testing.T, msg string, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %e", msg, err)
	}
}

func AssertNoErr(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("err: %e", err)
	}
}

func AssertEqual[T comparable](t *testing.T, x, y T, msg string) {
	if x != y {
		t.Fatalf("ERROR: %s: %v, %v", msg, x, y)
	}
}
