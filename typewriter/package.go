package typewriter

import (
	"go/ast"
	"go/token"
	"strings"

	_ "code.google.com/p/go.tools/go/gcimporter"
	"code.google.com/p/go.tools/go/types"
)

type evaluator interface {
	Eval(string) (Type, error)
}

type Package struct {
	*types.Package
}

func NewPackage(path, name string) *Package {
	return &Package{
		types.NewPackage(path, name),
	}
}

func getPackage(fset *token.FileSet, a *ast.Package) (*Package, error) {
	// pull map into a slice
	var files []*ast.File
	for _, f := range a.Files {
		files = append(files, f)
	}

	typesPkg, err := types.Check(a.Name, fset, files)

	if err != nil {
		return nil, err
	}

	return &Package{typesPkg}, nil
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
		Name:       strings.TrimLeft(name, Pointer(true).String()), // trims the * if it exists
		comparable: isComparable(t),
		numeric:    isNumeric(t),
		ordered:    isOrdered(t),
		Type:       t,
	}
	return
}
