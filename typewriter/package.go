package typewriter

import (
	"go/ast"
	"go/token"

	_ "code.google.com/p/go.tools/go/gcimporter"
	"code.google.com/p/go.tools/go/types"
)

type Package struct {
	*types.Package
	fset       *token.FileSet
	astPackage *ast.Package
	info       *types.Info
}

func NewPackage(path, name string) *Package {
	return &Package{
		Package: types.NewPackage(path, name),
	}
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

func (pkg *Package) GetSelectorsOn(search types.Type) []string {
	data := walkerData{pkg.Package, pkg.fset, search, pkg.info.Scopes, make(map[string]struct{})}
	ast.Walk(walker{pkg.Package.Scope(), &data}, pkg.astPackage)

	selectors := make([]string, 0)
	for k := range data.selectors {
		selectors = append(selectors, k)
	}
	return selectors
}

type walkerData struct {
	typesPkg  *types.Package
	fset      *token.FileSet
	search    types.Type
	scopes    map[ast.Node]*types.Scope
	selectors map[string]struct{}
}

type walker struct {
	scope *types.Scope
	data  *walkerData
}

func (w walker) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return w
	}
	d := w.data
	if newscope, ok := d.scopes[node]; ok {
		return walker{newscope, d}
	}
	switch n := node.(type) {
	case *ast.FuncDecl:
		newscope := d.scopes[n.Type]
		return walker{newscope, d}
	case *ast.SelectorExpr:
		ident := n.Sel.Name
		typ, _, err := types.EvalNode(d.fset, n.X, d.typesPkg, w.scope)
		if err == nil && types.Identical(typ, d.search) {
			d.selectors[ident] = struct{}{}
		}
	}
	return w
}
