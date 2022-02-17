package utils

import "testing"

func AssertStr(a string, b string, t *testing.T) {
	if a != b {
		t.Fatalf("(%s != %s) assertion failed", a, b)
	}
}

func AssertInt(a int, b int, t *testing.T) {
	if a != b {
		t.Fatalf("(%d != %d) assertion failed", a, b)
	}
}
