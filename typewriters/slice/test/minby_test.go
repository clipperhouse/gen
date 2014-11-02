package main

import "testing"

func TestMinBy(t *testing.T) {
	things := ThingSlice{
		{"First", 60},
		{"Second", -20},
		{"Third", 100},
	}

	min1, err1 := things.MinBy(func(a, b Thing) bool {
		return a.Number < b.Number
	})

	if err1 != nil {
		t.Errorf("MinBy Number should succeed")
	}

	expected1 := Thing{"Second", -20}
	if min1 != expected1 {
		t.Errorf("MinBy Number should return %v, got %v", expected1, min1)
	}

	_, err2 := ThingSlice{}.MinBy(func(a, b Thing) bool {
		return true
	})

	if err2 == nil {
		t.Errorf("MinBy Number should fail on empty slice")
	}
}
