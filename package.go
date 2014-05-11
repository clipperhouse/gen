package typewriter

import (
	_ "code.google.com/p/go.tools/go/gcimporter"
	"code.google.com/p/go.tools/go/types"
)

type Package struct {
	*types.Package
}

func (p *Package) Eval(name string) (result Type, err error) {
	t, _, typesErr := types.Eval(name, p.Package, p.Scope())
	if typesErr != nil {
		err = typesErr
		return
	}
	result = Type{
		Package:    p,
		Pointer:    isPointer(t),
		Name:       name,
		comparable: isComparable(t),
		numeric:    isNumeric(t),
		ordered:    isOrdered(t),
		Type:       t,
	}
	return
}
