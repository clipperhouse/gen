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

import "errors"

// The plural (slice) type of {{.Pointer}}{{.Singular}}, for use with gen methods below. Use this type where you would use []{{.Pointer}}{{.Singular}}. (This is required because slices cannot be method receivers.)
type {{.Plural}} []{{.Pointer}}{{.Singular}}

// Iterates over {{.Plural}}, operating on each element while accumulating results. For example, Sum & Min might be implemented like:
//	sum := func(_item {{.Pointer}}{{.Singular}}, accumulated int) int {
//		return accumulated + _item.Something
//	}
//	sumOfSomething := my{{.Plural}}.AggregateInt(sum)
//
//	min := func(_item {{.Pointer}}{{.Singular}}, accumulated int) int {
//		if _item.AnotherThing < accumulated {
//			return _item.AnotherThing
//		}
//		return accumulated
//	}
//	minOfAnotherThing := my{{.Plural}}.AggregateInt(min)
func ({{.Receiver}} {{.Plural}}) AggregateInt(fn func({{.Pointer}}{{.Singular}}, int) int) int {
	result := 0
	for _, {{.Loop}} := range {{.Receiver}} {
		result = fn({{.Loop}}, result)
	}
	return result
}

// Iterates over {{.Plural}}, operating on each element while accumulating results. For example, you might join strings like:
//	my{{.Plural}} := GetSome{{.Plural}}()
//	join := func(_item {{.Pointer}}{{.Singular}}, accumulated string) string {
//		if _item != my{{.Plural}}[0] {
//			accumulated += ", "
//		}
//		return accumulated + _item.Title
//	}
//	myList := my{{.Plural}}.AggregateString(join)
func ({{.Receiver}} {{.Plural}}) AggregateString(fn func({{.Pointer}}{{.Singular}}, string) string) (result string) {
	for _, {{.Loop}} := range {{.Receiver}} {
		result = fn({{.Loop}}, result)
	}
	return result
}

// Tests that all elements of {{.Plural}} are true for the passed func. Example:
//	good := func(_item {{.Pointer}}{{.Singular}}) bool {
//		return _item.Something > 42
//	}
//	allGood := my{{.Plural}}.All(good)
func ({{.Receiver}} {{.Plural}}) All(fn func({{.Pointer}}{{.Singular}}) bool) bool {
	for _, {{.Loop}} := range {{.Receiver}} {
		if !fn({{.Loop}}) {
			return false
		}
	}
	return true
}

// Tests that one or more elements of {{.Plural}} are true for the passed func. Example:
//	winner := func(_item {{.Pointer}}{{.Singular}}) bool {
//		return _item.Placement == "winner"
//	}
//	weHaveAWinner := my{{.Plural}}.Any(winner)
func ({{.Receiver}} {{.Plural}}) Any(fn func({{.Pointer}}{{.Singular}}) bool) bool {
	for _, {{.Loop}} := range {{.Receiver}} {
		if fn({{.Loop}}) {
			return true
		}
	}
	return false
}

// Counts the number elements of {{.Plural}} that are true for the passed func. Example:
//	dracula := func(_item {{.Pointer}}{{.Singular}}) bool {
//		return _item.IsDracula()
//	}
//	countDracula := my{{.Plural}}.Count(dracula)
func ({{.Receiver}} {{.Plural}}) Count(fn func({{.Pointer}}{{.Singular}}) bool) int {
	var count = func({{.Loop}} {{.Pointer}}{{.Singular}}, acc int) int {
		if fn({{.Loop}}) {
			acc++
		}
		return acc
	}
	return {{.Receiver}}.AggregateInt(count)
}

// Returns a new {{.Plural}} slice whose elements are unique. Keep in mind that pointers and values have different concepts of equality, and therefore distinctness. Example:
//	snowflakes := hipsters.Distinct()
func ({{.Receiver}} {{.Plural}}) Distinct() (result {{.Plural}}) {
	appended := make(map[{{.Pointer}}{{.Singular}}]bool)
	for _, {{.Loop}} := range {{.Receiver}} {
		if !appended[{{.Loop}}] {
			result = append(result, {{.Loop}})
			appended[{{.Loop}}] = true
		}
	}
	return result
}

// Returns a new {{.Plural}} slice whose elements are unique where equality is defined by a passed func. Example:
//	hairstyle := func(a *Fashionista, b *Fashionista) bool {
//		a.Hairstyle == b.Hairstyle
//	}
//	trendsetters := fashionistas.DistinctBy(hairstyle)
func ({{.Receiver}} {{.Plural}}) DistinctBy(equal func({{.Pointer}}{{.Singular}}, {{.Pointer}}{{.Singular}}) bool) (result {{.Plural}}) {
	for _, {{.Loop}} := range {{.Receiver}} {
		eq := func(_app {{.Pointer}}{{.Singular}}) bool {
			return equal({{.Loop}}, _app)
		}
		if !result.Any(eq) {
			result = append(result, {{.Loop}})
		}
	}
	return result
}

