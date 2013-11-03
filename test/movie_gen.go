// gen *models.Movie
// this file was auto-generated using github.com/clipperhouse/gen
// Sun, 03 Nov 2013 17:31:24 UTC

package models

// The plural (slice) type of *Movie, for use with gen methods below. Use this type where you would use []*Movie. (This is required because slices cannot be method receivers.)
type Movies []*Movie

// Iterates over Movies, operating on each element while accumulating results. For example, Sum & Min might be implemented like:
//	sum := func(_item *Movie, accumulated int) int {
//		return accumulated + _item.Something
//	}
//	sumOfSomething := myMovies.AggregateInt(sum)
//
//	min := func(_item *Movie, accumulated int) int {
//		if _item.AnotherThing < accumulated {
//			return _item.AnotherThing
//		}
//		return accumulated
//	}
//	minOfAnotherThing := myMovies.AggregateInt(min)
func (rcv Movies) AggregateInt(fn func(*Movie, int) int) int {
	result := 0
	for _, _item := range rcv {
		result = fn(_item, result)
	}
	return result
}

// Iterates over Movies, operating on each element while accumulating results. For example, you might join strings like:
//	myMovies := GetSomeMovies()
//	join := func(_item *Movie, accumulated string) string {
//		if _item != myMovies[0] {
//			accumulated += ", "
//		}
//		return accumulated + _item.Title
//	}
//	myList := myMovies.AggregateString(join)
func (rcv Movies) AggregateString(fn func(*Movie, string) string) (result string) {
	for _, _item := range rcv {
		result = fn(_item, result)
	}
	return result
}

// Tests that all elements of Movies are true for the passed func. Example:
//	good := func(_item *Movie) bool {
//		return _item.Something > 42
//	}
//	allGood := myMovies.All(good)
func (rcv Movies) All(fn func(*Movie) bool) bool {
	for _, _item := range rcv {
		if !fn(_item) {
			return false
		}
	}
	return true
}

// Tests that one or more elements of Movies are true for the passed func. Example:
//	winner := func(_item *Movie) bool {
//		return _item.Placement == "winner"
//	}
//	weHaveAWinner := myMovies.Any(winner)
func (rcv Movies) Any(fn func(*Movie) bool) bool {
	for _, _item := range rcv {
		if fn(_item) {
			return true
		}
	}
	return false
}

// Counts the number elements of Movies that are true for the passed func. Example:
//	dracula := func(_item *Movie) bool {
//		return _item.IsDracula()
//	}
//	countDracula := myMovies.Count(dracula)
func (rcv Movies) Count(fn func(*Movie) bool) int {
	var count = func(_item *Movie, acc int) int {
		if fn(_item) {
			acc++
		}
		return acc
	}
	return rcv.AggregateInt(count)
}

// Iterates over Movies and executes the passed func against each element.
func (rcv Movies) Each(fn func(*Movie)) {
	for _, _item := range rcv {
		fn(_item)
	}
}

// Groups Movies into a map of Movies, keyed by the result of the passed func. Example:
//	year := func(_item *Movie) int {
//		return _item.CreationDate.Year()
//	}
//	yearlyReport := myMovies.GroupByInt(year)
func (rcv Movies) GroupByInt(fn func(*Movie) int) map[int]Movies {
	result := make(map[int]Movies)
	for _, _item := range rcv {
		result[fn(_item)] = append(result[fn(_item)], _item)
	}
	return result
}

// Groups Movies into a map of Movies, keyed by the result of the passed func. Example:
//	dept := func(_item *Movie) string {
//		return _item.DepartmentName
//	}
//	byDepartment := myMovies.GroupByString(dept)
func (rcv Movies) GroupByString(fn func(*Movie) string) map[string]Movies {
	result := make(map[string]Movies)
	for _, _item := range rcv {
		result[fn(_item)] = append(result[fn(_item)], _item)
	}
	return result
}

// Returns the element of Movies containing the maximum value, when compared to other elements using a passed func defining ‘less’. Example:
//	byArea := func(_items Movies, a int, b int) bool {
//		return _items[a].Area() < _items[b].Area()
//	}
//	roomiest := myMovies.Max(byArea)
//
// In the case of multiple items being equally maximal, the last such element is returned.
// (Note: this is implemented by negating the passed ‘less’ func, effectively testing ‘greater than or equal to’.)
func (rcv Movies) Max(less func(Movies, int, int) bool) *Movie {
	return rcv.Min(not(less))
}

