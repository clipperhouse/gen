package main

import "testing"

func TestFirst(t *testing.T) {
	things := ThingSlice{
		{"First", 0},
		{"Second", 0},
		{"Third", 0},
	}

	f1, err1 := things.First(func(x Thing) bool {
		return x.Name == "Third"
	})

	if err1 != nil {
		t.Errorf("First should succeed when finding Name == Third")
	}

	expected1 := Thing{"Third", 0}
	if f1 != expected1 {
		t.Errorf("First should find %v, but found %v", expected1, f1)
	}

	_, err2 := things.First(func(x Thing) bool {
		return x.Name == "Dummy"
	})

	if err2 == nil {
		t.Errorf("First should fail when finding Name == Dummy")
	}

	_, err3 := ThingSlice{}.First(func(x Thing) bool {
		return true
	})

	if err3 == nil {
		t.Errorf("First should fail on empty slice")
	}
}