// Iterates over {{.Plural}} and executes the passed func against each element.
func ({{.Receiver}} {{.Plural}}) Each(fn func({{.Pointer}}{{.Singular}})) {
	for _, {{.Loop}} := range {{.Receiver}} {
		fn({{.Loop}})
	}
}

// Returns the first element that returns true for the passed func. Returns errors if no elements return true. Example:
//	winner := func(_item {{.Pointer}}{{.Singular}}) bool {
//		return _item.Placement == "winner"
//	}
//	theWinner, err := myMovies.First(winner)
func ({{.Receiver}} {{.Plural}}) First(fn func({{.Pointer}}{{.Singular}}) bool) ({{.Pointer}}{{.Singular}}, error) {
	for _, {{.Loop}} := range {{.Receiver}} {
		if fn({{.Loop}}) {
			return {{.Loop}}, nil
		}
	}
	return nil, errors.New("No {{.Plural}} elements return true for passed func")
}

// Groups {{.Plural}} into a map of Movies, keyed by the result of the passed func. Example:
//	year := func(_item {{.Pointer}}{{.Singular}}) int {
//		return _item.CreationDate.Year()
//	}
//	yearlyReport := my{{.Plural}}.GroupByInt(year)
func ({{.Receiver}} {{.Plural}}) GroupByInt(fn func({{.Pointer}}{{.Singular}}) int) map[int]{{.Plural}} {
	result := make(map[int]{{.Plural}})
	for _, _item := range {{.Receiver}} {
		result[fn(_item)] = append(result[fn(_item)], _item)
	}
	return result
}

// Groups {{.Plural}} into a map of Movies, keyed by the result of the passed func. Example:
//	dept := func(_item {{.Pointer}}{{.Singular}}) string {
//		return _item.DepartmentName
//	}
//	byDepartment := my{{.Plural}}.GroupByString(dept)
func ({{.Receiver}} {{.Plural}}) GroupByString(fn func({{.Pointer}}{{.Singular}}) string) map[string]{{.Plural}} {
	result := make(map[string]{{.Plural}})
	for _, _item := range {{.Receiver}} {
		result[fn(_item)] = append(result[fn(_item)], _item)
	}
	return result
}

// Returns the element of {{.Plural}} containing the maximum value, when compared to other elements using a passed func defining ‘less’. Example:
//	byArea := func(_items {{.Plural}}, a int, b int) bool {
//		return _items[a].Area() < _items[b].Area()
//	}
//	roomiest := my{{.Plural}}.Max(byArea)
//
// In the case of multiple items being equally maximal, the last such element is returned.
// (Note: this is implemented by negating the passed ‘less’ func, effectively testing ‘greater than or equal to’.)
func ({{.Receiver}} {{.Plural}}) Max(less func({{.Plural}}, int, int) bool) ({{.Pointer}}{{.Singular}}, error) {
	if len(rcv) == 0 {
		return nil, errors.New("Cannot determine the Max of an empty slice")
	}
	return rcv.Min(not(less))
}

// Returns the element of {{.Plural}} containing the minimum value, when compared to other elements using a passed func defining ‘less’. Example:
//	byPrice := func(_items {{.Plural}}, a int, b int) bool {
//		return _items[a].Price < _items[b].Price
//	}
//	cheapest := my{{.Plural}}.Min(byPrice)
//
// In the case of multiple items being equally minimal, the first such element is returned.
func ({{.Receiver}} {{.Plural}}) Min(less func({{.Plural}}, int, int) bool) ({{.Pointer}}{{.Singular}}, error) {
	l := len({{.Receiver}})
	if l == 0 {
		return nil, errors.New("Cannot determine the Min of an empty slice")
	}
	m := 0
	for i := 1; i < l; i++ {
		if less({{.Receiver}}, i, m) {
			m = i
		}
	}
	return {{.Receiver}}[m], nil
}

