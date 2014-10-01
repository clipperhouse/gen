package slice

import "github.com/clipperhouse/gen/typewriter"

var sortImplementation = &typewriter.Template{
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
`}
