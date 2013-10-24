// gen *models.Movie
// this file was auto-generated using github.com/clipperhouse/gen
// Thu, 24 Oct 2013 02:23:53 UTC

package models

type Movies []*Movie

func (rcv Movies) AggregateInt(fn func(*Movie, int) int) (result int) {
	for _, _item := range rcv {
		result = fn(_item, result)
	}
	return result
}

func (rcv Movies) AggregateString(fn func(*Movie, string) string) (result string) {
	for _, _item := range rcv {
		result = fn(_item, result)
	}
	return result
}

func (rcv Movies) All(fn func(*Movie) bool) bool {
	for _, _item := range rcv {
		if !fn(_item) {
			return false
		}
	}
	return true
}

func (rcv Movies) Any(fn func(*Movie) bool) bool {
	for _, _item := range rcv {
		if fn(_item) {
			return true
		}
	}
	return false
}

func (rcv Movies) Count(fn func(*Movie) bool) int {
	var count = func(_item *Movie, acc int) int {
		if fn(_item) {
			acc++
		}
		return acc
	}
	return rcv.AggregateInt(count)
}

func (rcv Movies) Each(fn func(*Movie)) {
	for _, _item := range rcv {
		fn(_item)
	}
}

func (rcv Movies) JoinString(fn func(*Movie) string, delimiter string) string {
	var join = func(_item *Movie, acc string) string {
		if _item != rcv[0] {
			acc += delimiter
		}
		return acc + fn(_item)
	}
	return rcv.AggregateString(join)
}

func (rcv Movies) SumInt(fn func(*Movie) int) int {
	var sum = func(_item *Movie, acc int) int {
		return acc + fn(_item)
	}
	return rcv.AggregateInt(sum)
}

func (rcv Movies) Where(fn func(*Movie) bool) (result Movies) {
	for _, _item := range rcv {
		if fn(_item) {
			result = append(result, _item)
		}
	}
	return result
}

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

func swap(rcv Movies, a, b int) {
	rcv[a], rcv[b] = rcv[b], rcv[a]
}

func less(rcv Movies, fn func(*Movie) string, a, b int) bool {
	return fn(rcv[a]) < fn(rcv[b])
}

// Insertion sort
func insertionSort(rcv Movies, fn func(*Movie) string, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && less(rcv, fn, j, j-1); j-- {
			swap(rcv, j, j-1)
		}
	}
}

// siftDown implements the heap property on rcv[lo, hi).
// first is an offset into the array where the root of the heap lies.
func siftDown(rcv Movies, fn func(*Movie) string, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && less(rcv, fn, first+child, first+child+1) {
			child++
		}
		if !less(rcv, fn, first+root, first+child) {
			return
		}
		swap(rcv, first+root, first+child)
		root = child
	}
}

func heapSort(rcv Movies, fn func(*Movie) string, a, b int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown(rcv, fn, i, hi, first)
	}

	// Pop elements, largest first, into end of rcv.
	for i := hi - 1; i >= 0; i-- {
		swap(rcv, first, first+i)
		siftDown(rcv, fn, lo, i, first)
	}
}

// Quicksort, following Bentley and McIlroy,
// Engineering a Sort Function, SP&E November 1993.

// medianOfThree moves the median of the three values rcv[a], rcv[b], rcv[c] into rcv[a].
func medianOfThree(rcv Movies, fn func(*Movie) string, a, b, c int) {
	m0 := b
	m1 := a
	m2 := c
	// bubble sort on 3 elements
	if less(rcv, fn, m1, m0) {
		swap(rcv, m1, m0)
	}
	if less(rcv, fn, m2, m1) {
		swap(rcv, m2, m1)
	}
	if less(rcv, fn, m1, m0) {
		swap(rcv, m1, m0)
	}
	// now rcv[m0] <= rcv[m1] <= rcv[m2]
}

func swapRange(rcv Movies, a, b, n int) {
	for i := 0; i < n; i++ {
		swap(rcv, a+i, b+i)
	}
}

func doPivot(rcv Movies, fn func(*Movie) string, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2 // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		// Tukey's Ninther, median of three medians of three.
		s := (hi - lo) / 8
		medianOfThree(rcv, fn, lo, lo+s, lo+2*s)
		medianOfThree(rcv, fn, m, m-s, m+s)
		medianOfThree(rcv, fn, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThree(rcv, fn, lo, m, hi-1)

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
			if less(rcv, fn, b, pivot) { // rcv[b] < pivot
				b++
			} else if !less(rcv, fn, pivot, b) { // rcv[b] = pivot
				swap(rcv, a, b)
				a++
				b++
			} else {
				break
			}
		}
		for b < c {
			if less(rcv, fn, pivot, c-1) { // rcv[c-1] > pivot
				c--
			} else if !less(rcv, fn, c-1, pivot) { // rcv[c-1] = pivot
				swap(rcv, c-1, d-1)
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
		swap(rcv, b, c-1)
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
	swapRange(rcv, lo, b-n, n)

	n = min(hi-d, d-c)
	swapRange(rcv, c, hi-n, n)

	return lo + b - a, hi - (d - c)
}

func quickSort(rcv Movies, fn func(*Movie) string, a, b, maxDepth int) {
	for b-a > 7 {
		if maxDepth == 0 {
			heapSort(rcv, fn, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivot(rcv, fn, a, b)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			quickSort(rcv, fn, a, mlo, maxDepth)
			a = mhi // i.e., quickSort(rcv, mhi, b)
		} else {
			quickSort(rcv, fn, mhi, b, maxDepth)
			b = mlo // i.e., quickSort(rcv, a, mlo)
		}
	}
	if b-a > 1 {
		insertionSort(rcv, fn, a, b)
	}
}

func (rcv Movies) SortByString(fn func(*Movie) string) Movies {
	result := make(Movies, len(rcv))
	copy(result, rcv)

	// Switch to heapsort if depth of 2*ceil(lg(n+1)) is reached.
	n := len(result)
	maxDepth := 0
	for i := n; i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSort(result, fn, 0, n, maxDepth)
	return result
}

// IsSorted reports whether rcv is sorted.
func (rcv Movies) IsSorted(fn func(*Movie) string) bool {
	n := len(rcv)
	for i := n - 1; i > 0; i-- {
		if less(rcv, fn, i, i-1) {
			return false
		}
	}
	return true
}
