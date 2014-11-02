package main

import "testing"

func TestAny(t *testing.T) {
	things := ThingSlice{
		{"First", 60},
		{"Second", -20},
		{"Third", 100},
	}

	any1 := things.Any(func(x Thing) bool {
		return x.Name == "Dummy"
	})

	if any1 {
		t.Errorf("Any should not evaulate true for Name == Dummy")
	}

	any2 := things.Any(func(x Thing) bool {
		return x.Number > 50
	})

	if !any2 {
		t.Errorf("Any should evaulate true for Number > 50")
	}

	any3 := ThingSlice{}.Any(func(x Thing) bool {
		return true
	})

	if any3 {
		t.Errorf("Any should evaulate false for empty slices")
	}
}