// Returns the element of Movies containing the minimum value, when compared to other elements using a passed func defining ‘less’. Example:
//	byPrice := func(_items Movies, a int, b int) bool {
//		return _items[a].Price < _items[b].Price
//	}
//	cheapest := myMovies.Min(byPrice)
//
// In the case of multiple items being equally minimal, the first such element is returned.
func (rcv Movies) Min(less func(Movies, int, int) bool) *Movie {
	var _nil *Movie
	l := len(rcv)
	if l == 0 {
		return _nil
	}
	m := 0
	for i := 1; i < l; i++ {
		if less(rcv, i, m) {
			m = i
		}
	}
	return rcv[m]
}

// Returns the sum of ints returned by passed func. Example:
//	itemTotal := func(_item *Movie) int {
//		return _item.Quantity * _item.UnitPrice
//	}
//	orderTotal := myMovies.SumInt(itemTotal)
func (rcv Movies) SumInt(fn func(*Movie) int) int {
	var sum = func(_item *Movie, acc int) int {
		return acc + fn(_item)
	}
	return rcv.AggregateInt(sum)
}

// Returns a new Movies slice whose elements return true for func. Example:
//	incredible := func(_item *Movie) bool {
//		return _item.Manufacturer == "Apple"
//	}
//	wishList := myMovies.Where(incredible)
func (rcv Movies) Where(fn func(*Movie) bool) (result Movies) {
	for _, _item := range rcv {
		if fn(_item) {
			result = append(result, _item)
		}
	}
	return result
}

// Sort functions below are a modification of http://golang.org/pkg/sort/#Sort
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Returns a new ordered Movies slice, determined by a func defining ‘less’. Example:
//	byName := func(_items Movies, a int, b int) bool {
//		return _items[a].LastName < _items[b].LastName
//	}
//	roster := myMovies.Sort(byName)
func (rcv Movies) Sort(less func(Movies, int, int) bool) Movies {
	result := make(Movies, len(rcv))
	copy(result, rcv)

	// Switch to heapsort if depth of 2*ceil(lg(n+1)) is reached.
	n := len(result)
	maxDepth := 0
	for i := n; i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSortMovies(result, less, 0, n, maxDepth)
	return result
}

// Reports whether an instance of Movies is sorted, using the pass func to define ‘less’. See Sort method below.
func (rcv Movies) IsSorted(less func(Movies, int, int) bool) bool {
	n := len(rcv)
	for i := n - 1; i > 0; i-- {
		if less(rcv, i, i-1) {
			return false
		}
	}
	return true
}

// Returns a new, descending-ordered Movies slice, determined by a func defining ‘less’. Example:
//	byPoints := func(_items Movies, a int, b int) bool {
//		return _items[a].Points < _items[b].Points
//	}
//	leaderboard := myMovies.SortDesc(byPoints)
// (Note: this is implemented by negating the passed ‘less’ func, effectively testing ‘greater than or equal to’.)
func (rcv Movies) SortDesc(less func(Movies, int, int) bool) Movies {
	return rcv.Sort(not(less))
}

// Reports whether an instance of Movies is sorted in descending order, using the pass func to define ‘less’. See SortDesc method below.
func (rcv Movies) IsSortedDesc(less func(Movies, int, int) bool) bool {
	return rcv.IsSorted(not(less))
}

func swapMovies(rcv Movies, a, b int) {
	rcv[a], rcv[b] = rcv[b], rcv[a]
}

// Insertion sort
func insertionSortMovies(rcv Movies, less func(Movies, int, int) bool, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && less(rcv, j, j-1); j-- {
			swapMovies(rcv, j, j-1)
		}
	}
}

// siftDown implements the heap property on rcv[lo, hi).
// first is an offset into the array where the root of the heap lies.
func siftDownMovies(rcv Movies, less func(Movies, int, int) bool, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && less(rcv, first+child, first+child+1) {
			child++
		}
		if !less(rcv, first+root, first+child) {
			return
		}
		swapMovies(rcv, first+root, first+child)
		root = child
	}
}

