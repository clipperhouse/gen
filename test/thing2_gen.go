// This file was auto-generated using github.com/clipperhouse/gen
// Modifying this file is not recommended as it will likely be overwritten in the future

// Sort (if included below) is a modification of http://golang.org/pkg/sort/#Sort
// List (if included below) is a modification of http://golang.org/pkg/container/list/
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

// Min returns the minimum value of Thing2s. In the case of multiple items being equally minimal, the first such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#Min
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

// Thing2Element is an element of a linked list.
type Thing2Element struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Thing2Element

	// The list to which this element belongs.
	list *Thing2List

	// The value stored with this element.
	Value Thing2
}

// Next returns the next list element or nil.
func (e *Thing2Element) Next() *Thing2Element {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list element or nil.
func (e *Thing2Element) Prev() *Thing2Element {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Thing2List represents a doubly linked list.
// The zero value for Thing2List is an empty list ready to use.
type Thing2List struct {
	root Thing2Element // sentinel list element, only &root, root.prev, and root.next are used
	len  int           // current list length excluding (this) sentinel element
}

// Init initializes or clears list l.
func (l *Thing2List) Init() *Thing2List {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// New returns an initialized list.
func New() *Thing2List { return new(Thing2List).Init() }

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *Thing2List) Len() int { return l.len }

// Front returns the first element of list l or nil
func (l *Thing2List) Front() *Thing2Element {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last element of list l or nil.
func (l *Thing2List) Back() *Thing2Element {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// lazyInit lazily initializes a zero Thing2List value.
func (l *Thing2List) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// insert inserts e after at, increments l.len, and returns e.
func (l *Thing2List) insert(e, at *Thing2Element) *Thing2Element {
	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e
	e.list = l
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&Thing2Element{Value: v}, at).
func (l *Thing2List) insertValue(v Thing2, at *Thing2Element) *Thing2Element {
	return l.insert(&Thing2Element{Value: v}, at)
}

// remove removes e from its list, decrements l.len, and returns e.
func (l *Thing2List) remove(e *Thing2Element) *Thing2Element {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
	return e
}

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value.
func (l *Thing2List) Remove(e *Thing2Element) Thing2 {
	if e.list == l {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero Thing2Element) and l.remove will crash
		l.remove(e)
	}
	return e.Value
}

// PushFront inserts a new element e with value v at the front of list l and returns e.
func (l *Thing2List) PushFront(v Thing2) *Thing2Element {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

// PushBack inserts a new element e with value v at the back of list l and returns e.
func (l *Thing2List) PushBack(v Thing2) *Thing2Element {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}

// InsertBefore inserts a new element e with value v immediately before mark and returns e.
// If mark is not an element of l, the list is not modified.
func (l *Thing2List) InsertBefore(v Thing2, mark *Thing2Element) *Thing2Element {
	if mark.list != l {
		return nil
	}
	// see comment in Thing2List.Remove about initialization of l
	return l.insertValue(v, mark.prev)
}

// InsertAfter inserts a new element e with value v immediately after mark and returns e.
// If mark is not an element of l, the list is not modified.
func (l *Thing2List) InsertAfter(v Thing2, mark *Thing2Element) *Thing2Element {
	if mark.list != l {
		return nil
	}
	// see comment in Thing2List.Remove about initialization of l
	return l.insertValue(v, mark)
}

// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
func (l *Thing2List) MoveToFront(e *Thing2Element) {
	if e.list != l || l.root.next == e {
		return
	}
	// see comment in Thing2List.Remove about initialization of l
	l.insert(l.remove(e), &l.root)
}

// MoveToBack moves element e to the back of list l.
// If e is not an element of l, the list is not modified.
func (l *Thing2List) MoveToBack(e *Thing2Element) {
	if e.list != l || l.root.prev == e {
		return
	}
	// see comment in Thing2List.Remove about initialization of l
	l.insert(l.remove(e), l.root.prev)
}

// MoveBefore moves element e to its new position before mark.
// If e is not an element of l, or e == mark, the list is not modified.
func (l *Thing2List) MoveBefore(e, mark *Thing2Element) {
	if e.list != l || e == mark {
		return
	}
	l.insert(l.remove(e), mark.prev)
}

// MoveAfter moves element e to its new position after mark.
// If e is not an element of l, or e == mark, the list is not modified.
func (l *Thing2List) MoveAfter(e, mark *Thing2Element) {
	if e.list != l || e == mark {
		return
	}
	l.insert(l.remove(e), mark)
}

// PushBackList inserts a copy of an other list at the back of list l.
// The lists l and other may be the same.
func (l *Thing2List) PushBackList(other *Thing2List) {
	l.lazyInit()
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Value, l.root.prev)
	}
}

// PushFrontList inserts a copy of an other list at the front of list l.
// The lists l and other may be the same.
func (l *Thing2List) PushFrontList(other *Thing2List) {
	l.lazyInit()
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Value, &l.root)
	}
}
