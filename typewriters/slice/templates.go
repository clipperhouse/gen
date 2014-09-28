package slice

import (
	"github.com/clipperhouse/gen/typewriter"
)

// a convenience for passing values into templates; in MVC it'd be called a view model
type model struct {
	Type      typewriter.Type
	SliceName string
	// these tempaltes only ever happen to use one type parameter
	TypeParameter typewriter.Type
	typewriter.TagValue
}

var templates = typewriter.TemplateSet{

	"slice": &typewriter.Template{
		Text: `// {{.SliceName}} is a slice of type {{.Type}}. Use it where you would use []{{.Type}}.
type {{.SliceName}} []{{.Type}}
`,
	},

	"All": &typewriter.Template{
		Text: `
// All verifies that all elements of {{.SliceName}} return true for the passed func. See: http://clipperhouse.github.io/gen/#All
func (rcv {{.SliceName}}) All(fn func({{.Type}}) bool) bool {
	for _, v := range rcv {
		if !fn(v) {
			return false
		}
	}
	return true
}
`},

	"Any": &typewriter.Template{
		Text: `
// Any verifies that one or more elements of {{.SliceName}} return true for the passed func. See: http://clipperhouse.github.io/gen/#Any
func (rcv {{.SliceName}}) Any(fn func({{.Type}}) bool) bool {
	for _, v := range rcv {
		if fn(v) {
			return true
		}
	}
	return false
}
`},

	"Average": &typewriter.Template{
		Text: `
// Average sums {{.SliceName}} over all elements and divides by len({{.SliceName}}). See: http://clipperhouse.github.io/gen/#Average
func (rcv {{.SliceName}}) Average() (result {{.Type}}, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine Average of zero-length {{.SliceName}}")
		return
	}
	for _, v := range rcv {
		result += v
	}
	result = result / {{.Type}}(l)
	return
}
`,
		TypeConstraint: typewriter.Constraint{Numeric: true},
	},

	"Count": &typewriter.Template{
		Text: `
// Count gives the number elements of {{.SliceName}} that return true for the passed func. See: http://clipperhouse.github.io/gen/#Count
func (rcv {{.SliceName}}) Count(fn func({{.Type}}) bool) (result int) {
	for _, v := range rcv {
		if fn(v) {
			result++
		}
	}
	return
}
`},

	"Distinct": &typewriter.Template{
		Text: `
// Distinct returns a new {{.SliceName}} whose elements are unique. See: http://clipperhouse.github.io/gen/#Distinct
func (rcv {{.SliceName}}) Distinct() (result {{.SliceName}}) {
	appended := make(map[{{.Type}}]bool)
	for _, v := range rcv {
		if !appended[v] {
			result = append(result, v)
			appended[v] = true
		}
	}
	return result
}
`,
		TypeConstraint: typewriter.Constraint{Comparable: true},
	},

	"DistinctBy": &typewriter.Template{
		Text: `
// DistinctBy returns a new {{.SliceName}} whose elements are unique, where equality is defined by a passed func. See: http://clipperhouse.github.io/gen/#DistinctBy
func (rcv {{.SliceName}}) DistinctBy(equal func({{.Type}}, {{.Type}}) bool) (result {{.SliceName}}) {
	for _, v := range rcv {
		eq := func(_app {{.Type}}) bool {
			return equal(v, _app)
		}
		if !result.Any(eq) {
			result = append(result, v)
		}
	}
	return result
}
`},

	"Each": &typewriter.Template{
		Text: `
// Each iterates over {{.SliceName}} and executes the passed func against each element. See: http://clipperhouse.github.io/gen/#Each
func (rcv {{.SliceName}}) Each(fn func({{.Type}})) {
	for _, v := range rcv {
		fn(v)
	}
}
`},

	"First": &typewriter.Template{
		Text: `
// First returns the first element that returns true for the passed func. Returns error if no elements return true. See: http://clipperhouse.github.io/gen/#First
func (rcv {{.SliceName}}) First(fn func({{.Type}}) bool) (result {{.Type}}, err error) {
	for _, v := range rcv {
		if fn(v) {
			result = v
			return
		}
	}
	err = errors.New("no {{.SliceName}} elements return true for passed func")
	return
}
`},

	"Max": &typewriter.Template{
		Text: `
	// Max returns the maximum value of {{.SliceName}}. In the case of multiple items being equally maximal, the first such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#Max
	func (rcv {{.SliceName}}) Max() (result {{.Type}}, err error) {
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
		TypeConstraint: typewriter.Constraint{Ordered: true},
	},

	"Min": &typewriter.Template{
		Text: `
	// Min returns the minimum value of {{.SliceName}}. In the case of multiple items being equally minimal, the first such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#Min
	func (rcv {{.SliceName}}) Min() (result {{.Type}}, err error) {
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
		TypeConstraint: typewriter.Constraint{Ordered: true},
	},

	"MaxBy": &typewriter.Template{
		Text: `
// MaxBy returns an element of {{.SliceName}} containing the maximum value, when compared to other elements using a passed func defining ‘less’. In the case of multiple items being equally maximal, the last such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#MaxBy
func (rcv {{.SliceName}}) MaxBy(less func({{.Type}}, {{.Type}}) bool) (result {{.Type}}, err error) {
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

	"MinBy": &typewriter.Template{
		Text: `
// MinBy returns an element of {{.SliceName}} containing the minimum value, when compared to other elements using a passed func defining ‘less’. In the case of multiple items being equally minimal, the first such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#MinBy
func (rcv {{.SliceName}}) MinBy(less func({{.Type}}, {{.Type}}) bool) (result {{.Type}}, err error) {
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

	"Single": &typewriter.Template{
		Text: `
// Single returns exactly one element of {{.SliceName}} that returns true for the passed func. Returns error if no or multiple elements return true. See: http://clipperhouse.github.io/gen/#Single
func (rcv {{.SliceName}}) Single(fn func({{.Type}}) bool) (result {{.Type}}, err error) {
	var candidate {{.Type}}
	found := false
	for _, v := range rcv {
		if fn(v) {
			if found {
				err = errors.New("multiple {{.SliceName}} elements return true for passed func")
				return
			}
			candidate = v
			found = true
		}
	}
	if found {
		result = candidate
	} else {
		err = errors.New("no {{.SliceName}} elements return true for passed func")
	}
	return
}
`},

	"Where": &typewriter.Template{
		Text: `
// Where returns a new {{.SliceName}} whose elements return true for func. See: http://clipperhouse.github.io/gen/#Where
func (rcv {{.SliceName}}) Where(fn func({{.Type}}) bool) (result {{.SliceName}}) {
	for _, v := range rcv {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
`},

	"Sort": &typewriter.Template{
		Text: `
// Sort returns a new ordered {{.SliceName}}. See: http://clipperhouse.github.io/gen/#Sort
func (rcv {{.SliceName}}) Sort() {{.SliceName}} {
	result := make({{.SliceName}}, len(rcv))
	copy(result, rcv)
	sort.Sort(result)
	return result
}
`,
		TypeConstraint: typewriter.Constraint{Ordered: true},
	},
	"IsSorted": &typewriter.Template{
		Text: `
// IsSorted reports whether {{.SliceName}} is sorted. See: http://clipperhouse.github.io/gen/#Sort
func (rcv {{.SliceName}}) IsSorted() bool {
	return sort.IsSorted(rcv)
}
`,
		TypeConstraint: typewriter.Constraint{Ordered: true},
	},
	"SortDesc": &typewriter.Template{
		Text: `
// SortDesc returns a new reverse-ordered {{.SliceName}}. See: http://clipperhouse.github.io/gen/#Sort
func (rcv {{.SliceName}}) SortDesc() {{.SliceName}} {
	result := make({{.SliceName}}, len(rcv))
	copy(result, rcv)
	sort.Sort(sort.Reverse(result))
	return result
}
`,
		TypeConstraint: typewriter.Constraint{Ordered: true},
	},
	"IsSortedDesc": &typewriter.Template{
		Text: `
// IsSortedDesc reports whether {{.SliceName}} is reverse-sorted. See: http://clipperhouse.github.io/gen/#Sort
func (rcv {{.SliceName}}) IsSortedDesc() bool {
	return sort.IsSorted(sort.Reverse(rcv))
}
`,
		TypeConstraint: typewriter.Constraint{Ordered: true},
	},

	"SortBy": &typewriter.Template{
		Text: `
// SortBy returns a new ordered {{.SliceName}}, determined by a func defining ‘less’. See: http://clipperhouse.github.io/gen/#SortBy
func (rcv {{.SliceName}}) SortBy(less func({{.Type}}, {{.Type}}) bool) {{.SliceName}} {
	result := make({{.SliceName}}, len(rcv))
	copy(result, rcv)
	// Switch to heapsort if depth of 2*ceil(lg(n+1)) is reached.
	n := len(result)
	maxDepth := 0
	for i := n; i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSort{{.SliceName}}(result, less, 0, n, maxDepth)
	return result
}
`},

	"IsSortedBy": &typewriter.Template{
		Text: `
// IsSortedBy reports whether an instance of {{.SliceName}} is sorted, using the pass func to define ‘less’. See: http://clipperhouse.github.io/gen/#SortBy
func (rcv {{.SliceName}}) IsSortedBy(less func({{.Type}}, {{.Type}}) bool) bool {
	n := len(rcv)
	for i := n - 1; i > 0; i-- {
		if less(rcv[i], rcv[i-1]) {
			return false
		}
	}
	return true
}
`},

	"SortByDesc": &typewriter.Template{
		Text: `
// SortByDesc returns a new, descending-ordered {{.SliceName}}, determined by a func defining ‘less’. See: http://clipperhouse.github.io/gen/#SortBy
func (rcv {{.SliceName}}) SortByDesc(less func({{.Type}}, {{.Type}}) bool) {{.SliceName}} {
	greater := func(a, b {{.Type}}) bool {
		return less(b, a)
	}
	return rcv.SortBy(greater)
}
`},

	"IsSortedByDesc": &typewriter.Template{
		Text: `
// IsSortedDesc reports whether an instance of {{.SliceName}} is sorted in descending order, using the pass func to define ‘less’. See: http://clipperhouse.github.io/gen/#SortBy
func (rcv {{.SliceName}}) IsSortedByDesc(less func({{.Type}}, {{.Type}}) bool) bool {
	greater := func(a, b {{.Type}}) bool {
		return less(b, a)
	}
	return rcv.IsSortedBy(greater)
}
`},

	"sortInterface": &typewriter.Template{
		Text: `
func (rcv {{.SliceName}}) Len() int {
	return len(rcv)
}
func (rcv {{.SliceName}}) Less(i, j int) bool {
	return rcv[i] < rcv[j]
}
func (rcv {{.SliceName}}) Swap(i, j int) {
	rcv[i], rcv[j] = rcv[j], rcv[i]
}
`},

	"sortSupport": &typewriter.Template{
		Text: `
// Sort implementation based on http://golang.org/pkg/sort/#Sort, see top of this file

func swap{{.SliceName}}(rcv {{.SliceName}}, a, b int) {
	rcv[a], rcv[b] = rcv[b], rcv[a]
}

// Insertion sort
func insertionSort{{.SliceName}}(rcv {{.SliceName}}, less func({{.Type}}, {{.Type}}) bool, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && less(rcv[j], rcv[j-1]); j-- {
			swap{{.SliceName}}(rcv, j, j-1)
		}
	}
}

// siftDown implements the heap property on rcv[lo, hi).
// first is an offset into the array where the root of the heap lies.
func siftDown{{.SliceName}}(rcv {{.SliceName}}, less func({{.Type}}, {{.Type}}) bool, lo, hi, first int) {
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
		swap{{.SliceName}}(rcv, first+root, first+child)
		root = child
	}
}

func heapSort{{.SliceName}}(rcv {{.SliceName}}, less func({{.Type}}, {{.Type}}) bool, a, b int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown{{.SliceName}}(rcv, less, i, hi, first)
	}

	// Pop elements, largest first, into end of rcv.
	for i := hi - 1; i >= 0; i-- {
		swap{{.SliceName}}(rcv, first, first+i)
		siftDown{{.SliceName}}(rcv, less, lo, i, first)
	}
}

// Quicksort, following Bentley and McIlroy,
// Engineering a Sort Function, SP&E November 1993.

// medianOfThree moves the median of the three values rcv[a], rcv[b], rcv[c] into rcv[a].
func medianOfThree{{.SliceName}}(rcv {{.SliceName}}, less func({{.Type}}, {{.Type}}) bool, a, b, c int) {
	m0 := b
	m1 := a
	m2 := c
	// bubble sort on 3 elements
	if less(rcv[m1], rcv[m0]) {
		swap{{.SliceName}}(rcv, m1, m0)
	}
	if less(rcv[m2], rcv[m1]) {
		swap{{.SliceName}}(rcv, m2, m1)
	}
	if less(rcv[m1], rcv[m0]) {
		swap{{.SliceName}}(rcv, m1, m0)
	}
	// now rcv[m0] <= rcv[m1] <= rcv[m2]
}

func swapRange{{.SliceName}}(rcv {{.SliceName}}, a, b, n int) {
	for i := 0; i < n; i++ {
		swap{{.SliceName}}(rcv, a+i, b+i)
	}
}

func doPivot{{.SliceName}}(rcv {{.SliceName}}, less func({{.Type}}, {{.Type}}) bool, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2 // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		// Tukey's Ninther, median of three medians of three.
		s := (hi - lo) / 8
		medianOfThree{{.SliceName}}(rcv, less, lo, lo+s, lo+2*s)
		medianOfThree{{.SliceName}}(rcv, less, m, m-s, m+s)
		medianOfThree{{.SliceName}}(rcv, less, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThree{{.SliceName}}(rcv, less, lo, m, hi-1)

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
				swap{{.SliceName}}(rcv, a, b)
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
				swap{{.SliceName}}(rcv, c-1, d-1)
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
		swap{{.SliceName}}(rcv, b, c-1)
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
	swapRange{{.SliceName}}(rcv, lo, b-n, n)

	n = min(hi-d, d-c)
	swapRange{{.SliceName}}(rcv, c, hi-n, n)

	return lo + b - a, hi - (d - c)
}

func quickSort{{.SliceName}}(rcv {{.SliceName}}, less func({{.Type}}, {{.Type}}) bool, a, b, maxDepth int) {
	for b-a > 7 {
		if maxDepth == 0 {
			heapSort{{.SliceName}}(rcv, less, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivot{{.SliceName}}(rcv, less, a, b)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			quickSort{{.SliceName}}(rcv, less, a, mlo, maxDepth)
			a = mhi // i.e., quickSort{{.SliceName}}(rcv, mhi, b)
		} else {
			quickSort{{.SliceName}}(rcv, less, mhi, b, maxDepth)
			b = mlo // i.e., quickSort{{.SliceName}}(rcv, a, mlo)
		}
	}
	if b-a > 1 {
		insertionSort{{.SliceName}}(rcv, less, a, b)
	}
}
`},

	"Aggregate[T]": &typewriter.Template{
		Text: `
// Aggregate{{.TypeParameter.LongName}} iterates over {{.SliceName}}, operating on each element while maintaining ‘state’. See: http://clipperhouse.github.io/gen/#Aggregate
func (rcv {{.SliceName}}) Aggregate{{.TypeParameter.LongName}}(fn func({{.TypeParameter}}, {{.Type}}) {{.TypeParameter}}) (result {{.TypeParameter}}) {
	for _, v := range rcv {
		result = fn(result, v)
	}
	return
}
`,
		TypeParameterConstraints: []typewriter.Constraint{
			// exactly one type parameter is required, but no constraints on that type
			{},
		},
	},

	"Average[T]": &typewriter.Template{
		Text: `
// Average{{.TypeParameter.LongName}} sums {{.Type}} over all elements and divides by len({{.SliceName}}). See: http://clipperhouse.github.io/gen/#Average
func (rcv {{.SliceName}}) Average{{.TypeParameter.LongName}}(fn func({{.Type}}) {{.TypeParameter}}) (result {{.TypeParameter}}, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine Average of zero-length {{.SliceName}}")
		return
	}
	for _, v := range rcv {
		result += fn(v)
	}
	result = result / {{.TypeParameter}}(l)
	return
}
`,
		TypeParameterConstraints: []typewriter.Constraint{
			// exactly one type parameter is required, and it must be numeric
			{Numeric: true},
		},
	},

	"GroupBy[T]": &typewriter.Template{
		Text: `
// GroupBy{{.TypeParameter.LongName}} groups elements into a map keyed by {{.TypeParameter}}. See: http://clipperhouse.github.io/gen/#GroupBy
func (rcv {{.SliceName}}) GroupBy{{.TypeParameter.LongName}}(fn func({{.Type}}) {{.TypeParameter}}) map[{{.TypeParameter}}]{{.SliceName}} {
	result := make(map[{{.TypeParameter}}]{{.SliceName}})
	for _, v := range rcv {
		key := fn(v)
		result[key] = append(result[key], v)
	}
	return result
}
`,
		TypeParameterConstraints: []typewriter.Constraint{
			// exactly one type parameter is required, and it must be comparable
			{Comparable: true},
		},
	},

	"Max[T]": &typewriter.Template{
		Text: `
// Max{{.TypeParameter.LongName}} selects the largest value of {{.TypeParameter}} in {{.SliceName}}. Returns error on {{.SliceName}} with no elements. See: http://clipperhouse.github.io/gen/#MaxCustom
func (rcv {{.SliceName}}) Max{{.TypeParameter.LongName}}(fn func({{.Type}}) {{.TypeParameter}}) (result {{.TypeParameter}}, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine Max of zero-length {{.SliceName}}")
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
`,
		TypeParameterConstraints: []typewriter.Constraint{
			// exactly one type parameter is required, and it must be ordered
			{Ordered: true},
		},
	},

	"Min[T]": &typewriter.Template{
		Text: `
// Min{{.TypeParameter.LongName}} selects the least value of {{.TypeParameter}} in {{.SliceName}}. Returns error on {{.SliceName}} with no elements. See: http://clipperhouse.github.io/gen/#MinCustom
func (rcv {{.SliceName}}) Min{{.TypeParameter.LongName}}(fn func({{.Type}}) {{.TypeParameter}}) (result {{.TypeParameter}}, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine Min of zero-length {{.SliceName}}")
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
`,
		TypeParameterConstraints: []typewriter.Constraint{
			// exactly one type parameter is required, and it must be ordered
			{Ordered: true},
		},
	},

	"Select[T]": &typewriter.Template{
		Text: `
// Select{{.TypeParameter.LongName}} returns a slice of {{.Type}} in {{.SliceName}}, projected by passed func. See: http://clipperhouse.github.io/gen/#Select
func (rcv {{.SliceName}}) Select{{.TypeParameter.LongName}}(fn func({{.Type}}) {{.TypeParameter}}) (result []{{.TypeParameter}}) {
	for _, v := range rcv {
		result = append(result, fn(v))
	}
	return
}
`,
		TypeParameterConstraints: []typewriter.Constraint{
			// exactly one type parameter is required, but no constraints on that type
			{},
		},
	},

	"Sum[T]": &typewriter.Template{
		Text: `
// Sum{{.TypeParameter.LongName}} sums {{.Type}} over elements in {{.SliceName}}. See: http://clipperhouse.github.io/gen/#Sum
func (rcv {{.SliceName}}) Sum{{.TypeParameter.LongName}}(fn func({{.Type}}) {{.TypeParameter}}) (result {{.TypeParameter}}) {
	for _, v := range rcv {
		result += fn(v)
	}
	return
}
`,
		TypeParameterConstraints: []typewriter.Constraint{
			// exactly one type parameter is required, and it must be numeric
			{Numeric: true},
		},
	},
}
