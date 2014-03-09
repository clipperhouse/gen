// This file was auto-generated using github.com/clipperhouse/gen
// Modifying this file is not recommended as it will likely be overwritten in the future

// Sort (if included below) is a modification of http://golang.org/pkg/sort/#Sort
// List (if included below) is a modification of http://golang.org/pkg/container/list/
// Ring (if included below) is a modification of http://golang.org/pkg/container/ring/
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Set (if included below) is a modification of https://github.com/deckarep/golang-set
// The MIT License (MIT)
// Copyright (c) 2013 Ralph Caraveo (deckarep@gmail.com)

package models

import (
	"errors"
)

// Thing1s is a slice of type Thing1, for use with gen methods below. Use this type where you would use []Thing1. (This is required because slices cannot be method receivers.)
type Thing1s []Thing1

// All verifies that all elements of Thing1s return true for the passed func. See: http://clipperhouse.github.io/gen/#All
func (rcv Thing1s) All(fn func(Thing1) bool) bool {
	for _, v := range rcv {
		if !fn(v) {
			return false
		}
	}
	return true
}

// Any verifies that one or more elements of Thing1s return true for the passed func. See: http://clipperhouse.github.io/gen/#Any
func (rcv Thing1s) Any(fn func(Thing1) bool) bool {
	for _, v := range rcv {
		if fn(v) {
			return true
		}
	}
	return false
}

// Count gives the number elements of Thing1s that return true for the passed func. See: http://clipperhouse.github.io/gen/#Count
func (rcv Thing1s) Count(fn func(Thing1) bool) (result int) {
	for _, v := range rcv {
		if fn(v) {
			result++
		}
	}
	return
}

// Distinct returns a new Thing1s slice whose elements are unique. See: http://clipperhouse.github.io/gen/#Distinct
func (rcv Thing1s) Distinct() (result Thing1s) {
	appended := make(map[Thing1]bool)
	for _, v := range rcv {
		if !appended[v] {
			result = append(result, v)
			appended[v] = true
		}
	}
	return result
}

// DistinctBy returns a new Thing1s slice whose elements are unique, where equality is defined by a passed func. See: http://clipperhouse.github.io/gen/#DistinctBy
func (rcv Thing1s) DistinctBy(equal func(Thing1, Thing1) bool) (result Thing1s) {
	for _, v := range rcv {
		eq := func(_app Thing1) bool {
			return equal(v, _app)
		}
		if !result.Any(eq) {
			result = append(result, v)
		}
	}
	return result
}

// Each iterates over Thing1s and executes the passed func against each element. See: http://clipperhouse.github.io/gen/#Each
func (rcv Thing1s) Each(fn func(Thing1)) {
	for _, v := range rcv {
		fn(v)
	}
}

// First returns the first element that returns true for the passed func. Returns error if no elements return true. See: http://clipperhouse.github.io/gen/#First
func (rcv Thing1s) First(fn func(Thing1) bool) (result Thing1, err error) {
	for _, v := range rcv {
		if fn(v) {
			result = v
			return
		}
	}
	err = errors.New("no Thing1s elements return true for passed func")
	return
}

// IsSortedBy reports whether an instance of Thing1s is sorted, using the pass func to define ‘less’. See: http://clipperhouse.github.io/gen/#SortBy
func (rcv Thing1s) IsSortedBy(less func(Thing1, Thing1) bool) bool {
	n := len(rcv)
	for i := n - 1; i > 0; i-- {
		if less(rcv[i], rcv[i-1]) {
			return false
		}
	}
	return true
}

// IsSortedDesc reports whether an instance of Thing1s is sorted in descending order, using the pass func to define ‘less’. See: http://clipperhouse.github.io/gen/#SortBy
func (rcv Thing1s) IsSortedByDesc(less func(Thing1, Thing1) bool) bool {
	greater := func(a, b Thing1) bool {
		return a != b && !less(a, b)
	}
	return rcv.IsSortedBy(greater)
}

