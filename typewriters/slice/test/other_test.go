package main

import (
	"testing"
)

// methods where underlying type is ordered
// +test slice:"Max,Min, Average, Sum, Sort,IsSorted,SortDesc,IsSortedDesc"
type Other Underlying

type Underlying int

var others = OtherSlice{50, 100, 9, 7, 100, 99}

func TestAverage(t *testing.T) {
	average1, err := others.Average()

	if err != nil {
		t.Errorf("Average should succeed")
	}

	avg1 := Other(60)

	if average1 != avg1 {
		t.Errorf("Average should be %v, got %v", avg1, average1)
	}

	average2, err := OtherSlice{}.Average()

	if err == nil || average2 != 0 {
		t.Errorf("Average should fail on empty slice")
	}
}

func TestMax(t *testing.T) {
	max1, err := others.Max()
	m1 := Other(100)

	if err != nil {
		t.Errorf("Max should succeed")
	}

	if max1 != m1 {
		t.Errorf("Max should return %v, got %v", m1, max1)
	}

	max2, err := OtherSlice{}.Max()
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

	min2, err := OtherSlice{}.Min()
	var m2 Other

	if err == nil || min2 != m2 {
		t.Errorf("Min should fail on empty slice")
	}
}

func TestSort(t *testing.T) {
	sort1 := others.Sort()
	s1 := OtherSlice{7, 9, 50, 99, 100, 100}

	if !otherSliceEqual(sort1, s1) {
		t.Errorf("Sort should result in %v, got %v", s1, sort1)
	}

	if !sort1.IsSorted() {
		t.Errorf("IsSorted should be true for %v", sort1)
	}

	sort2 := others.SortDesc()
	s2 := OtherSlice{100, 100, 99, 50, 9, 7}

	if !otherSliceEqual(sort2, s2) {
		t.Errorf("SortDesc should result in %v, got %v", s2, sort2)
	}

	if !sort2.IsSortedDesc() {
		t.Errorf("IsSortedDesc should be true for %v", sort1)
	}
}

func TestSum(t *testing.T) {
	sum := others.Sum()

	if sum != 365 {
		t.Errorf("Sum should result in %v, got %v", 365, sum)
	}
}
