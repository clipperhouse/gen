// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package models

func swap(data Movies, a, b int) {
	data[a], data[b] = data[b], data[a]
}

func less(data Movies, fn func(*Movie) string, a, b int) bool {
	return fn(data[a]) < fn(data[b])
}

// Insertion sort
func insertionSort(data Movies, fn func(*Movie) string, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && less(data, fn, j, j-1); j-- {
			swap(data, j, j-1)
		}
	}
}

// siftDown implements the heap property on data[lo, hi).
// first is an offset into the array where the root of the heap lies.
func siftDown(data Movies, fn func(*Movie) string, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && less(data, fn, first+child, first+child+1) {
			child++
		}
		if !less(data, fn, first+root, first+child) {
			return
		}
		swap(data, first+root, first+child)
		root = child
	}
}

func heapSort(data Movies, fn func(*Movie) string, a, b int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown(data, fn, i, hi, first)
	}

	// Pop elements, largest first, into end of data.
	for i := hi - 1; i >= 0; i-- {
		swap(data, first, first+i)
		siftDown(data, fn, lo, i, first)
	}
}

// Quicksort, following Bentley and McIlroy,
// ``Engineering a Sort Function,'' SP&E November 1993.

// medianOfThree moves the median of the three values data[a], data[b], data[c] into data[a].
func medianOfThree(data Movies, fn func(*Movie) string, a, b, c int) {
	m0 := b
	m1 := a
	m2 := c
	// bubble sort on 3 elements
	if less(data, fn, m1, m0) {
		swap(data, m1, m0)
	}
	if less(data, fn, m2, m1) {
		swap(data, m2, m1)
	}
	if less(data, fn, m1, m0) {
		swap(data, m1, m0)
	}
	// now data[m0] <= data[m1] <= data[m2]
}

func swapRange(data Movies, a, b, n int) {
	for i := 0; i < n; i++ {
		swap(data, a+i, b+i)
	}
}

func doPivot(data Movies, fn func(*Movie) string, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2 // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		// Tukey's ``Ninther,'' median of three medians of three.
		s := (hi - lo) / 8
		medianOfThree(data, fn, lo, lo+s, lo+2*s)
		medianOfThree(data, fn, m, m-s, m+s)
		medianOfThree(data, fn, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThree(data, fn, lo, m, hi-1)

	// Invariants are:
	//	data[lo] = pivot (set up by ChoosePivot)
	//	data[lo <= i < a] = pivot
	//	data[a <= i < b] < pivot
	//	data[b <= i < c] is unexamined
	//	data[c <= i < d] > pivot
	//	data[d <= i < hi] = pivot
	//
	// Once b meets c, can swap the "= pivot" sections
	// into the middle of the slice.
	pivot := lo
	a, b, c, d := lo+1, lo+1, hi, hi
	for {
		for b < c {
			if less(data, fn, b, pivot) { // data[b] < pivot
				b++
			} else if !less(data, fn, pivot, b) { // data[b] = pivot
				swap(data, a, b)
				a++
				b++
			} else {
				break
			}
		}
		for b < c {
			if less(data, fn, pivot, c-1) { // data[c-1] > pivot
				c--
			} else if !less(data, fn, c-1, pivot) { // data[c-1] = pivot
				swap(data, c-1, d-1)
				c--
				d--
			} else {
				break
			}
		}
		if b >= c {
			break
		}
		// data[b] > pivot; data[c-1] < pivot
		swap(data, b, c-1)
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
	swapRange(data, lo, b-n, n)

	n = min(hi-d, d-c)
	swapRange(data, c, hi-n, n)

	return lo + b - a, hi - (d - c)
}

func quickSort(data Movies, fn func(*Movie) string, a, b, maxDepth int) {
	for b-a > 7 {
		if maxDepth == 0 {
			heapSort(data, fn, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivot(data, fn, a, b)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			quickSort(data, fn, a, mlo, maxDepth)
			a = mhi // i.e., quickSort(data, mhi, b)
		} else {
			quickSort(data, fn, mhi, b, maxDepth)
			b = mlo // i.e., quickSort(data, a, mlo)
		}
	}
	if b-a > 1 {
		insertionSort(data, fn, a, b)
	}
}

// Sort sorts data.
// It makes one call to data.Len to determine n, and O(n*log(n)) calls to
// data.Less and data.Swap. The sort is not guaranteed to be stable.
func (data Movies) SortBy(fn func(*Movie) string) Movies {
	result := make(Movies, len(data))
	copy(result, data)

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

// IsSorted reports whether data is sorted.
func (data Movies) IsSorted(fn func(*Movie) string) bool {
	n := len(data)
	for i := n - 1; i > 0; i-- {
		if less(data, fn, i, i-1) {
			return false
		}
	}
	return true
}