func heapSortMovies(rcv Movies, less func(Movies, int, int) bool, a, b int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDownMovies(rcv, less, i, hi, first)
	}

	// Pop elements, largest first, into end of rcv.
	for i := hi - 1; i >= 0; i-- {
		swapMovies(rcv, first, first+i)
		siftDownMovies(rcv, less, lo, i, first)
	}
}

// Quicksort, following Bentley and McIlroy,
// Engineering a Sort Function, SP&E November 1993.

// medianOfThree moves the median of the three values rcv[a], rcv[b], rcv[c] into rcv[a].
func medianOfThreeMovies(rcv Movies, less func(Movies, int, int) bool, a, b, c int) {
	m0 := b
	m1 := a
	m2 := c
	// bubble sort on 3 elements
	if less(rcv, m1, m0) {
		swapMovies(rcv, m1, m0)
	}
	if less(rcv, m2, m1) {
		swapMovies(rcv, m2, m1)
	}
	if less(rcv, m1, m0) {
		swapMovies(rcv, m1, m0)
	}
	// now rcv[m0] <= rcv[m1] <= rcv[m2]
}

func swapRangeMovies(rcv Movies, a, b, n int) {
	for i := 0; i < n; i++ {
		swapMovies(rcv, a+i, b+i)
	}
}

func doPivotMovies(rcv Movies, less func(Movies, int, int) bool, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2 // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		// Tukey's Ninther, median of three medians of three.
		s := (hi - lo) / 8
		medianOfThreeMovies(rcv, less, lo, lo+s, lo+2*s)
		medianOfThreeMovies(rcv, less, m, m-s, m+s)
		medianOfThreeMovies(rcv, less, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThreeMovies(rcv, less, lo, m, hi-1)

	// Invariants are:
	//	rcv[lo] = pivot (set up by ChoosePivot)
	//	rcv[lo <= i < a] = pivot
	//	rcv[a <= i < b] < pivot
	//	rcv[b <= i < c] is unexamined
	//	rcv[c <= i < d] > pivot
	//	rcv[d <= i < hi] = pivot
	//
	// Once b meets c, can swap the "= pivot" sections
	// into the middle of the slice.
	pivot := lo
	a, b, c, d := lo+1, lo+1, hi, hi
	for {
		for b < c {
			if less(rcv, b, pivot) { // rcv[b] < pivot
				b++
			} else if !less(rcv, pivot, b) { // rcv[b] = pivot
				swapMovies(rcv, a, b)
				a++
				b++
			} else {
				break
			}
		}
		for b < c {
			if less(rcv, pivot, c-1) { // rcv[c-1] > pivot
				c--
			} else if !less(rcv, c-1, pivot) { // rcv[c-1] = pivot
				swapMovies(rcv, c-1, d-1)
				c--
				d--
			} else {
				break
			}
		}
		if b >= c {
			break
		}
		// rcv[b] > pivot; rcv[c-1] < pivot
		swapMovies(rcv, b, c-1)
		b++
		c--
	}

	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	n := min(b-a, a-lo)
	swapRangeMovies(rcv, lo, b-n, n)

	n = min(hi-d, d-c)
	swapRangeMovies(rcv, c, hi-n, n)

	return lo + b - a, hi - (d - c)
}

func quickSortMovies(rcv Movies, less func(Movies, int, int) bool, a, b, maxDepth int) {
	for b-a > 7 {
		if maxDepth == 0 {
			heapSortMovies(rcv, less, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivotMovies(rcv, less, a, b)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			quickSortMovies(rcv, less, a, mlo, maxDepth)
			a = mhi // i.e., quickSortMovies(rcv, mhi, b)
		} else {
			quickSortMovies(rcv, less, mhi, b, maxDepth)
			b = mlo // i.e., quickSortMovies(rcv, a, mlo)
		}
	}
	if b-a > 1 {
		insertionSortMovies(rcv, less, a, b)
	}
}

func not(less func(Movies, int, int) bool) func(Movies, int, int) bool {
	return func(z Movies, a int, b int) bool {
		return !less(z, a, b)
	}
}
