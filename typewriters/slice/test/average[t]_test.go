package main

import "testing"

func TestAverageOther(t *testing.T) {
	things := ThingSlice{
		{"First", 60},
		{"Second", -10},
		{"Third", 100},
	}

	number := func(x Thing) Other {
		return x.Number
	}

	average1, err := things.AverageOther(number)

	if err != nil {
		t.Errorf("Average should succeed")
	}

	expected1 := Other(50)

	if average1 != expected1 {
		t.Errorf("Average should be %v, got %v", expected1, average1)
	}

	average2, err := ThingSlice{}.AverageOther(number)

	if err == nil || average2 != 0 {
		t.Errorf("Average should fail on empty slice")
	}
}
