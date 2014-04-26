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

package typewriter

// Tags is a slice of type Tag, for use with gen methods below. Use this type where you would use []Tag. (This is required because slices cannot be method receivers.)
type Tags []Tag

// Where returns a new Tags slice whose elements return true for func. See: http://clipperhouse.github.io/gen/#Where
func (rcv Tags) Where(fn func(Tag) bool) (result Tags) {
	for _, v := range rcv {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
