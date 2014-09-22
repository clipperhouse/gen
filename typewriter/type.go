package typewriter

import (
	"fmt"
	"regexp"
	"strings"

	"code.google.com/p/go.tools/go/types"
)

type Type struct {
	Package                      *Package
	Pointer                      Pointer
	Name                         string
	Tags                         Tags
	comparable, numeric, ordered bool
	test                         test
	types.Type
}

type test bool

// a convenience for using bool in file name, see WriteAll
func (t test) String() string {
	if t {
		return "_test"
	}
	return ""
}

func (t Type) String() (result string) {
	return fmt.Sprintf("%s%s", t.Pointer.String(), t.Name)
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

func (typ Type) LongName() string {
	name := ""
	r := regexp.MustCompile(`[\[\]{}*]`)
	els := r.Split(typ.String(), -1)
	for _, s := range els {
		name += strings.Title(s)
	}
	return name
}

// Pointer exists as a type to allow simple use as bool or as String, which returns *
type Pointer bool

func (p Pointer) String() string {
	if p {
		return "*"
	}
	return ""
}
