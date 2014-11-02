package main

import "testing"

func TestSingle(t *testing.T) {
	things := ThingSlice{
		{"First", 0},
		{"Second", 0},
		{"Third", 0},
		{"Second", 1},
	}

	s1, err1 := things.Single(func(x Thing) bool {
		return x.Name == "Third"
	})

	if err1 != nil {
		t.Errorf("Single should succeed when finding Name == Third")
	}

	expected1 := Thing{"Third", 0}
	if s1 != expected1 {
		t.Errorf("Single should find %v, but found %v", expected1, s1)
	}

	_, err2 := things.Single(func(x Thing) bool {
		return x.Name == "Second"
	})

	if err2 == nil {
		t.Errorf("Single should fail when finding Name == Second")
	}

	_, err3 := things.Single(func(x Thing) bool {
		return x.Name == "Dummy"
	})

	if err3 == nil {
		t.Errorf("Single should fail when finding Name == Dummy")
	}

	_, err4 := ThingSlice{}.First(func(x Thing) bool {
		return true
	})

	if err4 == nil {
		t.Errorf("Single should fail on empty slice")
	}
}
