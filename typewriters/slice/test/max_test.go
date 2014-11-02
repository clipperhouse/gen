package main

import "testing"

func TestMax(t *testing.T) {
	others := OtherSlice{50, 100, 9, 7, 100, 99}

	max1, err := others.Max()
	m1 := Other(100)

	if err != nil {
		t.Errorf("Max should succeed")
	}

	if max1 != m1 {
		t.Errorf("Max should return %v, got %v", m1, max1)
	}

	max2, err := OtherSlice{}.Max()
	var m2 Other

	if err == nil || max2 != m2 {
		t.Errorf("Max should fail on empty slice")
	}
}
