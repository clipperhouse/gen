package main

import (
	"testing"
)

func TestAll(t *testing.T) {
	all1 := things.All(func(x Thing) bool {
		return x.Name == "First"
	})

	if all1 {
		t.Errorf("All should not evaulate true for Name == First")
	}

	all2 := things.All(func(x Thing) bool {
		return x.Number > 1
	})

	if !all2 {
		t.Errorf("All should evaulate true for Number > 1")
	}

	all3 := noThings.All(func(x Thing) bool {
		return false
	})

	if !all3 {
		t.Errorf("All should evaulate true for empty slices")
	}
}

func TestAny(t *testing.T) {
	any1 := things.Any(func(x Thing) bool {
		return x.Name == "Dummy"
	})

	if any1 {
		t.Errorf("Any should not evaulate true for Name == Dummy")
	}

	any2 := things.Any(func(x Thing) bool {
		return x.Number > 50
	})

	if !any2 {
		t.Errorf("Any should evaulate true for Number > 50")
	}

	any3 := noThings.Any(func(x Thing) bool {
		return true
	})

	if any3 {
		t.Errorf("Any should evaulate false for empty slices")
	}
}

func TestCount(t *testing.T) {
	count1 := things.Count(func(x Thing) bool {
		return x.Name == "Second"
	})

	if count1 != 1 {
		t.Errorf("Count should find one item Name == Second")
	}

	count2 := things.Count(func(x Thing) bool {
		return x.Number > 50
	})

	if count2 != 3 {
		t.Errorf("Count should find 3 items for Number > 50")
	}

	count3 := things.Count(func(x Thing) bool {
		return x.Name == "Dummy"
	})

	if count3 != 0 {
		t.Errorf("Count should no items for Name == Dummy")
	}

	count4 := noThings.Count(func(x Thing) bool {
		return true
	})

	if count4 != 0 {
		t.Errorf("Count should find no items in an empty slice")
	}
}

func TestDistinct(t *testing.T) {
	distinct1 := things.Distinct()

	if !sliceEqual(distinct1, Things{first, second, third, fourth}) {
		t.Errorf("Distinct should exclude %v, but found %v", anotherThird, distinct1)
	}
}

func TestDistinctBy(t *testing.T) {
	distinctby1 := things.DistinctBy(func(a, b Thing) bool {
		return a.Number == b.Number
	})

	if !sliceEqual(distinctby1, Things{first, second, third}) {
		t.Errorf("DistinctBy should exclude %v and %v, but found %v", anotherThird, fourth, things)
	}
}

func TestFirst(t *testing.T) {
	f1, err := things.First(func(x Thing) bool {
		return x.Name == "Third"
	})

	if err != nil {
		t.Errorf("First should succeed when finding Name == Third")
	}

	if f1 != third {
		t.Errorf("First should find %v, but found %v", third, f1)
	}

	f2, err := things.First(func(x Thing) bool {
		return x.Name == "Dummy"
	})

	if err == nil || f2 != zero {
		t.Errorf("First should fail when finding Name == Dummy")
	}

	f3, err := noThings.First(func(x Thing) bool {
		return true
	})

	if err == nil || f3 != zero {
		t.Errorf("First should fail on empty slice")
	}
}

func TestMinBy(t *testing.T) {
	min1, err := things.MinBy(func(a, b Thing) bool {
		return a.Number < b.Number
	})

	if err != nil {
		t.Errorf("MinBy Number should succeed")
	}

	if min1 != second {
		t.Errorf("MinBy Number should return %v, got %v", second, min1)
	}

	min2, err := noThings.MinBy(func(a, b Thing) bool {
		return true
	})

	if err == nil || min2 != zero {
		t.Errorf("MinBy Number should fail on empty slice")
	}
}

func TestMaxBy(t *testing.T) {
	max1, err := things.MaxBy(func(a, b Thing) bool {
		return a.Number < b.Number
	})

	if err != nil {
		t.Errorf("MaxBy Number should succeed")
	}

	if max1 != third {
		t.Errorf("MaxBy Number should return %v, got %v", third, max1)
	}

	max2, err := noThings.MaxBy(func(a, b Thing) bool {
		return true
	})

	if err == nil || max2 != zero {
		t.Errorf("MaxBy Number should fail on empty slice")
	}
}

func TestSingle(t *testing.T) {
	single1, err := things.Single(func(a Thing) bool {
		return a.Name == "Second"
	})

	if err != nil {
		t.Errorf("Single Name should succeed")
	}

	if single1 != second {
		t.Errorf("Single should return %v, got %v", second, single1)
	}

	single2, err := things.Single(func(a Thing) bool {
		return a.Name == "Third"
	})

	if err == nil || single2 != zero {
		t.Errorf("Single should error on Name == Third")
	}

	single3, err := noThings.Single(func(a Thing) bool {
		return true
	})

	if err == nil || single3 != zero {
		t.Errorf("Single should fail on empty slice")
	}
}

