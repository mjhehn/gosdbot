package main

import (
	"testing"
)

func TestRollDie(t *testing.T) {
	for i := 0; i < 12; i++ {
		roll := RollDie(int64(6))
		if roll < 1 || roll > 6 {
			t.Errorf("RollDie is returning numbers out of range.")
		}
	}
}
