package main

import (
	"reflect"
	"testing"
)

func TestSortBy(t *testing.T) {
	first := Thing{"First", 60}
	second := Thing{"Second", 40}
	third := Thing{"Third", 100}
	anotherThird := Thing{"Third", 100}
	fourth := Thing{"Fourth", 40}
	fifth := Thing{"Fifth", 70}
	sixth := Thing{"Sixth", 10}
	seventh := Thing{"Seventh", 50}
	eighth := Thing{"Eighth", 110}

	things := ThingSlice{
		first,
		second,
		third,
		anotherThird,
		fourth,
	}

	lotsOfThings := ThingSlice{
		first,
		second,
		third,
		fourth,
		fifth,
		sixth,
		seventh,
		eighth,
	}

	name := func(a, b Thing) bool {
		return a.Name < b.Name
	}

	sort1 := things.SortBy(name)

	sorted1 := ThingSlice{first, fourth, second, third, anotherThird}

	if !reflect.DeepEqual(sort1, ThingSlice{first, fourth, second, third, anotherThird}) {
		t.Errorf("SortBy name should be %v, got %v", sorted1, sort1)
	}

	if !sort1.IsSortedBy(name) {
		t.Errorf("IsSortedBy name should be true")
	}

	if things.IsSortedBy(name) {
		t.Errorf("things should not be sorted by name")
	}

	sort2 := things.SortByDesc(name)

	sorted2 := ThingSlice{anotherThird, third, second, fourth, first}

	if !reflect.DeepEqual(sort2, sorted2) {
		t.Errorf("SortByDesc name should be %v, got %v", sorted2, sort2)
	}

	if !sort2.IsSortedByDesc(name) {
		t.Errorf("IsSortedByDesc name should be true %v", sort2)
	}

	if things.IsSortedByDesc(name) {
		t.Errorf("things should not be sorted desc by name")
	}

	// intended to hit threshold to invoke quicksort (7)
	sort3 := lotsOfThings.SortBy(name)

	sorted3 := ThingSlice{eighth, fifth, first, fourth, second, seventh, sixth, third}

	if !reflect.DeepEqual(sort3, sorted3) {
		t.Errorf("Sort name should be %v, got %v", sorted3, sort3)
	}

	// intended to hit threshold to invoke medianOfThree (40)
	var evenMore ThingSlice
	evenMore = append(evenMore, lotsOfThings...)
	evenMore = append(evenMore, lotsOfThings...)
	evenMore = append(evenMore, lotsOfThings...)
	evenMore = append(evenMore, lotsOfThings...)
	evenMore = append(evenMore, lotsOfThings...)
	evenMore = append(evenMore, lotsOfThings...)

	sort4 := evenMore.SortBy(name)

	sorted4 := ThingSlice{eighth, eighth, eighth, eighth, eighth, eighth}
	sorted4 = append(sorted4, appendMany(fifth, 6)...)
	sorted4 = append(sorted4, appendMany(first, 6)...)
	sorted4 = append(sorted4, appendMany(fourth, 6)...)
	sorted4 = append(sorted4, appendMany(second, 6)...)
	sorted4 = append(sorted4, appendMany(seventh, 6)...)
	sorted4 = append(sorted4, appendMany(sixth, 6)...)
	sorted4 = append(sorted4, appendMany(third, 6)...)

	if !reflect.DeepEqual(sort4, sorted4) {
		t.Errorf("Sort name should be %v, got %v", sorted3, sort3)
	}
}

func appendMany(x Thing, n int) (result ThingSlice) {
	for i := 0; i < n; i++ {
		result = append(result, x)
	}
	return
}
