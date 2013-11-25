package main

import (
	"errors"
	"fmt"
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

// Tests that all elements of {{.Plural}} are true for the passed func. Example:
//	good := func({{.Loop}} {{.Pointer}}{{.Singular}}) bool {
//		return {{.Loop}}.Something > 42
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
//	winner := func({{.Loop}} {{.Pointer}}{{.Singular}}) bool {
//		return {{.Loop}}.Placement == "winner"
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
//	dracula := func({{.Loop}} {{.Pointer}}{{.Singular}}) bool {
//		return {{.Loop}}.IsDracula()
//	}
//	countDracula := my{{.Plural}}.Count(dracula)
func ({{.Receiver}} {{.Plural}}) Count(fn func({{.Pointer}}{{.Singular}}) bool) (result int) {
	for _, {{.Loop}} := range {{.Receiver}} {
		if fn({{.Loop}}) {
			result++
		}
	}
	return
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
//	winner := func({{.Loop}} {{.Pointer}}{{.Singular}}) bool {
//		return {{.Loop}}.Placement == "winner"
//	}
//	theWinner, err := myMovies.First(winner)
func ({{.Receiver}} {{.Plural}}) First(fn func({{.Pointer}}{{.Singular}}) bool) (result {{.Pointer}}{{.Singular}}, err error) {
	for _, {{.Loop}} := range {{.Receiver}} {
		if fn({{.Loop}}) {
			result = {{.Loop}}
			return
		}
	}
	err = errors.New("No {{.Plural}} elements return true for passed func")
	return
}

// Returns the element of {{.Plural}} containing the maximum value, when compared to other elements using a passed func defining ‘less’. Example:
//	byArea := func({{.Loop}}s {{.Plural}}, a int, b int) bool {
//		return {{.Loop}}s[a].Area() < {{.Loop}}s[b].Area()
//	}
//	roomiest := my{{.Plural}}.Max(byArea)
//
// In the case of multiple items being equally maximal, the last such element is returned.
// (Note: this is implemented by negating the passed ‘less’ func, effectively testing ‘greater than or equal to’.)
func ({{.Receiver}} {{.Plural}}) Max(less func({{.Plural}}, int, int) bool) (result {{.Pointer}}{{.Singular}}, err error) {
	if len(rcv) == 0 {
		err = errors.New("Cannot determine the Max of an empty slice")
		return
	}
	return rcv.Min(negate{{.Plural}}(less))
}

// Returns the element of {{.Plural}} containing the minimum value, when compared to other elements using a passed func defining ‘less’. Example:
//	byPrice := func({{.Loop}}s {{.Plural}}, a int, b int) bool {
//		return {{.Loop}}s[a].Price < {{.Loop}}s[b].Price
//	}
//	cheapest := my{{.Plural}}.Min(byPrice)
//
// In the case of multiple items being equally minimal, the first such element is returned.
func ({{.Receiver}} {{.Plural}}) Min(less func({{.Plural}}, int, int) bool) (result {{.Pointer}}{{.Singular}}, err error) {
	l := len({{.Receiver}})
	if l == 0 {
		err = errors.New("Cannot determine the Min of an empty slice")
		return
	}
	m := 0
	for i := 1; i < l; i++ {
		if less({{.Receiver}}, i, m) {
			m = i
		}
	}
	result = {{.Receiver}}[m]
	return
}

// Returns exactly one element that returns true for the passed func. Returns errors if no or multiple elements return true. Example:
//	byId := func({{.Loop}} {{.Pointer}}{{.Singular}}) bool {
//		return {{.Loop}}.Id == 5
//	}
//	single, err := myMovies.Single(byId)
func ({{.Receiver}} {{.Plural}}) Single(fn func({{.Pointer}}{{.Singular}}) bool) (result {{.Pointer}}{{.Singular}}, err error) {
	var candidate {{.Pointer}}{{.Singular}}
	found := false
	for _, {{.Loop}} := range {{.Receiver}} {
		if fn({{.Loop}}) {
			if found {
				err = errors.New("Multiple {{.Plural}} elements return true for passed func")
				return
			}
			candidate = {{.Loop}}
			found = true
		}
	}
	if found {
		result = candidate
	} else {
		err = errors.New("No {{.Plural}} elements return true for passed func")
	}
	return
}

// Returns a new {{.Plural}} slice whose elements return true for func. Example:
//	incredible := func({{.Loop}} {{.Pointer}}{{.Singular}}) bool {
//		return {{.Loop}}.Manufacturer == "Apple"
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
//	byName := func({{.Loop}}s {{.Plural}}, a int, b int) bool {
//		return {{.Loop}}s[a].LastName < {{.Loop}}s[b].LastName
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
//	byPoints := func({{.Loop}}s {{.Plural}}, a int, b int) bool {
//		return {{.Loop}}s[a].Points < {{.Loop}}s[b].Points
//	}
//	leaderboard := my{{.Plural}}.SortDesc(byPoints)
// (Note: this is implemented by negating the passed ‘less’ func, effectively testing ‘greater than or equal to’.)
func ({{.Receiver}} {{.Plural}}) SortDesc(less func({{.Plural}}, int, int) bool) {{.Plural}} {
	return {{.Receiver}}.Sort(negate{{.Plural}}(less))
}

// Reports whether an instance of {{.Plural}} is sorted in descending order, using the pass func to define ‘less’. See SortDesc method below.
func ({{.Receiver}} {{.Plural}}) IsSortedDesc(less func({{.Plural}}, int, int) bool) bool {
	return {{.Receiver}}.IsSorted(negate{{.Plural}}(less))
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

func negate{{.Plural}}(less func({{.Plural}}, int, int) bool) func({{.Plural}}, int, int) bool {
	return func(z {{.Plural}}, a int, b int) bool {
		return !less(z, a, b)
	}
}
`

func getCustomTemplate(name string) (result *template.Template, err error) {
	t, found := customTemplates[name]
	if found {
		result = template.Must(template.New(name).Parse(t))
	} else {
		err = errors.New(fmt.Sprintf("%s is not a known custom method", name))
	}
	return
}

var customTemplates = map[string]string{
	"Aggregate": `
func ({{.Parent.Receiver}} {{.Parent.Plural}}) Aggregate{{.Name}}(fn func({{.Pointer}}{{.Package}}{{.Type}}, {{.Pointer}}{{.Package}}{{.Type}}) {{.Pointer}}{{.Package}}{{.Type}}) (result {{.Pointer}}{{.Package}}{{.Type}}) {
	for _, {{.Parent.Loop}} := range {{.Parent.Receiver}} {
		result = fn(result, {{.Parent.Loop}}.{{.Name}})
	}
	return
}
`,
	"DistinctBy": `
func ({{.Parent.Receiver}} {{.Parent.Plural}}) DistinctBy{{.Name}}() {{.Parent.Plural}} {
	equal := func(a {{.Parent.Pointer}}{{.Parent.Singular}}, b {{.Parent.Pointer}}{{.Parent.Singular}}) bool {
		return a.{{.Name}} == b.{{.Name}}
	}
	return {{.Parent.Receiver}}.DistinctBy(equal)
}
`,
	"GroupBy": `
func ({{.Parent.Receiver}} {{.Parent.Plural}}) GroupBy{{.Name}}() map[{{.Pointer}}{{.Package}}{{.Type}}]{{.Parent.Plural}} {
	result := make(map[{{.Pointer}}{{.Package}}{{.Type}}]{{.Parent.Plural}})
	for _, {{.Parent.Loop}} := range {{.Parent.Receiver}} {
		result[{{.Parent.Loop}}.{{.Name}}] = append(result[{{.Parent.Loop}}.{{.Name}}], {{.Parent.Loop}})
	}
	return result
}
`,
	"Select": `
func ({{.Parent.Receiver}} {{.Parent.Plural}}) Select{{.Name}}() (result []{{.Pointer}}{{.Package}}{{.Type}}) {
	for _, {{.Parent.Loop}} := range {{.Parent.Receiver}} {
		result = append(result, {{.Parent.Loop}}.{{.Name}})
	}
	return
}
`,
	"SortBy": `
func ({{.Parent.Receiver}} {{.Parent.Plural}}) SortBy{{.Name}}() {{.Parent.Plural}} {
	less := func(z {{.Parent.Plural}}, a int, b int) bool {
		return z[a].{{.Name}} < z[b].{{.Name}}
	}
	return {{.Parent.Receiver}}.Sort(less)
}
`,
	"Sum": `
func ({{.Parent.Receiver}} {{.Parent.Plural}}) Sum{{.Name}}() (result {{.Pointer}}{{.Package}}{{.Type}}) {
	for _, {{.Parent.Loop}} := range {{.Parent.Receiver}} {
		result += {{.Parent.Loop}}.{{.Name}}
	}
	return
}
`,
}
