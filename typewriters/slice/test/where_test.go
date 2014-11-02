package main

import (
	"reflect"
	"testing"
)

func TestWhere(t *testing.T) {
	things := ThingSlice{
		{"First", 0},
		{"Second", 0},
		{"Third", 0},
		{"Second", 10},
	}

	where1 := things.Where(func(x Thing) bool {
		return x.Name == "Second"
	})

	expected1 := ThingSlice{
		{"Second", 0},
		{"Second", 10},
	}

	if !reflect.DeepEqual(where1, expected1) {
		t.Errorf("Where should result in %v, got %v", expected1, where1)
	}

	where2 := things.Where(func(x Thing) bool {
		return x.Name == "Dummy"
	})

	if len(where2) != 0 {
		t.Errorf("Where should result in empty slice, got %v", where2)
	}

	where3 := ThingSlice{}.Where(func(x Thing) bool {
		return true
	})

	if len(where3) != 0 {
		t.Errorf("Where should result in empty slice, got %v", where3)
	}
}
