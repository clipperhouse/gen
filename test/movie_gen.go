// gen *models.Movie
// this file was auto-generated using github.com/clipperhouse/gen
// Sun, 05 Jan 2014 05:35:46 UTC

// Sort functions are a modification of http://golang.org/pkg/sort/#Sort
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package models

import (
	"errors"
)

// The plural (slice) type of *Movie, for use with gen methods below. Use this type where you would use []*Movie. (This is required because slices cannot be method receivers.)
type Movies []*Movie

// Tests that all elements of Movies return true for the passed func. See: http://clipperhouse.github.io/gen/#All
func (rcv Movies) All(fn func(*Movie) bool) bool {
	for _, v := range rcv {
		if !fn(v) {
			return false
		}
	}
	return true
}

// Tests that one or more elements of Movies return true for the passed func. See: http://clipperhouse.github.io/gen/#Any
func (rcv Movies) Any(fn func(*Movie) bool) bool {
	for _, v := range rcv {
		if fn(v) {
			return true
		}
	}
	return false
}

// Counts the number elements of Movies that return true for the passed func. See: http://clipperhouse.github.io/gen/#Count
func (rcv Movies) Count(fn func(*Movie) bool) (result int) {
	for _, v := range rcv {
		if fn(v) {
			result++
		}
	}
	return
}

// Returns a new Movies slice whose elements are unique. See: http://clipperhouse.github.io/gen/#Distinct
func (rcv Movies) Distinct() (result Movies) {
	appended := make(map[*Movie]bool)
	for _, v := range rcv {
		if !appended[v] {
			result = append(result, v)
			appended[v] = true
		}
	}
	return result
}

// Returns a new Movies slice whose elements are unique, where equality is defined by a passed func. See: http://clipperhouse.github.io/gen/#DistinctBy
func (rcv Movies) DistinctBy(equal func(*Movie, *Movie) bool) (result Movies) {
	for _, v := range rcv {
		eq := func(_app *Movie) bool {
			return equal(v, _app)
		}
		if !result.Any(eq) {
			result = append(result, v)
		}
	}
	return result
}

// Iterates over Movies and executes the passed func against each element. See: http://clipperhouse.github.io/gen/#Each
func (rcv Movies) Each(fn func(*Movie)) {
	for _, v := range rcv {
		fn(v)
	}
}

// Returns the first element that returns true for the passed func. Returns error if no elements return true. See: http://clipperhouse.github.io/gen/#First
func (rcv Movies) First(fn func(*Movie) bool) (result *Movie, err error) {
	for _, v := range rcv {
		if fn(v) {
			result = v
			return
		}
	}
	err = errors.New("No Movies elements return true for passed func")
	return
}

// Reports whether an instance of Movies is sorted, using the pass func to define ‘less’. See: http://clipperhouse.github.io/gen/#Sort
func (rcv Movies) IsSorted(less func(*Movie, *Movie) bool) bool {
	n := len(rcv)
	for i := n - 1; i > 0; i-- {
		if less(rcv[i], rcv[i-1]) {
			return false
		}
	}
	return true
}

// Reports whether an instance of Movies is sorted in descending order, using the pass func to define ‘less’. See: http://clipperhouse.github.io/gen/#Sort
func (rcv Movies) IsSortedDesc(less func(*Movie, *Movie) bool) bool {
	greaterOrEqual := func(a, b *Movie) bool {
		return !less(a, b)
	}
	return rcv.IsSorted(greaterOrEqual)
}

// Returns an element of Movies containing the maximum value, when compared to other elements using a passed func defining ‘less’. In the case of multiple items being equally maximal, the last such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#Max
//
// (Note: this is implemented by negating the passed ‘less’ func, effectively testing ‘greater than or equal to’.)
func (rcv Movies) Max(less func(*Movie, *Movie) bool) (result *Movie, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("Cannot determine the Max of an empty slice")
		return
	}
	m := 0
	for i := 1; i < l; i++ {
		if !less(rcv[i], rcv[m]) {
			m = i
		}
	}
	result = rcv[m]
	return
}

// Returns an element of Movies containing the minimum value, when compared to other elements using a passed func defining ‘less’. In the case of multiple items being equally minimal, the first such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#Min
func (rcv Movies) Min(less func(*Movie, *Movie) bool) (result *Movie, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("Cannot determine the Min of an empty slice")
		return
	}
	m := 0
	for i := 1; i < l; i++ {
		if less(rcv[i], rcv[m]) {
			m = i
		}
	}
	result = rcv[m]
	return
}

// Returns exactly one element of Movies that returns true for the passed func. Returns error if no or multiple elements return true. See: http://clipperhouse.github.io/gen/#Single
func (rcv Movies) Single(fn func(*Movie) bool) (result *Movie, err error) {
	var candidate *Movie
	found := false
	for _, v := range rcv {
		if fn(v) {
			if found {
				err = errors.New("Multiple Movies elements return true for passed func")
				return
			}
			candidate = v
			found = true
		}
	}
	if found {
		result = candidate
	} else {
		err = errors.New("No Movies elements return true for passed func")
	}
	return
}

