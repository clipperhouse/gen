package models

import (
	"testing"
)

var (
	zero1, first1, second1, third1, anotherThird1, fourth1 Thing1
	fifth1, sixth1, seventh1, eighth1                      Thing1
	thing1s, no1s, lotsOfThing1s                           Thing1s
	thing2s                                                Thing2s
)

func init() {
	zero1 = Thing1{}
	first1 = Thing1{"First", 60}
	second1 = Thing1{"Second", 40}
	third1 = Thing1{"Third", 100}
	anotherThird1 = Thing1{"Third", 100}
	fourth1 = Thing1{"Fourth", 40}
	fifth1 = Thing1{"Fifth", 70}
	sixth1 = Thing1{"Sixth", 10}
	seventh1 = Thing1{"Seventh", 50}
	eighth1 = Thing1{"Eighth", 110}

	thing1s = Thing1s{
		first1,
		second1,
		third1,
		anotherThird1,
		fourth1,
	}
	no1s = Thing1s{}
	lotsOfThing1s = Thing1s{
		first1,
		second1,
		third1,
		fourth1,
		fifth1,
		sixth1,
		seventh1,
		eighth1,
	}
	thing2s = Thing2s{50.1, 7, 99.8, 100.4, 80}
}

func TestAll(t *testing.T) {
	all1 := thing1s.All(func(x Thing1) bool {
		return x.Name == "First"
	})

	if all1 {
		t.Errorf("All should not evaulate true for Name == First")
	}

	all2 := thing1s.All(func(x Thing1) bool {
		return x.Number > 1
	})

	if !all2 {
		t.Errorf("All should evaulate true for Number > 1")
	}

	all3 := no1s.All(func(x Thing1) bool {
		return false
	})

	if !all3 {
		t.Errorf("All should evaulate true for empty slices")
	}
}

func TestAny(t *testing.T) {
	any1 := thing1s.Any(func(x Thing1) bool {
		return x.Name == "Dummy"
	})

	if any1 {
		t.Errorf("Any should not evaulate true for Name == Dummy")
	}

	any2 := thing1s.Any(func(x Thing1) bool {
		return x.Number > 50
	})

	if !any2 {
		t.Errorf("Any should evaulate true for Number > 50")
	}

	any3 := no1s.Any(func(x Thing1) bool {
		return true
	})

	if any3 {
		t.Errorf("Any should evaulate false for empty slices")
	}
}

func TestCount(t *testing.T) {
	count1 := thing1s.Count(func(x Thing1) bool {
		return x.Name == "Second"
	})

	if count1 != 1 {
		t.Errorf("Count should find one item Name == Second")
	}

	count2 := thing1s.Count(func(x Thing1) bool {
		return x.Number > 50
	})

	if count2 != 3 {
		t.Errorf("Count should find 3 items for Number > 50")
	}

	count3 := thing1s.Count(func(x Thing1) bool {
		return x.Name == "Dummy"
	})

	if count3 != 0 {
		t.Errorf("Count should no items for Name == Dummy")
	}

	count4 := no1s.Count(func(x Thing1) bool {
		return true
	})

	if count4 != 0 {
		t.Errorf("Count should find no items in an empty slice")
	}
}

func TestDistinct(t *testing.T) {
	distinct1 := thing1s.Distinct()

	if !sliceEqual(distinct1, Thing1s{first1, second1, third1, fourth1}) {
		t.Errorf("Distinct should exclude %v, but found %v", anotherThird1, distinct1)
	}
}

func TestDistinctBy(t *testing.T) {
	distinctby1 := thing1s.DistinctBy(func(a, b Thing1) bool {
		return a.Number == b.Number
	})

	if !sliceEqual(distinctby1, Thing1s{first1, second1, third1}) {
		t.Errorf("DistinctBy should exclude %v and %v, but found %v", anotherThird1, fourth1, thing1s)
	}
}

func TestFirst(t *testing.T) {
	f1, err := thing1s.First(func(x Thing1) bool {
		return x.Name == "Third"
	})

	if err != nil {
		t.Errorf("First should succeed when finding Name == Third")
	}

	if f1 != third1 {
		t.Errorf("First should find %v, but found %v", third1, f1)
	}

	f2, err := thing1s.First(func(x Thing1) bool {
		return x.Name == "Dummy"
	})

	if err == nil || f2 != zero1 {
		t.Errorf("First should fail when finding Name == Dummy")
	}

	f3, err := no1s.First(func(x Thing1) bool {
		return true
	})

	if err == nil || f3 != zero1 {
		t.Errorf("First should fail on empty slice")
	}
}

func TestMin(t *testing.T) {
	min1, err := thing2s.Min()
	m1 := Thing2(7)

	if err != nil {
		t.Errorf("Min should succeed")
	}

	if min1 != m1 {
		t.Errorf("Min should return %v, got %v", m1, min1)
	}

	min2, err := Thing2s{}.Min()
	var m2 Thing2

	if err == nil || min2 != m2 {
		t.Errorf("Min should fail on empty slice")
	}
}

