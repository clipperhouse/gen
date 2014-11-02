package main

import (
	"reflect"
	"testing"
)

func TestSelectOther(t *testing.T) {
	things := ThingSlice{
		{"First", 60},
		{"Second", -20},
		{"Third", 100},
	}

	number := func(x Thing) Other {
		return x.Number
	}

	select1 := things.SelectOther(number)
	expected1 := []Other{60, -20, 100}

	if !reflect.DeepEqual(select1, expected1) {
		t.Errorf("SelectOther should result in %v, got %v", expected1, select1)
	}
}
