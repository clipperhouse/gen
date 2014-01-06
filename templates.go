package main

import (
	"errors"
	"fmt"
	"text/template"
)

func getHeaderTemplate() *template.Template {
	return template.Must(template.New("header").Parse(header))
}

const header = `// {{.Command}}
// this file was auto-generated using github.com/clipperhouse/gen
// {{.Generated}}

// Sort functions are a modification of http://golang.org/pkg/sort/#Sort
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package {{.Package}}
{{if gt (len .Imports) 0}}
import ({{range .Imports}}
	"{{.}}"
{{end}})
{{end}}
// The plural (slice) type of {{.Pointer}}{{.Singular}}, for use with gen methods below. Use this type where you would use []{{.Pointer}}{{.Singular}}. (This is required because slices cannot be method receivers.)
type {{.Plural}} []{{.Pointer}}{{.Singular}}
`

func getStandardMethodKeys() (result []string) {
	for k := range standardTemplates {
		result = append(result, k)
	}
	return
}

func getStandardTemplate(name string) (result *template.Template, err error) {
	t, found := standardTemplates[name]
	if found {
		result = template.Must(template.New(name).Parse(t))
	} else {
		err = errors.New(fmt.Sprintf("%s is not a known method", name))
	}
	return
}

var standardTemplates = map[string]string{
	"Equals": `
// Tests that all elements of one {{.Plural}} match the corresponding elements of another {{.Plural}}. See: http://clipperhouse.github.io/gen/#Equals
func ({{.Receiver}} {{.Plural}}) Equals(other {{.Plural}}) bool {
	if len({{.Receiver}}) != len(other) {
		return false
	}
	for i, v := range rcv {
		if v != other[i] {
			return false
		}
	}
	return true
}
`,
	"All": `
// Tests that all elements of {{.Plural}} return true for the passed func. See: http://clipperhouse.github.io/gen/#All
func ({{.Receiver}} {{.Plural}}) All(fn func({{.Pointer}}{{.Singular}}) bool) bool {
	for _, {{.Loop}} := range {{.Receiver}} {
		if !fn({{.Loop}}) {
			return false
		}
	}
	return true
}
`,
	"Any": `
// Tests that one or more elements of {{.Plural}} return true for the passed func. See: http://clipperhouse.github.io/gen/#Any
func ({{.Receiver}} {{.Plural}}) Any(fn func({{.Pointer}}{{.Singular}}) bool) bool {
	for _, {{.Loop}} := range {{.Receiver}} {
		if fn({{.Loop}}) {
			return true
		}
	}
	return false
}
`,
	"Count": `
// Counts the number elements of {{.Plural}} that return true for the passed func. See: http://clipperhouse.github.io/gen/#Count
func ({{.Receiver}} {{.Plural}}) Count(fn func({{.Pointer}}{{.Singular}}) bool) (result int) {
	for _, {{.Loop}} := range {{.Receiver}} {
		if fn({{.Loop}}) {
			result++
		}
	}
	return
}
`,
	"Distinct": `
// Returns a new {{.Plural}} slice whose elements are unique. See: http://clipperhouse.github.io/gen/#Distinct
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
`,
	"DistinctBy": `
// Returns a new {{.Plural}} slice whose elements are unique, where equality is defined by a passed func. See: http://clipperhouse.github.io/gen/#DistinctBy
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
`,
	"Each": `
// Iterates over {{.Plural}} and executes the passed func against each element. See: http://clipperhouse.github.io/gen/#Each
func ({{.Receiver}} {{.Plural}}) Each(fn func({{.Pointer}}{{.Singular}})) {
	for _, {{.Loop}} := range {{.Receiver}} {
		fn({{.Loop}})
	}
}
`,
	"First": `
// Returns the first element that returns true for the passed func. Returns error if no elements return true. See: http://clipperhouse.github.io/gen/#First
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
`,
	"Max": `
// Returns an element of {{.Plural}} containing the maximum value, when compared to other elements using a passed func defining ‘less’. In the case of multiple items being equally maximal, the last such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#Max
//
// (Note: this is implemented by negating the passed ‘less’ func, effectively testing ‘greater than or equal to’.)
func ({{.Receiver}} {{.Plural}}) Max(less func({{.Pointer}}{{.Singular}}, {{.Pointer}}{{.Singular}}) bool) (result {{.Pointer}}{{.Singular}}, err error) {
	l := len({{.Receiver}})
	if l == 0 {
		err = errors.New("Cannot determine the Max of an empty slice")
		return
	}
	m := 0
	for i := 1; i < l; i++ {
		if !less({{.Receiver}}[i], {{.Receiver}}[m]) {
			m = i
		}
	}
	result = {{.Receiver}}[m]
	return
}
`,
	"Min": `
// Returns an element of {{.Plural}} containing the minimum value, when compared to other elements using a passed func defining ‘less’. In the case of multiple items being equally minimal, the first such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#Min
func ({{.Receiver}} {{.Plural}}) Min(less func({{.Pointer}}{{.Singular}}, {{.Pointer}}{{.Singular}}) bool) (result {{.Pointer}}{{.Singular}}, err error) {
	l := len({{.Receiver}})
	if l == 0 {
		err = errors.New("Cannot determine the Min of an empty slice")
		return
	}
	m := 0
	for i := 1; i < l; i++ {
		if less({{.Receiver}}[i], {{.Receiver}}[m]) {
			m = i
		}
	}
	result = {{.Receiver}}[m]
	return
}
`,
	"Single": `
// Returns exactly one element of {{.Plural}} that returns true for the passed func. Returns error if no or multiple elements return true. See: http://clipperhouse.github.io/gen/#Single
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
`,
	"Where": `
// Returns a new {{.Plural}} slice whose elements return true for func. See: http://clipperhouse.github.io/gen/#Where
func ({{.Receiver}} {{.Plural}}) Where(fn func({{.Pointer}}{{.Singular}}) bool) (result {{.Plural}}) {
	for _, {{.Loop}} := range {{.Receiver}} {
		if fn({{.Loop}}) {
			result = append(result, {{.Loop}})
		}
	}
	return result
}
`,
	"Sort": `
// Returns a new ordered {{.Plural}} slice, determined by a func defining ‘less’. See: http://clipperhouse.github.io/gen/#Sort
func ({{.Receiver}} {{.Plural}}) Sort(less func({{.Pointer}}{{.Singular}}, {{.Pointer}}{{.Singular}}) bool) {{.Plural}} {
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
`,
	"IsSorted": `
// Reports whether an instance of {{.Plural}} is sorted, using the pass func to define ‘less’. See: http://clipperhouse.github.io/gen/#Sort
func ({{.Receiver}} {{.Plural}}) IsSorted(less func({{.Pointer}}{{.Singular}}, {{.Pointer}}{{.Singular}}) bool) bool {
	n := len({{.Receiver}})
	for i := n - 1; i > 0; i-- {
		if less({{.Receiver}}[i], {{.Receiver}}[i-1]) {
			return false
		}
	}
	return true
}
`,
	"SortDesc": `
// Returns a new, descending-ordered {{.Plural}} slice, determined by a func defining ‘less’. See: http://clipperhouse.github.io/gen/#Sort
//
// (Note: this is implemented by negating the passed ‘less’ func, effectively testing ‘greater than or equal to’.)
func ({{.Receiver}} {{.Plural}}) SortDesc(less func({{.Pointer}}{{.Singular}}, {{.Pointer}}{{.Singular}}) bool) {{.Plural}} {
	greaterOrEqual := func(a, b {{.Pointer}}{{.Singular}}) bool {
		return !less(a, b)
	}
	return {{.Receiver}}.Sort(greaterOrEqual)
}
`,
	"IsSortedDesc": `
// Reports whether an instance of {{.Plural}} is sorted in descending order, using the pass func to define ‘less’. See: http://clipperhouse.github.io/gen/#Sort
func ({{.Receiver}} {{.Plural}}) IsSortedDesc(less func({{.Pointer}}{{.Singular}}, {{.Pointer}}{{.Singular}}) bool) bool {
	greaterOrEqual := func(a, b {{.Pointer}}{{.Singular}}) bool {
		return !less(a, b)
	}
	return {{.Receiver}}.IsSorted(greaterOrEqual)
}
`,
}

