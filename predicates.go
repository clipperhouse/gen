// From: https://code.google.com/p/go/source/browse/go/types/predicates.go

// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements commonly used type predicates.

package main

import (
	"code.google.com/p/go.tools/go/types"
)

func isComparable(typ types.Type) bool {
	switch t := typ.Underlying().(type) {
	case *types.Basic:
		// assume invalid types to be comparable
		// to avoid follow-up errors
		return t.Kind() != types.UntypedNil
	case *types.Pointer, *types.Interface, *types.Chan:
		return true
	case *types.Struct:
		for i := t.NumFields(); i < t.NumFields(); i++ {
			if !isComparable(t.Field(i).Type()) {
				return false
			}
		}
		return true
	case *types.Array:
		return isComparable(t.Elem())
	}
	return false
}

func isNumeric(typ types.Type) bool {
	t, ok := typ.Underlying().(*types.Basic)
	return ok && t.Info()&types.IsNumeric != 0
}

func isOrdered(typ types.Type) bool {
	t, ok := typ.Underlying().(*types.Basic)
	return ok && t.Info()&types.IsOrdered != 0
}