// Returns a new ordered Movies slice, determined by a func defining ‘less’. See: http://clipperhouse.github.io/gen/#Sort
func (rcv Movies) Sort(less func(*Movie, *Movie) bool) Movies {
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

// Returns a new, descending-ordered Movies slice, determined by a func defining ‘less’. See: http://clipperhouse.github.io/gen/#Sort
//
// (Note: this is implemented by negating the passed ‘less’ func, effectively testing ‘greater than or equal to’.)
func (rcv Movies) SortDesc(less func(*Movie, *Movie) bool) Movies {
	greaterOrEqual := func(a, b *Movie) bool {
		return !less(a, b)
	}
	return rcv.Sort(greaterOrEqual)
}

// Returns a new Movies slice whose elements return true for func. See: http://clipperhouse.github.io/gen/#Where
func (rcv Movies) Where(fn func(*Movie) bool) (result Movies) {
	for _, v := range rcv {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// Iterates over Movies, operating on each element while maintaining ‘state’. See: http://clipperhouse.github.io/gen/#Aggregate
func (rcv Movies) AggregateInt(fn func(int, *Movie) int) (result int) {
	for _, v := range rcv {
		result = fn(result, v)
	}
	return
}

// Sums int over all elements and divides by len(Movies). See: http://clipperhouse.github.io/gen/#Average
func (rcv Movies) AverageInt(fn func(*Movie) int) (result int, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine AverageInt of zero-length Movies")
		return
	}
	for _, v := range rcv {
		result += fn(v)
	}
	result = result / int(l)
	return
}

// Groups elements into a map keyed by int. See: http://clipperhouse.github.io/gen/#GroupBy
func (rcv Movies) GroupByInt(fn func(*Movie) int) map[int]Movies {
	result := make(map[int]Movies)
	for _, v := range rcv {
		key := fn(v)
		result[key] = append(result[key], v)
	}
	return result
}

// Selects the largest value of int in Movies. Returns error on Movies with no elements. See: http://clipperhouse.github.io/gen/#MaxCustom
func (rcv Movies) MaxInt(fn func(*Movie) int) (result int, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine MaxInt of zero-length Movies")
		return
	}
	result = fn(rcv[0])
	if l > 1 {
		for _, v := range rcv[1:] {
			f := fn(v)
			if f > result {
				result = f
			}
		}
	}
	return
}

// Selects the least value of int in Movies. Returns error on Movies with no elements. See: http://clipperhouse.github.io/gen/#MinCustom
func (rcv Movies) MinInt(fn func(*Movie) int) (result int, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine MinInt of zero-length Movies")
		return
	}
	result = fn(rcv[0])
	if l > 1 {
		for _, v := range rcv[1:] {
			f := fn(v)
			if f < result {
				result = f
			}
		}
	}
	return
}

// Returns a slice containing all values of SelectInt in Movies. See: http://clipperhouse.github.io/gen/#Select
func (rcv Movies) SelectInt(fn func(*Movie) int) (result []int) {
	for _, v := range rcv {
		result = append(result, fn(v))
	}
	return
}

// Sums SumInt over all elements in Movies. See: http://clipperhouse.github.io/gen/#Sum
func (rcv Movies) SumInt(fn func(*Movie) int) (result int) {
	for _, v := range rcv {
		result += fn(v)
	}
	return
}

// Sort support methods

func swapMovies(rcv Movies, a, b int) {
	rcv[a], rcv[b] = rcv[b], rcv[a]
}

// Insertion sort
func insertionSortMovies(rcv Movies, less func(*Movie, *Movie) bool, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && less(rcv[j], rcv[j-1]); j-- {
			swapMovies(rcv, j, j-1)
		}
	}
}

// siftDown implements the heap property on rcv[lo, hi).
// first is an offset into the array where the root of the heap lies.
func siftDownMovies(rcv Movies, less func(*Movie, *Movie) bool, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && less(rcv[first+child], rcv[first+child+1]) {
			child++
		}
		if !less(rcv[first+root], rcv[first+child]) {
			return
		}
		swapMovies(rcv, first+root, first+child)
		root = child
	}
}

func heapSortMovies(rcv Movies, less func(*Movie, *Movie) bool, a, b int) {
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
func medianOfThreeMovies(rcv Movies, less func(*Movie, *Movie) bool, a, b, c int) {
	m0 := b
	m1 := a
	m2 := c
	// bubble sort on 3 elements
	if less(rcv[m1], rcv[m0]) {
		swapMovies(rcv, m1, m0)
	}
	if less(rcv[m2], rcv[m1]) {
		swapMovies(rcv, m2, m1)
	}
	if less(rcv[m1], rcv[m0]) {
		swapMovies(rcv, m1, m0)
	}
	// now rcv[m0] <= rcv[m1] <= rcv[m2]
}

func swapRangeMovies(rcv Movies, a, b, n int) {
	for i := 0; i < n; i++ {
		swapMovies(rcv, a+i, b+i)
	}
}

func doPivotMovies(rcv Movies, less func(*Movie, *Movie) bool, lo, hi int) (midlo, midhi int) {
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
			if less(rcv[b], rcv[pivot]) { // rcv[b] < pivot
				b++
			} else if !less(rcv[pivot], rcv[b]) { // rcv[b] = pivot
				swapMovies(rcv, a, b)
				a++
				b++
			} else {
				break
			}
		}
		for b < c {
			if less(rcv[pivot], rcv[c-1]) { // rcv[c-1] > pivot
				c--
			} else if !less(rcv[c-1], rcv[pivot]) { // rcv[c-1] = pivot
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

func quickSortMovies(rcv Movies, less func(*Movie, *Movie) bool, a, b, maxDepth int) {
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
