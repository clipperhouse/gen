package main

import (
	"fmt"
	"sort"
	"text/template"
)

func getStandardTemplate(name string) (result *template.Template, err error) {
	t, found := StandardTemplates[name]
	if found {
		result = template.Must(template.New(name).Parse(t.Text))
	} else {
		err = fmt.Errorf("%s is not a known method", name)
	}
	return
}

func isStandardMethod(s string) bool {
	_, ok := StandardTemplates[s]
	return ok
}

func getStandardMethodKeys() (result []string) {
	for k := range StandardTemplates {
		result = append(result, k)
	}
	sort.Strings(result)
	return
}

var StandardTemplates = map[string]*Template{

	"All": &Template{
		Text: `
// All verifies that all elements of {{.Plural}} return true for the passed func. See: http://clipperhouse.github.io/gen/#All
func (rcv {{.Plural}}) All(fn func({{.Pointer}}{{.Name}}) bool) bool {
	for _, v := range rcv {
		if !fn(v) {
			return false
		}
	}
	return true
}
`},

	"Any": &Template{
		Text: `
// Any verifies that one or more elements of {{.Plural}} return true for the passed func. See: http://clipperhouse.github.io/gen/#Any
func (rcv {{.Plural}}) Any(fn func({{.Pointer}}{{.Name}}) bool) bool {
	for _, v := range rcv {
		if fn(v) {
			return true
		}
	}
	return false
}
`},

	"Count": &Template{
		Text: `
// Count gives the number elements of {{.Plural}} that return true for the passed func. See: http://clipperhouse.github.io/gen/#Count
func (rcv {{.Plural}}) Count(fn func({{.Pointer}}{{.Name}}) bool) (result int) {
	for _, v := range rcv {
		if fn(v) {
			result++
		}
	}
	return
}
`},

	"Distinct": &Template{
		Text: `
// Distinct returns a new {{.Plural}} slice whose elements are unique. See: http://clipperhouse.github.io/gen/#Distinct
func (rcv {{.Plural}}) Distinct() (result {{.Plural}}) {
	appended := make(map[{{.Pointer}}{{.Name}}]bool)
	for _, v := range rcv {
		if !appended[v] {
			result = append(result, v)
			appended[v] = true
		}
	}
	return result
}
`,
		RequiresComparable: true,
	},

	"DistinctBy": &Template{
		Text: `
// DistinctBy returns a new {{.Plural}} slice whose elements are unique, where equality is defined by a passed func. See: http://clipperhouse.github.io/gen/#DistinctBy
func (rcv {{.Plural}}) DistinctBy(equal func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool) (result {{.Plural}}) {
	for _, v := range rcv {
		eq := func(_app {{.Pointer}}{{.Name}}) bool {
			return equal(v, _app)
		}
		if !result.Any(eq) {
			result = append(result, v)
		}
	}
	return result
}
`},

	"Each": &Template{
		Text: `
// Each iterates over {{.Plural}} and executes the passed func against each element. See: http://clipperhouse.github.io/gen/#Each
func (rcv {{.Plural}}) Each(fn func({{.Pointer}}{{.Name}})) {
	for _, v := range rcv {
		fn(v)
	}
}
`},

	"First": &Template{
		Text: `
// First returns the first element that returns true for the passed func. Returns error if no elements return true. See: http://clipperhouse.github.io/gen/#First
func (rcv {{.Plural}}) First(fn func({{.Pointer}}{{.Name}}) bool) (result {{.Pointer}}{{.Name}}, err error) {
	for _, v := range rcv {
		if fn(v) {
			result = v
			return
		}
	}
	err = errors.New("no {{.Plural}} elements return true for passed func")
	return
}
`},

	"Max": &Template{
		Text: `
// Max returns the maximum value of {{.Plural}}. In the case of multiple items being equally maximal, the first such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#Max
func (rcv {{.Plural}}) Max() (result {{.Pointer}}{{.Name}}, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine the Max of an empty slice")
		return
	}
	result = rcv[0]
	for _, v := range rcv {
		if v > result {
			result = v
		}
	}
	return
}
`,
		RequiresOrdered: true,
	},

	"Min": &Template{
		Text: `
// Min returns the minimum value of {{.Plural}}. In the case of multiple items being equally minimal, the first such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#Min
func (rcv {{.Plural}}) Min() (result {{.Pointer}}{{.Name}}, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine the Min of an empty slice")
		return
	}
	result = rcv[0]
	for _, v := range rcv {
		if v < result {
			result = v
		}
	}
	return
}
`,
		RequiresOrdered: true,
	},

	"MaxBy": &Template{
		Text: `
// MaxBy returns an element of {{.Plural}} containing the maximum value, when compared to other elements using a passed func defining ‘less’. In the case of multiple items being equally maximal, the last such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#MaxBy
func (rcv {{.Plural}}) MaxBy(less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool) (result {{.Pointer}}{{.Name}}, err error) {
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
`},

	"MinBy": &Template{
		Text: `
// MinBy returns an element of {{.Plural}} containing the minimum value, when compared to other elements using a passed func defining ‘less’. In the case of multiple items being equally minimal, the first such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#MinBy
func (rcv {{.Plural}}) MinBy(less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool) (result {{.Pointer}}{{.Name}}, err error) {
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
`},

	"Single": &Template{
		Text: `
// Single returns exactly one element of {{.Plural}} that returns true for the passed func. Returns error if no or multiple elements return true. See: http://clipperhouse.github.io/gen/#Single
func (rcv {{.Plural}}) Single(fn func({{.Pointer}}{{.Name}}) bool) (result {{.Pointer}}{{.Name}}, err error) {
	var candidate {{.Pointer}}{{.Name}}
	found := false
	for _, v := range rcv {
		if fn(v) {
			if found {
				err = errors.New("multiple {{.Plural}} elements return true for passed func")
				return
			}
			candidate = v
			found = true
		}
	}
	if found {
		result = candidate
	} else {
		err = errors.New("no {{.Plural}} elements return true for passed func")
	}
	return
}
`},

	"Where": &Template{
		Text: `
// Where returns a new {{.Plural}} slice whose elements return true for func. See: http://clipperhouse.github.io/gen/#Where
func (rcv {{.Plural}}) Where(fn func({{.Pointer}}{{.Name}}) bool) (result {{.Plural}}) {
	for _, v := range rcv {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
`},

	"Sort": &Template{
		Text: `
// Sort returns a new ordered {{.Plural}} slice. See: http://clipperhouse.github.io/gen/#Sort
func (rcv {{.Plural}}) Sort() {{.Plural}} {
	result := make({{.Plural}}, len(rcv))
	copy(result, rcv)
	sort.Sort(result)
	return result
}
`,
		RequiresOrdered: true,
	},
	"IsSorted": &Template{
		Text: `
// IsSorted reports whether {{.Plural}} is sorted. See: http://clipperhouse.github.io/gen/#Sort
func (rcv {{.Plural}}) IsSorted() bool {
	return sort.IsSorted(rcv)
}
`,
		RequiresOrdered: true,
	},
	"SortDesc": &Template{
		Text: `
// SortDesc returns a new reverse-ordered {{.Plural}} slice. See: http://clipperhouse.github.io/gen/#Sort
func (rcv {{.Plural}}) SortDesc() {{.Plural}} {
	result := make({{.Plural}}, len(rcv))
	copy(result, rcv)
	sort.Sort(sort.Reverse(result))
	return result
}
`,
		RequiresOrdered: true,
	},
	"IsSortedDesc": &Template{
		Text: `
// IsSortedDesc reports whether {{.Plural}} is reverse-sorted. See: http://clipperhouse.github.io/gen/#Sort
func (rcv {{.Plural}}) IsSortedDesc() bool {
	return sort.IsSorted(sort.Reverse(rcv))
}
`,
		RequiresOrdered: true,
	},

	"SortBy": &Template{
		Text: `
// SortBy returns a new ordered {{.Plural}} slice, determined by a func defining ‘less’. See: http://clipperhouse.github.io/gen/#SortBy
func (rcv {{.Plural}}) SortBy(less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool) {{.Plural}} {
	result := make({{.Plural}}, len(rcv))
	copy(result, rcv)
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
`},

	"IsSortedBy": &Template{
		Text: `
// IsSortedBy reports whether an instance of {{.Plural}} is sorted, using the pass func to define ‘less’. See: http://clipperhouse.github.io/gen/#SortBy
func (rcv {{.Plural}}) IsSortedBy(less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool) bool {
	n := len(rcv)
	for i := n - 1; i > 0; i-- {
		if less(rcv[i], rcv[i-1]) {
			return false
		}
	}
	return true
}
`},

	"SortByDesc": &Template{
		Text: `
// SortByDesc returns a new, descending-ordered {{.Plural}} slice, determined by a func defining ‘less’. See: http://clipperhouse.github.io/gen/#SortBy
func (rcv {{.Plural}}) SortByDesc(less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool) {{.Plural}} {
	greater := func(a, b {{.Pointer}}{{.Name}}) bool {
		return a != b && !less(a, b)
	}
	return rcv.SortBy(greater)
}
`},

	"IsSortedByDesc": &Template{
		Text: `
// IsSortedDesc reports whether an instance of {{.Plural}} is sorted in descending order, using the pass func to define ‘less’. See: http://clipperhouse.github.io/gen/#SortBy
func (rcv {{.Plural}}) IsSortedByDesc(less func({{.Pointer}}{{.Name}}, {{.Pointer}}{{.Name}}) bool) bool {
	greater := func(a, b {{.Pointer}}{{.Name}}) bool {
		return a != b && !less(a, b)
	}
	return rcv.IsSortedBy(greater)
}
`},
}