func getSortSupportTemplate() *template.Template {
	return template.Must(template.New("sortSupport").Parse(sortSupport))
}

const sortSupport = `
// Sort support methods

func swap{{.Plural}}({{.Receiver}} {{.Plural}}, a, b int) {
	{{.Receiver}}[a], {{.Receiver}}[b] = {{.Receiver}}[b], {{.Receiver}}[a]
}

// Insertion sort
func insertionSort{{.Plural}}({{.Receiver}} {{.Plural}}, less func({{.Pointer}}{{.Singular}}, {{.Pointer}}{{.Singular}}) bool, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && less({{.Receiver}}[j], {{.Receiver}}[j-1]); j-- {
			swap{{.Plural}}({{.Receiver}}, j, j-1)
		}
	}
}

// siftDown implements the heap property on {{.Receiver}}[lo, hi).
// first is an offset into the array where the root of the heap lies.
func siftDown{{.Plural}}({{.Receiver}} {{.Plural}}, less func({{.Pointer}}{{.Singular}}, {{.Pointer}}{{.Singular}}) bool, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && less({{.Receiver}}[first+child], {{.Receiver}}[first+child+1]) {
			child++
		}
		if !less({{.Receiver}}[first+root], {{.Receiver}}[first+child]) {
			return
		}
		swap{{.Plural}}({{.Receiver}}, first+root, first+child)
		root = child
	}
}

func heapSort{{.Plural}}({{.Receiver}} {{.Plural}}, less func({{.Pointer}}{{.Singular}}, {{.Pointer}}{{.Singular}}) bool, a, b int) {
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
func medianOfThree{{.Plural}}({{.Receiver}} {{.Plural}}, less func({{.Pointer}}{{.Singular}}, {{.Pointer}}{{.Singular}}) bool, a, b, c int) {
	m0 := b
	m1 := a
	m2 := c
	// bubble sort on 3 elements
	if less({{.Receiver}}[m1], {{.Receiver}}[m0]) {
		swap{{.Plural}}({{.Receiver}}, m1, m0)
	}
	if less({{.Receiver}}[m2], {{.Receiver}}[m1]) {
		swap{{.Plural}}({{.Receiver}}, m2, m1)
	}
	if less({{.Receiver}}[m1], {{.Receiver}}[m0]) {
		swap{{.Plural}}({{.Receiver}}, m1, m0)
	}
	// now {{.Receiver}}[m0] <= {{.Receiver}}[m1] <= {{.Receiver}}[m2]
}

func swapRange{{.Plural}}({{.Receiver}} {{.Plural}}, a, b, n int) {
	for i := 0; i < n; i++ {
		swap{{.Plural}}({{.Receiver}}, a+i, b+i)
	}
}

func doPivot{{.Plural}}({{.Receiver}} {{.Plural}}, less func({{.Pointer}}{{.Singular}}, {{.Pointer}}{{.Singular}}) bool, lo, hi int) (midlo, midhi int) {
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
			if less({{.Receiver}}[b], {{.Receiver}}[pivot]) { // {{.Receiver}}[b] < pivot
				b++
			} else if !less({{.Receiver}}[pivot], {{.Receiver}}[b]) { // {{.Receiver}}[b] = pivot
				swap{{.Plural}}({{.Receiver}}, a, b)
				a++
				b++
			} else {
				break
			}
		}
		for b < c {
			if less({{.Receiver}}[pivot], {{.Receiver}}[c-1]) { // {{.Receiver}}[c-1] > pivot
				c--
			} else if !less({{.Receiver}}[c-1], {{.Receiver}}[pivot]) { // {{.Receiver}}[c-1] = pivot
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

func quickSort{{.Plural}}({{.Receiver}} {{.Plural}}, less func({{.Pointer}}{{.Singular}}, {{.Pointer}}{{.Singular}}) bool, a, b, maxDepth int) {
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
// Iterates over {{.Parent.Plural}}, operating on each element while maintaining ‘state’. See: http://clipperhouse.github.io/gen/#Aggregate
func ({{.Parent.Receiver}} {{.Parent.Plural}}) Aggregate{{.Name}}(fn func({{.Type}}, {{.Type}}) {{.Type}}) (result {{.Type}}) {
	for _, {{.Parent.Loop}} := range {{.Parent.Receiver}} {
		result = fn(result, {{.Parent.Loop}}.{{.Name}})
	}
	return
}
`,
	"Average": `
// Sums {{.Name}} over all elements and divides by len({{.Parent.Plural}}). See: http://clipperhouse.github.io/gen/#Average
func ({{.Parent.Receiver}} {{.Parent.Plural}}) Average{{.Name}}() (result {{.Type}}, err error) {
	l := len({{.Parent.Receiver}})
	if l == 0 {
		err = errors.New("cannot determine Average{{.Name}} of zero-length {{.Parent.Plural}}")
		return
	}
	for _, {{.Parent.Loop}} := range {{.Parent.Receiver}} {
		result += {{.Parent.Loop}}.{{.Name}}
	}
	result = result / l
	return
}
`,
	"GroupBy": `
// Groups elements into a map keyed by {{.Name}}’s value. See: http://clipperhouse.github.io/gen/#GroupBy
func ({{.Parent.Receiver}} {{.Parent.Plural}}) GroupBy{{.Name}}() map[{{.Type}}]{{.Parent.Plural}} {
	result := make(map[{{.Type}}]{{.Parent.Plural}})
	for _, {{.Parent.Loop}} := range {{.Parent.Receiver}} {
		result[{{.Parent.Loop}}.{{.Name}}] = append(result[{{.Parent.Loop}}.{{.Name}}], {{.Parent.Loop}})
	}
	return result
}
`,
	"Max": `
// Selects the largest value of {{.Name}} in {{.Parent.Plural}}. Returns error on {{.Parent.Plural}} with no elements. See: http://clipperhouse.github.io/gen/#MaxCustom
func ({{.Parent.Receiver}} {{.Parent.Plural}}) Max{{.Name}}() (result {{.Type}}, err error) {
	l := len({{.Parent.Receiver}})
	if l == 0 {
		err = errors.New("cannot determine Max{{.Name}} of zero-length {{.Parent.Plural}}")
		return
	}
	result = {{.Parent.Receiver}}[0].{{.Name}}
	if l > 1 {
		for _, {{.Parent.Loop}} := range {{.Parent.Receiver}}[1:] {
			if {{.Parent.Loop}}.{{.Name}} > result {
				result = {{.Parent.Loop}}.{{.Name}}
			}
		}
	}
	return
}
`,
	"Min": `
// Selects the least value of {{.Name}} in {{.Parent.Plural}}. Returns error on {{.Parent.Plural}} with no elements. See: http://clipperhouse.github.io/gen/#MinCustom
func ({{.Parent.Receiver}} {{.Parent.Plural}}) Min{{.Name}}() (result {{.Type}}, err error) {
	l := len({{.Parent.Receiver}})
	if l == 0 {
		err = errors.New("cannot determine Min{{.Name}} of zero-length {{.Parent.Plural}}")
		return
	}
	result = {{.Parent.Receiver}}[0].{{.Name}}
	if l > 1 {
		for _, {{.Parent.Loop}} := range {{.Parent.Receiver}}[1:] {
			if {{.Parent.Loop}}.{{.Name}} < result {
				result = {{.Parent.Loop}}.{{.Name}}
			}
		}
	}
	return
}
`,
	"Select": `
// Returns a slice containing all values of {{.Name}} in {{.Parent.Plural}}. See: http://clipperhouse.github.io/gen/#Select
func ({{.Parent.Receiver}} {{.Parent.Plural}}) Select{{.Name}}() (result []{{.Type}}) {
	for _, {{.Parent.Loop}} := range {{.Parent.Receiver}} {
		result = append(result, {{.Parent.Loop}}.{{.Name}})
	}
	return
}
`,
	"Sum": `
// Sums {{.Name}} over all elements in {{.Parent.Plural}}. See: http://clipperhouse.github.io/gen/#Sum
func ({{.Parent.Receiver}} {{.Parent.Plural}}) Sum{{.Name}}() (result {{.Type}}) {
	for _, {{.Parent.Loop}} := range {{.Parent.Receiver}} {
		result += {{.Parent.Loop}}.{{.Name}}
	}
	return
}
`,
}
