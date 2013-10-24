package main

import (
	"text/template"
)

func getTemplate() *template.Template {
	return template.Must(template.New("gen").Parse(tmpl))
}

const tmpl = `// {{.Command}}
// this file was auto-generated using github.com/clipperhouse/gen
// {{.Generated}}

package {{.Package}}

type {{.Plural}} []{{.Pointer}}{{.Singular}}

func ({{.Receiver}} {{.Plural}}) AggregateInt(fn func({{.Pointer}}{{.Singular}}, int) int) (result int) {
	for _, {{.Loop}} := range {{.Receiver}} {
		result = fn({{.Loop}}, result)
	}
	return result
}

func ({{.Receiver}} {{.Plural}}) AggregateString(fn func({{.Pointer}}{{.Singular}}, string) string) (result string) {
	for _, {{.Loop}} := range {{.Receiver}} {
		result = fn({{.Loop}}, result)
	}
	return result
}

func ({{.Receiver}} {{.Plural}}) All(fn func({{.Pointer}}{{.Singular}}) bool) bool {
	for _, {{.Loop}} := range {{.Receiver}} {
		if !fn({{.Loop}}) {
			return false
		}
	}
	return true
}

func ({{.Receiver}} {{.Plural}}) Any(fn func({{.Pointer}}{{.Singular}}) bool) bool {
	for _, {{.Loop}} := range {{.Receiver}} {
		if fn({{.Loop}}) {
			return true
		}
	}
	return false
}

func ({{.Receiver}} {{.Plural}}) Count(fn func({{.Pointer}}{{.Singular}}) bool) int {
	var count = func({{.Loop}} {{.Pointer}}{{.Singular}}, acc int) int {
		if fn({{.Loop}}) {
			acc++
		}
		return acc
	}
	return {{.Receiver}}.AggregateInt(count)
}

func ({{.Receiver}} {{.Plural}}) Each(fn func({{.Pointer}}{{.Singular}})) {
	for _, {{.Loop}} := range {{.Receiver}} {
		fn({{.Loop}})
	}
}

func ({{.Receiver}} {{.Plural}}) JoinString(fn func({{.Pointer}}{{.Singular}}) string, delimiter string) string {
	var join = func({{.Loop}} {{.Pointer}}{{.Singular}}, acc string) string {
		if {{.Loop}} != {{.Receiver}}[0] {
			acc += delimiter
		}
		return acc + fn({{.Loop}})
	}
	return {{.Receiver}}.AggregateString(join)
}

func ({{.Receiver}} {{.Plural}}) SumInt(fn func({{.Pointer}}{{.Singular}}) int) int {
	var sum = func({{.Loop}} {{.Pointer}}{{.Singular}}, acc int) int {
		return acc + fn({{.Loop}})
	}
	return {{.Receiver}}.AggregateInt(sum)
}

func ({{.Receiver}} {{.Plural}}) Where(fn func({{.Pointer}}{{.Singular}}) bool) (result {{.Plural}}) {
	for _, {{.Loop}} := range {{.Receiver}} {
		if fn({{.Loop}}) {
			result = append(result, {{.Loop}})
		}
	}
	return result
}

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

func swap({{.Receiver}} {{.Plural}}, a, b int) {
	{{.Receiver}}[a], {{.Receiver}}[b] = {{.Receiver}}[b], {{.Receiver}}[a]
}

func less({{.Receiver}} {{.Plural}}, fn func({{.Pointer}}{{.Singular}}) string, a, b int) bool {
	return fn({{.Receiver}}[a]) < fn({{.Receiver}}[b])
}

// Insertion sort
func insertionSort({{.Receiver}} {{.Plural}}, fn func({{.Pointer}}{{.Singular}}) string, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && less({{.Receiver}}, fn, j, j-1); j-- {
			swap({{.Receiver}}, j, j-1)
		}
	}
}

// siftDown implements the heap property on {{.Receiver}}[lo, hi).
// first is an offset into the array where the root of the heap lies.
func siftDown({{.Receiver}} {{.Plural}}, fn func({{.Pointer}}{{.Singular}}) string, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && less({{.Receiver}}, fn, first+child, first+child+1) {
			child++
		}
		if !less({{.Receiver}}, fn, first+root, first+child) {
			return
		}
		swap({{.Receiver}}, first+root, first+child)
		root = child
	}
}

func heapSort({{.Receiver}} {{.Plural}}, fn func({{.Pointer}}{{.Singular}}) string, a, b int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown({{.Receiver}}, fn, i, hi, first)
	}

	// Pop elements, largest first, into end of {{.Receiver}}.
	for i := hi - 1; i >= 0; i-- {
		swap({{.Receiver}}, first, first+i)
		siftDown({{.Receiver}}, fn, lo, i, first)
	}
}

// Quicksort, following Bentley and McIlroy,
// Engineering a Sort Function, SP&E November 1993.

// medianOfThree moves the median of the three values {{.Receiver}}[a], {{.Receiver}}[b], {{.Receiver}}[c] into {{.Receiver}}[a].
func medianOfThree({{.Receiver}} {{.Plural}}, fn func({{.Pointer}}{{.Singular}}) string, a, b, c int) {
	m0 := b
	m1 := a
	m2 := c
	// bubble sort on 3 elements
	if less({{.Receiver}}, fn, m1, m0) {
		swap({{.Receiver}}, m1, m0)
	}
	if less({{.Receiver}}, fn, m2, m1) {
		swap({{.Receiver}}, m2, m1)
	}
	if less({{.Receiver}}, fn, m1, m0) {
		swap({{.Receiver}}, m1, m0)
	}
	// now {{.Receiver}}[m0] <= {{.Receiver}}[m1] <= {{.Receiver}}[m2]
}

func swapRange({{.Receiver}} {{.Plural}}, a, b, n int) {
	for i := 0; i < n; i++ {
		swap({{.Receiver}}, a+i, b+i)
	}
}

func doPivot({{.Receiver}} {{.Plural}}, fn func({{.Pointer}}{{.Singular}}) string, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2 // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		// Tukey's Ninther, median of three medians of three.
		s := (hi - lo) / 8
		medianOfThree({{.Receiver}}, fn, lo, lo+s, lo+2*s)
		medianOfThree({{.Receiver}}, fn, m, m-s, m+s)
		medianOfThree({{.Receiver}}, fn, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThree({{.Receiver}}, fn, lo, m, hi-1)

	// Invariants are:
	//	{{.Receiver}}[lo] = pivot (set up by ChoosePivot)
	//	{{.Receiver}}[lo <= i < a] = pivot
	//	{{.Receiver}}[a <= i < b] < pivot
	//	{{.Receiver}}[b <= i < c] is unexamined
	//	{{.Receiver}}[c <= i < d] > pivot
	//	{{.Receiver}}[d <= i < hi] = pivot
	//
	// Once b meets c, can swap the "= pivot" sections
	// into the middle of the slice.
	pivot := lo
	a, b, c, d := lo+1, lo+1, hi, hi
	for {
		for b < c {
			if less({{.Receiver}}, fn, b, pivot) { // {{.Receiver}}[b] < pivot
				b++
			} else if !less({{.Receiver}}, fn, pivot, b) { // {{.Receiver}}[b] = pivot
				swap({{.Receiver}}, a, b)
				a++
				b++
			} else {
				break
			}
		}
		for b < c {
			if less({{.Receiver}}, fn, pivot, c-1) { // {{.Receiver}}[c-1] > pivot
				c--
			} else if !less({{.Receiver}}, fn, c-1, pivot) { // {{.Receiver}}[c-1] = pivot
				swap({{.Receiver}}, c-1, d-1)
				c--
				d--
			} else {
				break
			}
		}
		if b >= c {
			break
		}
		// {{.Receiver}}[b] > pivot; {{.Receiver}}[c-1] < pivot
		swap({{.Receiver}}, b, c-1)
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
	swapRange({{.Receiver}}, lo, b-n, n)

	n = min(hi-d, d-c)
	swapRange({{.Receiver}}, c, hi-n, n)

	return lo + b - a, hi - (d - c)
}

func quickSort({{.Receiver}} {{.Plural}}, fn func({{.Pointer}}{{.Singular}}) string, a, b, maxDepth int) {
	for b-a > 7 {
		if maxDepth == 0 {
			heapSort({{.Receiver}}, fn, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivot({{.Receiver}}, fn, a, b)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			quickSort({{.Receiver}}, fn, a, mlo, maxDepth)
			a = mhi // i.e., quickSort({{.Receiver}}, mhi, b)
		} else {
			quickSort({{.Receiver}}, fn, mhi, b, maxDepth)
			b = mlo // i.e., quickSort({{.Receiver}}, a, mlo)
		}
	}
	if b-a > 1 {
		insertionSort({{.Receiver}}, fn, a, b)
	}
}

func ({{.Receiver}} {{.Plural}}) SortByString(fn func({{.Pointer}}{{.Singular}}) string) {{.Plural}} {
	result := make({{.Plural}}, len({{.Receiver}}))
	copy(result, {{.Receiver}})

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

// IsSorted reports whether {{.Receiver}} is sorted.
func ({{.Receiver}} {{.Plural}}) IsSorted(fn func({{.Pointer}}{{.Singular}}) string) bool {
	n := len({{.Receiver}})
	for i := n - 1; i > 0; i-- {
		if less({{.Receiver}}, fn, i, i-1) {
			return false
		}
	}
	return true
}
`