func TestSortBy(t *testing.T) {
	name := func(a, b Thing) bool {
		return a.Name < b.Name
	}

	sort1 := things.SortBy(name)

	sorted1 := Things{first, fourth, second, third, anotherThird}

	if !sliceEqual(sort1, Things{first, fourth, second, third, anotherThird}) {
		t.Errorf("SortBy name should be %v, got %v", sorted1, sort1)
	}

	if !sort1.IsSortedBy(name) {
		t.Errorf("IsSortedBy name should be true")
	}

	if things.IsSortedBy(name) {
		t.Errorf("things should not be sorted by name")
	}

	sort2 := things.SortByDesc(name)

	sorted2 := Things{anotherThird, third, second, fourth, first}

	if !sliceEqual(sort2, sorted2) {
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

	sorted3 := Things{eighth, fifth, first, fourth, second, seventh, sixth, third}

	if !sliceEqual(sort3, sorted3) {
		t.Errorf("Sort name should be %v, got %v", sorted3, sort3)
	}

	// intended to hit threshold to invoke medianOfThree (40)
	var evenMore Things
	evenMore = append(evenMore, lotsOfThings...)
	evenMore = append(evenMore, lotsOfThings...)
	evenMore = append(evenMore, lotsOfThings...)
	evenMore = append(evenMore, lotsOfThings...)
	evenMore = append(evenMore, lotsOfThings...)
	evenMore = append(evenMore, lotsOfThings...)

	sort4 := evenMore.SortBy(name)

	sorted4 := Things{eighth, eighth, eighth, eighth, eighth, eighth}
	sorted4 = append(sorted4, appendMany(fifth, 6)...)
	sorted4 = append(sorted4, appendMany(first, 6)...)
	sorted4 = append(sorted4, appendMany(fourth, 6)...)
	sorted4 = append(sorted4, appendMany(second, 6)...)
	sorted4 = append(sorted4, appendMany(seventh, 6)...)
	sorted4 = append(sorted4, appendMany(sixth, 6)...)
	sorted4 = append(sorted4, appendMany(third, 6)...)

	if !sliceEqual(sort4, sorted4) {
		t.Errorf("Sort name should be %v, got %v", sorted3, sort3)
	}
}

func TestWhere(t *testing.T) {
	where1 := things.Where(func(x Thing) bool {
		return x.Name == "Third"
	})

	w1 := Things{third, anotherThird}

	if !sliceEqual(where1, w1) {
		t.Errorf("Where should result in %v, got %v", w1, where1)
	}

	where2 := things.Where(func(x Thing) bool {
		return x.Name == "Dummy"
	})

	if len(where2) != 0 {
		t.Errorf("Where should result in empty slice, got %v", where2)
	}

	where3 := noThings.Where(func(x Thing) bool {
		return true
	})

	if len(where3) != 0 {
		t.Errorf("Where should result in empty slice, got %v", where3)
	}
}

func TestAggregateOther(t *testing.T) {
	sum := func(state Other, x Thing) Other {
		return state + x.Number
	}

	aggregate1 := things.AggregateOther(sum)
	agg1 := Other(340)

	if aggregate1 != agg1 {
		t.Errorf("AggregateOther should be %v, got %v", agg1, aggregate1)
	}
}

func TestAverageOther(t *testing.T) {
	number := func(x Thing) Other {
		return x.Number
	}

	average1, err := things.AverageOther(number)

	if err != nil {
		t.Errorf("Average should succeed")
	}

	avg1 := Other(68)

	if average1 != avg1 {
		t.Errorf("Average should be %v, got %v", avg1, average1)
	}

	average2, err := noThings.AverageOther(number)

	if err == nil || average2 != 0 {
		t.Errorf("Average should fail on empty slice")
	}
}

func TestGroupByOther(t *testing.T) {
	number := func(x Thing) Other {
		return x.Number
	}

	groupby1 := things.GroupByOther(number)
	g1 := map[Other]Things{
		40:  {second, fourth},
		60:  {first},
		100: {third, anotherThird},
	}

	if len(groupby1) != len(g1) {
		t.Errorf("GroupByInt result should have %d elements, has %d", len(g1), len(groupby1))
	}

	for k, v := range g1 {
		g, ok := groupby1[k]

		if !ok {
			t.Errorf("GroupByOther result should have %d element, but is %v", k, len(groupby1))
		}

		if !sliceEqual(v, g) {
			t.Errorf("GroupByOther result [%d] should have %v but has %v", k, v, g)
		}
	}
}

func TestMaxInt(t *testing.T) {
	number := func(x Thing) Other {
		return x.Number
	}

	max1, err := things.MaxOther(number)

	if err != nil {
		t.Errorf("Max should succeed")
	}

	if max1 != 100 {
		t.Errorf("Max should be %v, got %v", 100, max1)
	}

	max2, err := noThings.MaxOther(number)

	if err == nil || max2 != 0 {
		t.Errorf("Max should fail on empty slice")
	}
}

func TestMinOther(t *testing.T) {
	number := func(x Thing) Other {
		return x.Number
	}

	min1, err := things.MinOther(number)

	if err != nil {
		t.Errorf("Min should succeed")
	}

	if min1 != 40 {
		t.Errorf("Min should be %v, got %v", 40, min1)
	}

	min2, err := noThings.MinOther(number)

	if err == nil || min2 != 0 {
		t.Errorf("Min should fail on empty slice")
	}
}

func TestSelectOther(t *testing.T) {
	number := func(x Thing) Other {
		return x.Number
	}

	select1 := things.SelectOther(number)
	s1 := []Other{60, 40, 100, 100, 40}

	if !otherSliceEqual(select1, s1) {
		t.Errorf("Select should result in %v, got %v", s1, select1)
	}
}

func TestSum(t *testing.T) {
	number := func(x Thing) Other {
		return x.Number
	}

	sum1 := things.SumOther(number)

	if sum1 != 340 {
		t.Errorf("Sum should result in %v, got %v", 340, sum1)
	}
}

func appendMany(x Thing, n int) (result Things) {
	for i := 0; i < n; i++ {
		result = append(result, x)
	}
	return
}

func sliceEqual(a, b Things) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func intSliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func otherSliceEqual(a, b Others) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