// Returns exactly one element that returns true for the passed func. Returns errors if no or multiple elements return true. Example:
//	byId := func(_item {{.Pointer}}{{.Singular}}) bool {
//		return _item.Id == 5
//	}
//	single, err := myMovies.Single(byId)
func ({{.Receiver}} {{.Plural}}) Single(fn func({{.Pointer}}{{.Singular}}) bool) ({{.Pointer}}{{.Singular}}, error) {
	var result {{.Pointer}}{{.Singular}}
	found := false
	for _, {{.Loop}} := range {{.Receiver}} {
		if fn({{.Loop}}) {
			if found {
				return nil, errors.New("Multiple {{.Plural}} elements return true for passed func")
			}
			result = {{.Loop}}
			found = true
		}
	}
	if !found {
		return nil, errors.New("No {{.Plural}} elements return true for passed func")
	}
	return result, nil
}

// Returns the sum of ints returned by passed func. Example:
//	itemTotal := func(_item {{.Pointer}}{{.Singular}}) int {
//		return _item.Quantity * _item.UnitPrice
//	}
//	orderTotal := my{{.Plural}}.SumInt(itemTotal)
func ({{.Receiver}} {{.Plural}}) SumInt(fn func({{.Pointer}}{{.Singular}}) int) int {
	var sum = func({{.Loop}} {{.Pointer}}{{.Singular}}, acc int) int {
		return acc + fn({{.Loop}})
	}
	return {{.Receiver}}.AggregateInt(sum)
}

// Returns a new {{.Plural}} slice whose elements return true for func. Example:
//	incredible := func(_item {{.Pointer}}{{.Singular}}) bool {
//		return _item.Manufacturer == "Apple"
//	}
//	wishList := my{{.Plural}}.Where(incredible)
func ({{.Receiver}} {{.Plural}}) Where(fn func({{.Pointer}}{{.Singular}}) bool) (result {{.Plural}}) {
	for _, {{.Loop}} := range {{.Receiver}} {
		if fn({{.Loop}}) {
			result = append(result, {{.Loop}})
		}
	}
	return result
}

// Sort functions below are a modification of http://golang.org/pkg/sort/#Sort
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Returns a new ordered {{.Plural}} slice, determined by a func defining ‘less’. Example:
//	byName := func(_items {{.Plural}}, a int, b int) bool {
//		return _items[a].LastName < _items[b].LastName
//	}
//	roster := my{{.Plural}}.Sort(byName)
func ({{.Receiver}} {{.Plural}}) Sort(less func({{.Plural}}, int, int) bool) {{.Plural}} {
	result := make({{.Plural}}, len({{.Receiver}}))
	copy(result, {{.Receiver}})

	// Switch to heapsort if depth of 2*ceil(lg(n+1)) is reached.
	n := len(result)
	maxDepth := 0
	for i := n; i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSort{{.Plural}}(result, less, 0, n, maxDepth)
	return result
}

// Reports whether an instance of {{.Plural}} is sorted, using the pass func to define ‘less’. See Sort method below.
func ({{.Receiver}} {{.Plural}}) IsSorted(less func({{.Plural}}, int, int) bool) bool {
	n := len({{.Receiver}})
	for i := n - 1; i > 0; i-- {
		if less({{.Receiver}}, i, i-1) {
			return false
		}
	}
	return true
}

// Returns a new, descending-ordered {{.Plural}} slice, determined by a func defining ‘less’. Example:
//	byPoints := func(_items {{.Plural}}, a int, b int) bool {
//		return _items[a].Points < _items[b].Points
//	}
//	leaderboard := my{{.Plural}}.SortDesc(byPoints)
// (Note: this is implemented by negating the passed ‘less’ func, effectively testing ‘greater than or equal to’.)
func ({{.Receiver}} {{.Plural}}) SortDesc(less func({{.Plural}}, int, int) bool) {{.Plural}} {
	return {{.Receiver}}.Sort(not(less))
}

// Reports whether an instance of {{.Plural}} is sorted in descending order, using the pass func to define ‘less’. See SortDesc method below.
func ({{.Receiver}} {{.Plural}}) IsSortedDesc(less func({{.Plural}}, int, int) bool) bool {
	return {{.Receiver}}.IsSorted(not(less))
}

func swap{{.Plural}}({{.Receiver}} {{.Plural}}, a, b int) {
	{{.Receiver}}[a], {{.Receiver}}[b] = {{.Receiver}}[b], {{.Receiver}}[a]
}

// Insertion sort
func insertionSort{{.Plural}}({{.Receiver}} {{.Plural}}, less func({{.Plural}}, int, int) bool, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && less({{.Receiver}}, j, j-1); j-- {
			swap{{.Plural}}({{.Receiver}}, j, j-1)
		}
	}
}

// siftDown implements the heap property on {{.Receiver}}[lo, hi).
// first is an offset into the array where the root of the heap lies.
func siftDown{{.Plural}}({{.Receiver}} {{.Plural}}, less func({{.Plural}}, int, int) bool, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && less({{.Receiver}}, first+child, first+child+1) {
			child++
		}
		if !less({{.Receiver}}, first+root, first+child) {
			return
		}
		swap{{.Plural}}({{.Receiver}}, first+root, first+child)
		root = child
	}
}

