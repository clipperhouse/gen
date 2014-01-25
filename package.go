package main

import (
	_ "code.google.com/p/go.tools/go/gcimporter"
	"code.google.com/p/go.tools/go/types"
	"errors"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

type Package struct {
	Name  string
	Types []*Type
}

type GenSpec struct {
	Pointer, Name        string // Name is included mainly for informative error messages
	Methods, Projections *GenTag
}

type GenTag struct {
	Items []string
}

// Returns one gen Package per Go package found in current directory
func getPackages() (result []*Package) {
	fset := token.NewFileSet()
	astPackages, err := parser.ParseDir(fset, "./", nil, parser.ParseComments)
	if err != nil {
		errs = append(errs, err)
	}

	for name, astPackage := range astPackages {
		pkg := &Package{Name: name}

		typesPkg, err := types.Check(name, fset, getAstFiles(astPackage))
		if err != nil {
			errs = append(errs, err)
		}

		// fall back to Universe scope if types.Check fails; "best-effort" to handle primitives, at least
		scope := types.Universe
		if typesPkg != nil {
			scope = typesPkg.Scope()
		}

		docPkg := doc.New(astPackage, name, doc.AllDecls)
		for _, docType := range docPkg.Types {

			// look for deprecated struct tags, used for 'custom methods' in older version of gen
			if t, _, err := types.Eval(docType.Name, typesPkg, scope); err == nil {
				checkDeprecatedTags(t)
			}

			// identify marked-up types
			spec, found := getGenSpec(docType.Doc, docType.Name)
			if !found {
				continue
			}

			standardMethods, projectionMethods, err := determineMethods(spec)
			if err != nil {
				errs = append(errs, err)
			}

			typ := &Type{Package: pkg, Pointer: spec.Pointer, Name: docType.Name, StandardMethods: standardMethods}

			// assemble projections with type verification
			if spec.Projections != nil {
				for _, s := range spec.Projections.Items {
					numeric := false
					comparable := true // sensible default?
					ordered := false

					t, _, err := types.Eval(s, typesPkg, scope)
					known := err == nil

					if !known {
						addError(fmt.Sprintf("unable to identify type %s, projected on %s (%s)", s, docType.Name, err))
					} else {
						numeric = isNumeric(t)
						comparable = isComparable(t)
						ordered = isOrdered(t)
					}

					for _, m := range projectionMethods {
						pm, ok := ProjectionMethods[m]

						if !ok {
							addError(fmt.Sprintf("unknown projection method %v", m))
							continue
						}

						valid := (!pm.RequiresNumeric || numeric) && (!pm.RequiresComparable || comparable) && (!pm.RequiresOrdered || ordered)

						if valid {
							typ.AddProjection(m, s)
						}
					}
				}
			}

			determineImports(typ)

			pkg.Types = append(pkg.Types, typ)
		}

		// only add it to the results if there is something there
		if len(pkg.Types) > 0 {
			result = append(result, pkg)
		}
	}

	return
}

// getGenSpec identifies gen-marked types and parses tags
func getGenSpec(doc, name string) (result *GenSpec, found bool) {
	lines := strings.Split(doc, "\n")
	for _, line := range lines {
		if line = strings.TrimLeft(line, "/ "); strings.HasPrefix(line, "+gen") {
			// parse out tags & pointer
			spaces := regexp.MustCompile(" +")
			parts := spaces.Split(line, -1)

			var pointer string
			var subsettedMethods, projectedTypes *GenTag

			for _, s := range parts {
				if s == "*" {
					pointer = s
				}
				if x, found := parseTag("methods", s); found {
					subsettedMethods = &GenTag{x}
				}
				if x, found := parseTag("projections", s); found {
					projectedTypes = &GenTag{x}
				}
			}

			found = true
			result = &GenSpec{pointer, name, subsettedMethods, projectedTypes}
			return
		}
	}
	return
}

func determineMethods(spec *GenSpec) (standardMethods, projectionMethods []string, err error) {
	if spec.Methods != nil {
		// categorize subsetted methods as standard or projection
		for _, m := range spec.Methods.Items {
			if isStandardMethod(m) {
				standardMethods = append(standardMethods, m)
			}
			if isProjectionMethod(m) {
				projectionMethods = append(projectionMethods, m)
			}
			if !isStandardMethod(m) && !isProjectionMethod(m) {
				err = errors.New(fmt.Sprintf("method %s is unknown", m, spec.Name))
			}
		}

		if spec.Projections != nil && len(projectionMethods) == 0 {
			err = errors.New(fmt.Sprintf("you've included projection types without specifying projection methods on type %s", spec.Name))
		}

		if len(projectionMethods) > 0 && spec.Projections == nil {
			err = errors.New(fmt.Sprintf("you've included projection methods without specifying projection types on type %s", spec.Name))
		}
	} else {
		// default to all if not subsetted
		standardMethods = getStandardMethodKeys()
		if spec.Projections != nil {
			projectionMethods = getProjectionMethodKeys()
		}
	}
	return
}

func getAstFiles(p *ast.Package) (result []*ast.File) {
	// pull map of *ast.File into a slice
	for _, f := range p.Files {
		result = append(result, f)
	}
	return
}

func parseTag(name, s string) (result []string, found bool) {
	pattern := fmt.Sprintf(`%s:"(.*)"`, name)
	r := regexp.MustCompile(pattern)
	if matches := r.FindStringSubmatch(s); matches != nil && len(matches) > 1 {
		found = true
		match := matches[1]
		if len(match) > 0 {
			result = strings.Split(match, ",")
		}
	}
	return
}

func determineImports(t *Type) {
	imports := make(map[string]bool)
	methodRequiresErrors := map[string]bool{
		"First":   true,
		"Single":  true,
		"MaxBy":   true,
		"MinBy":   true,
		"Average": true,
	}

	for _, m := range t.StandardMethods {
		if methodRequiresErrors[m] {
			imports["errors"] = true
		}
	}

	for _, f := range t.Projections {
		if methodRequiresErrors[f.Method] {
			imports["errors"] = true
		}
	}

	for s := range imports {
		t.Imports = append(t.Imports, s)
	}
}

func (t *Type) requiresSortSupport() bool {
	for _, m := range t.StandardMethods {
		if strings.HasPrefix(m, "SortBy") {
			return true
		}
	}
	return false
}

func checkDeprecatedTags(t types.Type) {
	// give informative errors for use of deprecated custom methods
	switch x := t.Underlying().(type) {
	case *types.Struct:
		for i := 0; i < x.NumFields(); i++ {
			_, found := parseTag("gen", x.Tag(i))
			if found {
				addError(fmt.Sprintf(`custom methods (%s on %s) have been deprecated, see %s`, x.Tag(i), x.Field(i).Name(), deprecationUrl))
			}
		}
	}
}
