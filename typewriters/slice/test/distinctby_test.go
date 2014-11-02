package main

import (
	"reflect"
	"testing"
)

func TestDistinctBy(t *testing.T) {
	things := ThingSlice{
		{"First", 0},
		{"Second", 9},
		{"First", 4},
		{"Third", 9},
		{"Fourth", 5},
		{"Fifth", 4},
	}

	expected := ThingSlice{
		{"First", 0},
		{"Second", 9},
		{"First", 4},
		{"Fourth", 5},
	}

	distinctby1 := things.DistinctBy(func(a, b Thing) bool {
		return a.Number == b.Number
	})

	if !reflect.DeepEqual(distinctby1, expected) {
		t.Errorf("DistinctBy should be %v, but got %v", expected, distinctby1)
	}
}