func TestMinBy(t *testing.T) {
	min1, err := thing1s.MinBy(func(a, b Thing1) bool {
		return a.Number < b.Number
	})

	if err != nil {
		t.Errorf("MinBy Number should succeed")
	}

	if min1 != second1 {
		t.Errorf("MinBy Number should return %v, got %v", second1, min1)
	}

	min2, err := no1s.MinBy(func(a, b Thing1) bool {
		return true
	})

	if err == nil || min2 != zero1 {
		t.Errorf("MinBy Number should fail on empty slice")
	}
}

func TestMax(t *testing.T) {
	max1, err := thing2s.Max()
	m1 := Thing2(100.4)

	if err != nil {
		t.Errorf("Max should succeed")
	}

	if max1 != m1 {
		t.Errorf("Max should return %v, got %v", m1, max1)
	}

	max2, err := Thing2s{}.Max()
	var m2 Thing2

	if err == nil || max2 != m2 {
		t.Errorf("Max should fail on empty slice")
	}
}

func TestMaxBy(t *testing.T) {
	max1, err := thing1s.MaxBy(func(a, b Thing1) bool {
		return a.Number < b.Number
	})

	if err != nil {
		t.Errorf("MaxBy Number should succeed")
	}

	if max1 != third1 {
		t.Errorf("MaxBy Number should return %v, got %v", third1, max1)
	}

	max2, err := no1s.MaxBy(func(a, b Thing1) bool {
		return true
	})

	if err == nil || max2 != zero1 {
		t.Errorf("MaxBy Number should fail on empty slice")
	}
}

func TestSingle(t *testing.T) {
	single1, err := thing1s.Single(func(a Thing1) bool {
		return a.Name == "Second"
	})

	if err != nil {
		t.Errorf("Single Name should succeed")
	}

	if single1 != second1 {
		t.Errorf("Single should return %v, got %v", second1, single1)
	}

	single2, err := thing1s.Single(func(a Thing1) bool {
		return a.Name == "Third"
	})

	if err == nil || single2 != zero1 {
		t.Errorf("Single should error on Name == Third")
	}

	single3, err := no1s.Single(func(a Thing1) bool {
		return true
	})

	if err == nil || single3 != zero1 {
		t.Errorf("Single should fail on empty slice")
	}
}

func TestSort(t *testing.T) {
	sort1 := thing2s.Sort()
	s1 := Thing2s{7, 50.1, 80, 99.8, 100.4}

	if !thing2SliceEqual(sort1, s1) {
		t.Errorf("Sort should result in %v, got %v", s1, sort1)
	}

	if !sort1.IsSorted() {
		t.Errorf("IsSorted should be true for %v", sort1)
	}

	sort2 := thing2s.SortDesc()
	s2 := Thing2s{100.4, 99.8, 80, 50.1, 7}

	if !thing2SliceEqual(sort2, s2) {
		t.Errorf("SortDesc should result in %v, got %v", s2, sort2)
	}

	if !sort2.IsSortedDesc() {
		t.Errorf("IsSortedDesc should be true for %v", sort1)
	}
}

func TestSortBy(t *testing.T) {
	name := func(a, b Thing1) bool {
		return a.Name < b.Name
	}

	sort1 := thing1s.SortBy(name)

	sorted1 := Thing1s{first1, fourth1, second1, third1, anotherThird1}

	if !sliceEqual(sort1, Thing1s{first1, fourth1, second1, third1, anotherThird1}) {
		t.Errorf("SortBy name should be %v, got %v", sorted1, sort1)
	}

	if !sort1.IsSortedBy(name) {
		t.Errorf("IsSortedBy name should be true")
	}

	if thing1s.IsSortedBy(name) {
		t.Errorf("thing1s should not be sorted by name")
	}

	sort2 := thing1s.SortByDesc(name)

	sorted2 := Thing1s{anotherThird1, third1, second1, fourth1, first1}

	if !sliceEqual(sort2, sorted2) {
		t.Errorf("SortByDesc name should be %v, got %v", sorted2, sort2)
	}

	if !sort2.IsSortedByDesc(name) {
		t.Errorf("IsSortedByDesc name should be true %v", sort2)
	}

	if thing1s.IsSortedByDesc(name) {
		t.Errorf("thing1s should not be sorted desc by name")
	}

	// intended to hit threshold to invoke quicksort (7)
	sort3 := lotsOfThing1s.SortBy(name)

	sorted3 := Thing1s{eighth1, fifth1, first1, fourth1, second1, seventh1, sixth1, third1}

	if !sliceEqual(sort3, sorted3) {
		t.Errorf("Sort name should be %v, got %v", sorted3, sort3)
	}

	// intended to hit threshold to invoke medianOfThree (40)
	var evenMore Thing1s
	evenMore = append(evenMore, lotsOfThing1s...)
	evenMore = append(evenMore, lotsOfThing1s...)
	evenMore = append(evenMore, lotsOfThing1s...)
	evenMore = append(evenMore, lotsOfThing1s...)
	evenMore = append(evenMore, lotsOfThing1s...)
	evenMore = append(evenMore, lotsOfThing1s...)

	sort4 := evenMore.SortBy(name)

	sorted4 := Thing1s{eighth1, eighth1, eighth1, eighth1, eighth1, eighth1}
	sorted4 = append(sorted4, appendMany(fifth1, 6)...)
	sorted4 = append(sorted4, appendMany(first1, 6)...)
	sorted4 = append(sorted4, appendMany(fourth1, 6)...)
	sorted4 = append(sorted4, appendMany(second1, 6)...)
	sorted4 = append(sorted4, appendMany(seventh1, 6)...)
	sorted4 = append(sorted4, appendMany(sixth1, 6)...)
	sorted4 = append(sorted4, appendMany(third1, 6)...)

	if !sliceEqual(sort4, sorted4) {
		t.Errorf("Sort name should be %v, got %v", sorted3, sort3)
	}
}

