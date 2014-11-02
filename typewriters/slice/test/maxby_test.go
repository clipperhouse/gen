package main

import "testing"

func TestMaxBy(t *testing.T) {
	things := ThingSlice{
		{"First", 60},
		{"Second", -20},
		{"Third", 100},
		{"Fourth", 10},
	}

	max1, err1 := things.MaxBy(func(a, b Thing) bool {
		return a.Number < b.Number
	})

	if err1 != nil {
		t.Errorf("MaxBy Number should succeed")
	}

	expected1 := Thing{"Third", 100}

	if max1 != expected1 {
		t.Errorf("MaxBy Number should return %v, got %v", expected1, max1)
	}

	_, err2 := ThingSlice{}.MaxBy(func(a, b Thing) bool {
		return true
	})

	if err2 == nil {
		t.Errorf("MaxBy Number should fail on empty slice")
	}
}
