package typewriter

import (
	_ "code.google.com/p/go.tools/go/gcimporter"
	"code.google.com/p/go.tools/go/types"
)

type Package struct {
	Name     string
	Types    []Type
	typesPkg *types.Package // reference held for Eval below
}

func (p *Package) Eval(name string) (result *Type, err error) {
	t, _, typesErr := types.Eval(name, p.typesPkg, p.typesPkg.Scope())
	if typesErr != nil {
		err = typesErr
		return
	}
	result = &Type{
		Package:    p,
		Pointer:    isPointer(t),
		Name:       name,
		comparable: isComparable(t),
		numeric:    isNumeric(t),
		ordered:    isOrdered(t),
	}
	return
}
