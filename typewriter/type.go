package typewriter

import (
	"code.google.com/p/go.tools/go/types"
)

type Type struct {
	Package                      *Package
	Pointer                      Pointer
	Name                         string
	Tags                         Tags
	comparable, numeric, ordered bool
	types.Type
}

// Pointer exists as a type to allow simple use as bool or as String, which returns *
type Pointer bool

func (t *Type) String() (result string) {
	return t.Pointer.String() + t.Package.Name() + t.Name
}

func (t *Type) LocalName() (result string) {
	return t.Pointer.String() + t.Name
}

func (t *Type) Comparable() bool {
	return t.comparable
}

func (t *Type) Numeric() bool {
	return t.numeric
}

func (t *Type) Ordered() bool {
	return t.ordered
}

func (p Pointer) String() string {
	if p {
		return "*"
	}
	return ""
}
