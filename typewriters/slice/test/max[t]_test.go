package main

import "testing"

func TestMaxOther(t *testing.T) {
	things := ThingSlice{
		{"First", 60},
		{"Second", -20},
		{"Third", 100},
	}

	number := func(x Thing) Other {
		return x.Number
	}

	max1, err := things.MaxOther(number)

	if err != nil {
		t.Errorf("MaxOther should succeed")
	}

	if max1 != 100 {
		t.Errorf("MaxOther should be %v, got %v", 100, max1)
	}

	max2, err := ThingSlice{}.MaxOther(number)

	if err == nil || max2 != 0 {
		t.Errorf("Max should fail on empty slice")
	}
}