func TestWhere(t *testing.T) {
	where1 := thing1s.Where(func(x Thing1) bool {
		return x.Name == "Third"
	})

	w1 := Thing1s{third1, anotherThird1}

	if !sliceEqual(where1, w1) {
		t.Errorf("Where should result in %v, got %v", w1, where1)
	}

	where2 := thing1s.Where(func(x Thing1) bool {
		return x.Name == "Dummy"
	})

	if len(where2) != 0 {
		t.Errorf("Where should result in empty slice, got %v", where2)
	}

	where3 := no1s.Where(func(x Thing1) bool {
		return true
	})

	if len(where3) != 0 {
		t.Errorf("Where should result in empty slice, got %v", where3)
	}
}

func TestAggregate(t *testing.T) {
	join := func(state string, x Thing1) string {
		if len(state) > 0 {
			state += ", "
		}
		return state + x.Name
	}

	aggregate1 := thing1s.AggregateString(join)
	agg1 := "First, Second, Third, Third, Fourth"

	if aggregate1 != agg1 {
		t.Errorf("AggregateString should be %v, got %v", agg1, aggregate1)
	}
}

func TestAverage(t *testing.T) {
	number := func(x Thing1) int {
		return x.Number
	}

	average1, err := thing1s.AverageInt(number)

	if err != nil {
		t.Errorf("Average should succeed")
	}

	avg1 := 68

	if average1 != avg1 {
		t.Errorf("Average should be %v, got %v", avg1, average1)
	}

	average2, err := no1s.AverageInt(number)

	if err == nil || average2 != 0 {
		t.Errorf("Average should fail on empty slice")
	}
}

func TestGroupBy(t *testing.T) {
	number := func(x Thing1) int {
		return x.Number
	}

	groupby1 := thing1s.GroupByInt(number)
	g1 := map[int]Thing1s{
		40:  {second1, fourth1},
		60:  {first1},
		100: {third1, anotherThird1},
	}

	if len(groupby1) != len(g1) {
		t.Errorf("GroupByInt result should have %d elements, has %d", len(g1), len(groupby1))
	}

	for k, v := range g1 {
		g, ok := groupby1[k]

		if !ok {
			t.Errorf("GroupByInt result should have %d element, but is %v", k, len(groupby1))
		}

		if !sliceEqual(v, g) {
			t.Errorf("GroupByInt result [%d] should have %v but has %v", k, v, g)
		}
	}
}

func TestMaxInt(t *testing.T) {
	number := func(x Thing1) int {
		return x.Number
	}

	max1, err := thing1s.MaxInt(number)

	if err != nil {
		t.Errorf("Max should succeed")
	}

	if max1 != 100 {
		t.Errorf("Max should be %v, got %v", 100, max1)
	}

	max2, err := no1s.MaxInt(number)

	if err == nil || max2 != 0 {
		t.Errorf("Max should fail on empty slice")
	}
}

func TestMinInt(t *testing.T) {
	number := func(x Thing1) int {
		return x.Number
	}

	min1, err := thing1s.MinInt(number)

	if err != nil {
		t.Errorf("Min should succeed")
	}

	if min1 != 40 {
		t.Errorf("Min should be %v, got %v", 40, min1)
	}

	min2, err := no1s.MinInt(number)

	if err == nil || min2 != 0 {
		t.Errorf("Min should fail on empty slice")
	}
}

func TestSelect(t *testing.T) {
	number := func(x Thing1) int {
		return x.Number
	}

	select1 := thing1s.SelectInt(number)
	s1 := []int{60, 40, 100, 100, 40}

	if !intSliceEqual(select1, s1) {
		t.Errorf("Select should result in %v, got %v", s1, select1)
	}
}

func TestSum(t *testing.T) {
	number := func(x Thing1) int {
		return x.Number
	}

	sum1 := thing1s.SumInt(number)

	if sum1 != 340 {
		t.Errorf("Sum should result in %v, got %v", 340, sum1)
	}
}

func appendMany(x Thing1, n int) (result Thing1s) {
	for i := 0; i < n; i++ {
		result = append(result, x)
	}
	return
}

func sliceEqual(a, b Thing1s) bool {
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

func thing2SliceEqual(a, b Thing2s) bool {
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
