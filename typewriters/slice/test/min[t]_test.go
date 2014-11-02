package main

import "testing"

func TestMinOther(t *testing.T) {
	things := ThingSlice{
		{"First", 60},
		{"Second", -20},
		{"Third", 100},
	}

	number := func(x Thing) Other {
		return x.Number
	}

	min1, err := things.MinOther(number)

	if err != nil {
		t.Errorf("MinOther should succeed")
	}

	if min1 != -20 {
		t.Errorf("MinOther should be %v, got %v", -20, min1)
	}

	min2, err := ThingSlice{}.MinOther(number)

	if err == nil || min2 != 0 {
		t.Errorf("MinOther should fail on empty slice")
	}
}
