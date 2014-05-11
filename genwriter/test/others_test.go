package main

import (
	"testing"
)

func TestMax(t *testing.T) {
	max1, err := others.Max()
	m1 := Other(100)

	if err != nil {
		t.Errorf("Max should succeed")
	}

	if max1 != m1 {
		t.Errorf("Max should return %v, got %v", m1, max1)
	}

	max2, err := Others{}.Max()
	var m2 Other

	if err == nil || max2 != m2 {
		t.Errorf("Max should fail on empty slice")
	}
}

func TestMin(t *testing.T) {
	min1, err := others.Min()
	m1 := Other(7)

	if err != nil {
		t.Errorf("Min should succeed")
	}

	if min1 != m1 {
		t.Errorf("Min should return %v, got %v", m1, min1)
	}

	min2, err := Others{}.Min()
	var m2 Other

	if err == nil || min2 != m2 {
		t.Errorf("Min should fail on empty slice")
	}
}

func TestSort(t *testing.T) {
	sort1 := others.Sort()
	s1 := Others{7, 9, 50, 99, 100, 100}

	if !otherSliceEqual(sort1, s1) {
		t.Errorf("Sort should result in %v, got %v", s1, sort1)
	}

	if !sort1.IsSorted() {
		t.Errorf("IsSorted should be true for %v", sort1)
	}

	sort2 := others.SortDesc()
	s2 := Others{100, 100, 99, 50, 9, 7}

	if !otherSliceEqual(sort2, s2) {
		t.Errorf("SortDesc should result in %v, got %v", s2, sort2)
	}

	if !sort2.IsSortedDesc() {
		t.Errorf("IsSortedDesc should be true for %v", sort1)
	}
}
