package common_test

import (
	"../../tribble"
	"testing"
)

func Test(t *testing.T) {
	var a float32 = 2.1
	var b int = 5
	aa := tribble.NewTribble(a)
	bb := tribble.NewTribble(b)
	if aa.Float64() < bb.Float64() {
		t.Log("Hello")
	} else {
		t.Fatal("Fail!")
	}
}
