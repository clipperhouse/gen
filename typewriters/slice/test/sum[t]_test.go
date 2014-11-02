package main

import "testing"

func TestSumOther(t *testing.T) {
	things := ThingSlice{
		{"First", 60},
		{"Second", -20},
		{"Third", 100},
	}

	number := func(x Thing) Other {
		return x.Number
	}

	sum1 := things.SumOther(number)

	if sum1 != 140 {
		t.Errorf("SumOther should result in %v, got %v", 340, sum1)
	}
}
