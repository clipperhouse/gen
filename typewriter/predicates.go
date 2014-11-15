// From: https://code.google.com/p/go/source/browse/go/types/predicates.go?repo=tools

// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements commonly used type predicates.

package typewriter

import (
	"golang.org/x/tools/go/types"
)

func isComparable(typ types.Type) bool {
	return types.Comparable(typ)
}

func isNumeric(typ types.Type) bool {
	t, ok := typ.Underlying().(*types.Basic)
	return ok && t.Info()&types.IsNumeric != 0
}

func isOrdered(typ types.Type) bool {
	t, ok := typ.Underlying().(*types.Basic)
	return ok && t.Info()&types.IsOrdered != 0
}

func isPointer(typ types.Type) Pointer {
	_, ok := typ.Underlying().(*types.Pointer)
	return Pointer(ok)
}