// MaxBy returns an element of Thing1s containing the maximum value, when compared to other elements using a passed func defining ‘less’. In the case of multiple items being equally maximal, the last such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#MaxBy
func (rcv Thing1s) MaxBy(less func(Thing1, Thing1) bool) (result Thing1, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine the MaxBy of an empty slice")
		return
	}
	m := 0
	for i := 1; i < l; i++ {
		if rcv[i] != rcv[m] && !less(rcv[i], rcv[m]) {
			m = i
		}
	}
	result = rcv[m]
	return
}

// MinBy returns an element of Thing1s containing the minimum value, when compared to other elements using a passed func defining ‘less’. In the case of multiple items being equally minimal, the first such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#MinBy
func (rcv Thing1s) MinBy(less func(Thing1, Thing1) bool) (result Thing1, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine the Min of an empty slice")
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

// Single returns exactly one element of Thing1s that returns true for the passed func. Returns error if no or multiple elements return true. See: http://clipperhouse.github.io/gen/#Single
func (rcv Thing1s) Single(fn func(Thing1) bool) (result Thing1, err error) {
	var candidate Thing1
	found := false
	for _, v := range rcv {
		if fn(v) {
			if found {
				err = errors.New("multiple Thing1s elements return true for passed func")
				return
			}
			candidate = v
			found = true
		}
	}
	if found {
		result = candidate
	} else {
		err = errors.New("no Thing1s elements return true for passed func")
	}
	return
}

// SortBy returns a new ordered Thing1s slice, determined by a func defining ‘less’. See: http://clipperhouse.github.io/gen/#SortBy
func (rcv Thing1s) SortBy(less func(Thing1, Thing1) bool) Thing1s {
	result := make(Thing1s, len(rcv))
	copy(result, rcv)
	// Switch to heapsort if depth of 2*ceil(lg(n+1)) is reached.
	n := len(result)
	maxDepth := 0
	for i := n; i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSortThing1s(result, less, 0, n, maxDepth)
	return result
}

// SortByDesc returns a new, descending-ordered Thing1s slice, determined by a func defining ‘less’. See: http://clipperhouse.github.io/gen/#SortBy
func (rcv Thing1s) SortByDesc(less func(Thing1, Thing1) bool) Thing1s {
	greater := func(a, b Thing1) bool {
		return a != b && !less(a, b)
	}
	return rcv.SortBy(greater)
}

// Where returns a new Thing1s slice whose elements return true for func. See: http://clipperhouse.github.io/gen/#Where
func (rcv Thing1s) Where(fn func(Thing1) bool) (result Thing1s) {
	for _, v := range rcv {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// AggregateInt iterates over Thing1s, operating on each element while maintaining ‘state’. See: http://clipperhouse.github.io/gen/#Aggregate
func (rcv Thing1s) AggregateInt(fn func(int, Thing1) int) (result int) {
	for _, v := range rcv {
		result = fn(result, v)
	}
	return
}

// AverageInt sums int over all elements and divides by len(Thing1s). See: http://clipperhouse.github.io/gen/#Average
func (rcv Thing1s) AverageInt(fn func(Thing1) int) (result int, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine AverageInt of zero-length Thing1s")
		return
	}
	for _, v := range rcv {
		result += fn(v)
	}
	result = result / int(l)
	return
}

// GroupByInt groups elements into a map keyed by int. See: http://clipperhouse.github.io/gen/#GroupBy
func (rcv Thing1s) GroupByInt(fn func(Thing1) int) map[int]Thing1s {
	result := make(map[int]Thing1s)
	for _, v := range rcv {
		key := fn(v)
		result[key] = append(result[key], v)
	}
	return result
}

// MaxInt selects the largest value of int in Thing1s. Returns error on Thing1s with no elements. See: http://clipperhouse.github.io/gen/#MaxCustom
func (rcv Thing1s) MaxInt(fn func(Thing1) int) (result int, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine MaxInt of zero-length Thing1s")
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

// MinInt selects the least value of int in Thing1s. Returns error on Thing1s with no elements. See: http://clipperhouse.github.io/gen/#MinCustom
func (rcv Thing1s) MinInt(fn func(Thing1) int) (result int, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine MinInt of zero-length Thing1s")
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

// SelectInt returns a slice of int in Thing1s, projected by passed func. See: http://clipperhouse.github.io/gen/#Select
func (rcv Thing1s) SelectInt(fn func(Thing1) int) (result []int) {
	for _, v := range rcv {
		result = append(result, fn(v))
	}
	return
}

// SumInt sums int over elements in Thing1s. See: http://clipperhouse.github.io/gen/#Sum
func (rcv Thing1s) SumInt(fn func(Thing1) int) (result int) {
	for _, v := range rcv {
		result += fn(v)
	}
	return
}

// AggregateThing2 iterates over Thing1s, operating on each element while maintaining ‘state’. See: http://clipperhouse.github.io/gen/#Aggregate
func (rcv Thing1s) AggregateThing2(fn func(Thing2, Thing1) Thing2) (result Thing2) {
	for _, v := range rcv {
		result = fn(result, v)
	}
	return
}

// AverageThing2 sums Thing2 over all elements and divides by len(Thing1s). See: http://clipperhouse.github.io/gen/#Average
func (rcv Thing1s) AverageThing2(fn func(Thing1) Thing2) (result Thing2, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine AverageThing2 of zero-length Thing1s")
		return
	}
	for _, v := range rcv {
		result += fn(v)
	}
	result = result / Thing2(l)
	return
}

// GroupByThing2 groups elements into a map keyed by Thing2. See: http://clipperhouse.github.io/gen/#GroupBy
func (rcv Thing1s) GroupByThing2(fn func(Thing1) Thing2) map[Thing2]Thing1s {
	result := make(map[Thing2]Thing1s)
	for _, v := range rcv {
		key := fn(v)
		result[key] = append(result[key], v)
	}
	return result
}

// MaxThing2 selects the largest value of Thing2 in Thing1s. Returns error on Thing1s with no elements. See: http://clipperhouse.github.io/gen/#MaxCustom
func (rcv Thing1s) MaxThing2(fn func(Thing1) Thing2) (result Thing2, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine MaxThing2 of zero-length Thing1s")
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

// MinThing2 selects the least value of Thing2 in Thing1s. Returns error on Thing1s with no elements. See: http://clipperhouse.github.io/gen/#MinCustom
func (rcv Thing1s) MinThing2(fn func(Thing1) Thing2) (result Thing2, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine MinThing2 of zero-length Thing1s")
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

// SelectThing2 returns a slice of Thing2 in Thing1s, projected by passed func. See: http://clipperhouse.github.io/gen/#Select
func (rcv Thing1s) SelectThing2(fn func(Thing1) Thing2) (result []Thing2) {
	for _, v := range rcv {
		result = append(result, fn(v))
	}
	return
}

// SumThing2 sums Thing2 over elements in Thing1s. See: http://clipperhouse.github.io/gen/#Sum
func (rcv Thing1s) SumThing2(fn func(Thing1) Thing2) (result Thing2) {
	for _, v := range rcv {
		result += fn(v)
	}
	return
}

// AggregateString iterates over Thing1s, operating on each element while maintaining ‘state’. See: http://clipperhouse.github.io/gen/#Aggregate
func (rcv Thing1s) AggregateString(fn func(string, Thing1) string) (result string) {
	for _, v := range rcv {
		result = fn(result, v)
	}
	return
}

// GroupByString groups elements into a map keyed by string. See: http://clipperhouse.github.io/gen/#GroupBy
func (rcv Thing1s) GroupByString(fn func(Thing1) string) map[string]Thing1s {
	result := make(map[string]Thing1s)
	for _, v := range rcv {
		key := fn(v)
		result[key] = append(result[key], v)
	}
	return result
}

// MaxString selects the largest value of string in Thing1s. Returns error on Thing1s with no elements. See: http://clipperhouse.github.io/gen/#MaxCustom
func (rcv Thing1s) MaxString(fn func(Thing1) string) (result string, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine MaxString of zero-length Thing1s")
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

// MinString selects the least value of string in Thing1s. Returns error on Thing1s with no elements. See: http://clipperhouse.github.io/gen/#MinCustom
func (rcv Thing1s) MinString(fn func(Thing1) string) (result string, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine MinString of zero-length Thing1s")
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

// SelectString returns a slice of string in Thing1s, projected by passed func. See: http://clipperhouse.github.io/gen/#Select
func (rcv Thing1s) SelectString(fn func(Thing1) string) (result []string) {
	for _, v := range rcv {
		result = append(result, fn(v))
	}
	return
}

// Sort support methods

func swapThing1s(rcv Thing1s, a, b int) {
	rcv[a], rcv[b] = rcv[b], rcv[a]
}

// Insertion sort
func insertionSortThing1s(rcv Thing1s, less func(Thing1, Thing1) bool, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && less(rcv[j], rcv[j-1]); j-- {
			swapThing1s(rcv, j, j-1)
		}
	}
}

// siftDown implements the heap property on rcv[lo, hi).
// first is an offset into the array where the root of the heap lies.
func siftDownThing1s(rcv Thing1s, less func(Thing1, Thing1) bool, lo, hi, first int) {
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
		swapThing1s(rcv, first+root, first+child)
		root = child
	}
}

func heapSortThing1s(rcv Thing1s, less func(Thing1, Thing1) bool, a, b int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDownThing1s(rcv, less, i, hi, first)
	}

	// Pop elements, largest first, into end of rcv.
	for i := hi - 1; i >= 0; i-- {
		swapThing1s(rcv, first, first+i)
		siftDownThing1s(rcv, less, lo, i, first)
	}
}

// Quicksort, following Bentley and McIlroy,
// Engineering a Sort Function, SP&E November 1993.

// medianOfThree moves the median of the three values rcv[a], rcv[b], rcv[c] into rcv[a].
func medianOfThreeThing1s(rcv Thing1s, less func(Thing1, Thing1) bool, a, b, c int) {
	m0 := b
	m1 := a
	m2 := c
	// bubble sort on 3 elements
	if less(rcv[m1], rcv[m0]) {
		swapThing1s(rcv, m1, m0)
	}
	if less(rcv[m2], rcv[m1]) {
		swapThing1s(rcv, m2, m1)
	}
	if less(rcv[m1], rcv[m0]) {
		swapThing1s(rcv, m1, m0)
	}
	// now rcv[m0] <= rcv[m1] <= rcv[m2]
}

func swapRangeThing1s(rcv Thing1s, a, b, n int) {
	for i := 0; i < n; i++ {
		swapThing1s(rcv, a+i, b+i)
	}
}

func doPivotThing1s(rcv Thing1s, less func(Thing1, Thing1) bool, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2 // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		// Tukey's Ninther, median of three medians of three.
		s := (hi - lo) / 8
		medianOfThreeThing1s(rcv, less, lo, lo+s, lo+2*s)
		medianOfThreeThing1s(rcv, less, m, m-s, m+s)
		medianOfThreeThing1s(rcv, less, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThreeThing1s(rcv, less, lo, m, hi-1)

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
				swapThing1s(rcv, a, b)
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
				swapThing1s(rcv, c-1, d-1)
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
		swapThing1s(rcv, b, c-1)
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
	swapRangeThing1s(rcv, lo, b-n, n)

	n = min(hi-d, d-c)
	swapRangeThing1s(rcv, c, hi-n, n)

	return lo + b - a, hi - (d - c)
}

func quickSortThing1s(rcv Thing1s, less func(Thing1, Thing1) bool, a, b, maxDepth int) {
	for b-a > 7 {
		if maxDepth == 0 {
			heapSortThing1s(rcv, less, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivotThing1s(rcv, less, a, b)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			quickSortThing1s(rcv, less, a, mlo, maxDepth)
			a = mhi // i.e., quickSortThing1s(rcv, mhi, b)
		} else {
			quickSortThing1s(rcv, less, mhi, b, maxDepth)
			b = mlo // i.e., quickSortThing1s(rcv, a, mlo)
		}
	}
	if b-a > 1 {
		insertionSortThing1s(rcv, less, a, b)
	}
}
