package utils

import (
	"testing"
	"time"
)

func AssertStr(a, b string, t *testing.T) {
	if a != b {
		t.Fatalf("(%s != %s) assertion failed", a, b)
	}
}

func AssertInt(a, b int, t *testing.T) {
	if a != b {
		t.Fatalf("(%d != %d) assertion failed", a, b)
	}
}

func AssertDuration(a, b time.Duration, t *testing.T) {
	if a != b {
		t.Fatalf("(%s != %s) assertion failed", a.String(), b.String())
	}
}

func AssertTrue(a bool, t *testing.T) {
	if !a {
		t.Fatalf("%+v is not true", a)
	}
}

func AssertFalse(a bool, t *testing.T) {
	if a {
		t.Fatalf("%+v is not false", a)
	}
}
