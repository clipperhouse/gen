package main

import "testing"

func TestAverage(t *testing.T) {
	others := OtherSlice{50, 100, 9, 7, 100, 99}

	average1, err := others.Average()

	if err != nil {
		t.Errorf("Average should succeed")
	}

	avg1 := Other(60)

	if average1 != avg1 {
		t.Errorf("Average should be %v, got %v", avg1, average1)
	}

	average2, err := OtherSlice{}.Average()

	if err == nil || average2 != 0 {
		t.Errorf("Average should fail on empty slice")
	}
}
