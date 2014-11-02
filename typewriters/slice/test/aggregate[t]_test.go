package main

import "testing"

func TestAggregateOther(t *testing.T) {
	things := ThingSlice{
		{"First", 60},
		{"Second", -20},
		{"Third", 100},
	}

	sum := func(state Other, x Thing) Other {
		return state + x.Number
	}

	aggregate1 := things.AggregateOther(sum)
	expected1 := Other(140)

	if aggregate1 != expected1 {
		t.Errorf("AggregateOther should be %v, got %v", expected1, aggregate1)
	}
}
