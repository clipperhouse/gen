// This file was auto-generated using github.com/clipperhouse/gen
// Modifying this file is not recommended as it will likely be overwritten in the future

// Sort functions are a modification of http://golang.org/pkg/sort/#Sort
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package models

import (
	"errors"
	"sort"
)

// Thing2s is a slice of type Thing2, for use with gen methods below. Use this type where you would use []Thing2. (This is required because slices cannot be method receivers.)
type Thing2s []Thing2

// Max returns the maximum value of Thing2s. In the case of multiple items being equally maximal, the first such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#Max
func (rcv Thing2s) Max() (result Thing2, err error) {
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

// Min returns the minimum value of Thing2s. In the case of multiple items being equally minimal, the first such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#Max
func (rcv Thing2s) Min() (result Thing2, err error) {
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

// Sort returns a new ordered Thing2s slice. See: http://clipperhouse.github.io/gen/#Sort
func (rcv Thing2s) Sort() Thing2s {
	result := make(Thing2s, len(rcv))
	copy(result, rcv)
	sort.Sort(result)
	return result
}

// IsSorted reports whether Thing2s is sorted. See: http://clipperhouse.github.io/gen/#Sort
func (rcv Thing2s) IsSorted() bool {
	return sort.IsSorted(rcv)
}

// SortDesc returns a new reverse-ordered Thing2s slice. See: http://clipperhouse.github.io/gen/#Sort
func (rcv Thing2s) SortDesc() Thing2s {
	result := make(Thing2s, len(rcv))
	copy(result, rcv)
	sort.Sort(sort.Reverse(result))
	return result
}

// IsSortedDesc reports whether Thing2s is reverse-sorted. See: http://clipperhouse.github.io/gen/#Sort
func (rcv Thing2s) IsSortedDesc() bool {
	return sort.IsSorted(sort.Reverse(rcv))
}

func (rcv Thing2s) Len() int {
	return len(rcv)
}
func (rcv Thing2s) Less(i, j int) bool {
	return rcv[i] < rcv[j]
}
func (rcv Thing2s) Swap(i, j int) {
	rcv[i], rcv[j] = rcv[j], rcv[i]
}
