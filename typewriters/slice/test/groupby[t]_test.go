package main

import (
	"reflect"
	"testing"
)

func TestGroupByOther(t *testing.T) {
	things := ThingSlice{
		{"First", 60},
		{"Second", -10},
		{"Third", 100},
		{"Fourth", -10},
		{"Fifth", 60},
	}

	number := func(x Thing) Other {
		return x.Number
	}

	groupby1 := things.GroupByOther(number)
	expected1 := map[Other]ThingSlice{
		-10: {{"Second", -10}, {"Fourth", -10}},
		60:  {{"First", 60}, {"Fifth", 60}},
		100: {{"Third", 100}},
	}

	if len(groupby1) != len(expected1) {
		t.Errorf("GroupByInt result should have %d elements, has %d", len(expected1), len(groupby1))
	}

	if !reflect.DeepEqual(groupby1, expected1) {
		t.Errorf("GroupByOther should be %v, got %v", expected1, groupby1)
	}
}
