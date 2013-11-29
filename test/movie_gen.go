// gen *models.Movie
// this file was auto-generated using github.com/clipperhouse/gen
// Fri, 29 Nov 2013 06:24:36 UTC

// Sort functions are a modification of http://golang.org/pkg/sort/#Sort
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package models

import "errors"

// The plural (slice) type of *Movie, for use with gen methods below. Use this type where you would use []*Movie. (This is required because slices cannot be method receivers.)
type Movies []*Movie

// Tests that all elements of Movies are true for the passed func. Example:
//	good := func(v *Movie) bool {
//		return v.Something > 42
//	}
//	allGood := myMovies.All(good)
func (rcv Movies) All(fn func(*Movie) bool) bool {
	for _, v := range rcv {
		if !fn(v) {
			return false
		}
	}
	return true
}

// Tests that one or more elements of Movies are true for the passed func. Example:
//	winner := func(v *Movie) bool {
//		return v.Placement == "winner"
//	}
//	weHaveAWinner := myMovies.Any(winner)
func (rcv Movies) Any(fn func(*Movie) bool) bool {
	for _, v := range rcv {
		if fn(v) {
			return true
		}
	}
	return false
}

// Counts the number elements of Movies that are true for the passed func. Example:
//	dracula := func(v *Movie) bool {
//		return v.IsDracula()
//	}
//	countDracula := myMovies.Count(dracula)
func (rcv Movies) Count(fn func(*Movie) bool) (result int) {
	for _, v := range rcv {
		if fn(v) {
			result++
		}
	}
	return
}

// Returns a new Movies slice whose elements are unique. Keep in mind that pointers and values have different concepts of equality, and therefore distinctness. Example:
//	snowflakes := hipsters.Distinct()
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

// Returns a new Movies slice whose elements are unique where equality is defined by a passed func. Example:
//	hairstyle := func(a *Fashionista, b *Fashionista) bool {
//		a.Hairstyle == b.Hairstyle
//	}
//	trendsetters := fashionistas.DistinctBy(hairstyle)
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

// Iterates over Movies and executes the passed func against each element.
func (rcv Movies) Each(fn func(*Movie)) {
	for _, v := range rcv {
		fn(v)
	}
}

// Returns the first element that returns true for the passed func. Returns errors if no elements return true. Example:
//	winner := func(v *Movie) bool {
//		return v.Placement == "winner"
//	}
//	theWinner, err := myMovies.First(winner)
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

// Returns the element of Movies containing the maximum value, when compared to other elements using a passed func defining ‘less’. Example:
//	byArea := func(a, b *Movie) bool {
//		return a.Area() < b.Area()
//	}
//	roomiest := myMovies.Max(byArea)
//
// In the case of multiple items being equally maximal, the last such element is returned.
// (Note: this is implemented by negating the passed ‘less’ func, effectively testing ‘greater than or equal to’.)
func (rcv Movies) Max(less func(*Movie, *Movie) bool) (result *Movie, err error) {
	if len(rcv) == 0 {
		err = errors.New("Cannot determine the Max of an empty slice")
		return
	}
	return rcv.Min(negateMovies(less))
}

// Returns the element of Movies containing the minimum value, when compared to other elements using a passed func defining ‘less’. Example:
//	byPrice := func(a, b *Movie) bool {
//		return a.Price < b.Price
//	}
//	cheapest := myMovies.Min(byPrice)
//
// In the case of multiple items being equally minimal, the first such element is returned.
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

// Returns exactly one element that returns true for the passed func. Returns errors if no or multiple elements return true. Example:
//	byId := func(v *Movie) bool {
//		return v.Id == 5
//	}
//	single, err := myMovies.Single(byId)
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

// Returns a new Movies slice whose elements return true for func. Example:
//	incredible := func(v *Movie) bool {
//		return v.Manufacturer == "Apple"
//	}
//	wishList := myMovies.Where(incredible)
func (rcv Movies) Where(fn func(*Movie) bool) (result Movies) {
	for _, v := range rcv {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// Returns a new ordered Movies slice, determined by a func defining ‘less’. Example:
//	byName := func(a, b *Movie) bool {
//		return a.LastName < b.LastName
//	}
//	roster := myMovies.Sort(byName)
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

// Reports whether an instance of Movies is sorted, using the pass func to define ‘less’. See Sort method below.
func (rcv Movies) IsSorted(less func(*Movie, *Movie) bool) bool {
	n := len(rcv)
	for i := n - 1; i > 0; i-- {
		if less(rcv[i], rcv[i-1]) {
			return false
		}
	}
	return true
}

// Returns a new, descending-ordered Movies slice, determined by a func defining ‘less’. Example:
//	byPoints := func(vs Movies, a int, b int) bool {
//		return vs[a].Points < vs[b].Points
//	}
//	leaderboard := myMovies.SortDesc(byPoints)
// (Note: this is implemented by negating the passed ‘less’ func, effectively testing ‘greater than or equal to’.)
func (rcv Movies) SortDesc(less func(*Movie, *Movie) bool) Movies {
	return rcv.Sort(negateMovies(less))
}

// Reports whether an instance of Movies is sorted in descending order, using the pass func to define ‘less’. See SortDesc method below.
func (rcv Movies) IsSortedDesc(less func(*Movie, *Movie) bool) bool {
	return rcv.IsSorted(negateMovies(less))
}

func (rcv Movies) SelectTitle() (result []string) {
	for _, v := range rcv {
		result = append(result, v.Title)
	}
	return
}

func (rcv Movies) AggregateTheaters(fn func(int, int) int) (result int) {
	for _, v := range rcv {
		result = fn(result, v.Theaters)
	}
	return
}

func (rcv Movies) SumTheaters() (result int) {
	for _, v := range rcv {
		result += v.Theaters
	}
	return
}

func (rcv Movies) GroupByStudio() map[string]Movies {
	result := make(map[string]Movies)
	for _, v := range rcv {
		result[v.Studio] = append(result[v.Studio], v)
	}
	return result
}

func (rcv Movies) GroupByBoxOfficeMillions() map[int]Movies {
	result := make(map[int]Movies)
	for _, v := range rcv {
		result[v.BoxOfficeMillions] = append(result[v.BoxOfficeMillions], v)
	}
	return result
}

// ====================
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

func negateMovies(less func(*Movie, *Movie) bool) func(*Movie, *Movie) bool {
	return func(a, b *Movie) bool {
		return !less(a, b)
	}
}
