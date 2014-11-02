package main

import "testing"

func TestSum(t *testing.T) {
	others := OtherSlice{50, 100, 9, 7, 100, 99}

	sum := others.Sum()

	if sum != 365 {
		t.Errorf("Sum should result in %v, got %v", 365, sum)
	}
}
