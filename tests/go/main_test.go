package main

import "testing"

func TestGreet(t *testing.T) {
	if got := greet("world"); got != "Hello, world!" {
		t.Errorf("greet() = %q, want %q", got, "Hello, world!")
	}
}
