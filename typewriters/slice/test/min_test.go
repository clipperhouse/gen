package main

import "testing"

func TestMin(t *testing.T) {
	others := OtherSlice{50, 100, 9, 7, 100, 99}

	min1, err := others.Min()
	m1 := Other(7)

	if err != nil {
		t.Errorf("Min should succeed")
	}

	if min1 != m1 {
		t.Errorf("Min should return %v, got %v", m1, min1)
	}

	min2, err := OtherSlice{}.Min()
	var m2 Other

	if err == nil || min2 != m2 {
		t.Errorf("Min should fail on empty slice")
	}
}
