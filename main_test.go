package main

import (
	"testing"
	"time"
)

// Unit test

func TestComplexFunction(t *testing.T) {
	ch := make(chan int)
	go complexFunction(2, 3, ch)

	select {
	case result := <-ch:
		if result != 5 {
			t.Errorf("Expected 5 but got %d", result)
		}
	case <-time.After(3 * time.Second):
		t.Error("Test timed out")
	}
}
