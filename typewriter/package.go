package typewriter

import (
	"go/ast"
	"go/token"

	_ "code.google.com/p/go.tools/go/gcimporter"
	"code.google.com/p/go.tools/go/types"
)

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
		Name:       name,
		comparable: isComparable(t),
		numeric:    isNumeric(t),
		ordered:    isOrdered(t),
		Type:       t,
	}
	return
}
