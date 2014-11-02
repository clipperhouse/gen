package main

import (
	"reflect"
	"testing"
)

func TestSort(t *testing.T) {
	others := OtherSlice{50, 100, 9, 7, 100, 99}

	sort1 := others.Sort()
	s1 := OtherSlice{7, 9, 50, 99, 100, 100}

	if !reflect.DeepEqual(sort1, s1) {
		t.Errorf("Sort should result in %v, got %v", s1, sort1)
	}

	if !sort1.IsSorted() {
		t.Errorf("IsSorted should be true for %v", sort1)
	}

	sort2 := others.SortDesc()
	s2 := OtherSlice{100, 100, 99, 50, 9, 7}

	if !reflect.DeepEqual(sort2, s2) {
		t.Errorf("SortDesc should result in %v, got %v", s2, sort2)
	}

	if !sort2.IsSortedDesc() {
		t.Errorf("IsSortedDesc should be true for %v", sort1)
	}
}
