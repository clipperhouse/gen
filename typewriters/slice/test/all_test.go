package main

import "testing"

func TestAll(t *testing.T) {
	things := ThingSlice{
		{"First", 60},
		{"Second", -20},
		{"Third", 100},
	}

	all1 := things.All(func(x Thing) bool {
		return x.Name == "First"
	})

	if all1 {
		t.Errorf("All should be false for Name == 'First'")
	}

	all2 := things.All(func(x Thing) bool {
		return x.Name == "First" || x.Name == "Second" || x.Name == "Third"
	})

	if !all2 {
		t.Errorf("All should be true")
	}

	all3 := ThingSlice{}.All(func(x Thing) bool {
		return false
	})

	if !all3 {
		t.Errorf("All should evaulate true for empty slices")
	}
}