func heapSort{{.Plural}}({{.Receiver}} {{.Plural}}, less func({{.Plural}}, int, int) bool, a, b int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown{{.Plural}}({{.Receiver}}, less, i, hi, first)
	}

	// Pop elements, largest first, into end of {{.Receiver}}.
	for i := hi - 1; i >= 0; i-- {
		swap{{.Plural}}({{.Receiver}}, first, first+i)
		siftDown{{.Plural}}({{.Receiver}}, less, lo, i, first)
	}
}

// Quicksort, following Bentley and McIlroy,
// Engineering a Sort Function, SP&E November 1993.

// medianOfThree moves the median of the three values {{.Receiver}}[a], {{.Receiver}}[b], {{.Receiver}}[c] into {{.Receiver}}[a].
func medianOfThree{{.Plural}}({{.Receiver}} {{.Plural}}, less func({{.Plural}}, int, int) bool, a, b, c int) {
	m0 := b
	m1 := a
	m2 := c
	// bubble sort on 3 elements
	if less({{.Receiver}}, m1, m0) {
		swap{{.Plural}}({{.Receiver}}, m1, m0)
	}
	if less({{.Receiver}}, m2, m1) {
		swap{{.Plural}}({{.Receiver}}, m2, m1)
	}
	if less({{.Receiver}}, m1, m0) {
		swap{{.Plural}}({{.Receiver}}, m1, m0)
	}
	// now {{.Receiver}}[m0] <= {{.Receiver}}[m1] <= {{.Receiver}}[m2]
}

func swapRange{{.Plural}}({{.Receiver}} {{.Plural}}, a, b, n int) {
	for i := 0; i < n; i++ {
		swap{{.Plural}}({{.Receiver}}, a+i, b+i)
	}
}

func doPivot{{.Plural}}({{.Receiver}} {{.Plural}}, less func({{.Plural}}, int, int) bool, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2 // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		// Tukey's Ninther, median of three medians of three.
		s := (hi - lo) / 8
		medianOfThree{{.Plural}}({{.Receiver}}, less, lo, lo+s, lo+2*s)
		medianOfThree{{.Plural}}({{.Receiver}}, less, m, m-s, m+s)
		medianOfThree{{.Plural}}({{.Receiver}}, less, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThree{{.Plural}}({{.Receiver}}, less, lo, m, hi-1)

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
			if less({{.Receiver}}, b, pivot) { // {{.Receiver}}[b] < pivot
				b++
			} else if !less({{.Receiver}}, pivot, b) { // {{.Receiver}}[b] = pivot
				swap{{.Plural}}({{.Receiver}}, a, b)
				a++
				b++
			} else {
				break
			}
		}
		for b < c {
			if less({{.Receiver}}, pivot, c-1) { // {{.Receiver}}[c-1] > pivot
				c--
			} else if !less({{.Receiver}}, c-1, pivot) { // {{.Receiver}}[c-1] = pivot
				swap{{.Plural}}({{.Receiver}}, c-1, d-1)
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
		swap{{.Plural}}({{.Receiver}}, b, c-1)
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
	swapRange{{.Plural}}({{.Receiver}}, lo, b-n, n)

	n = min(hi-d, d-c)
	swapRange{{.Plural}}({{.Receiver}}, c, hi-n, n)

	return lo + b - a, hi - (d - c)
}

func quickSort{{.Plural}}({{.Receiver}} {{.Plural}}, less func({{.Plural}}, int, int) bool, a, b, maxDepth int) {
	for b-a > 7 {
		if maxDepth == 0 {
			heapSort{{.Plural}}({{.Receiver}}, less, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivot{{.Plural}}({{.Receiver}}, less, a, b)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			quickSort{{.Plural}}({{.Receiver}}, less, a, mlo, maxDepth)
			a = mhi // i.e., quickSort{{.Plural}}({{.Receiver}}, mhi, b)
		} else {
			quickSort{{.Plural}}({{.Receiver}}, less, mhi, b, maxDepth)
			b = mlo // i.e., quickSort{{.Plural}}({{.Receiver}}, a, mlo)
		}
	}
	if b-a > 1 {
		insertionSort{{.Plural}}({{.Receiver}}, less, a, b)
	}
}

func not(less func({{.Plural}}, int, int) bool) func({{.Plural}}, int, int) bool {
	return func(z {{.Plural}}, a int, b int) bool {
		return !less(z, a, b)
	}
}
`
