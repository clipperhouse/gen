package main

import (
	_ "code.google.com/p/go.tools/go/gcimporter"
	"code.google.com/p/go.tools/go/types"
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

		docPkg := doc.New(astPackage, name, doc.AllDecls)
		for _, docType := range docPkg.Types {
			// identify marked-up types
			genLine, found := getGenLine(docType)
			if !found {
				continue
			}

			// parse out tags & pointer
			spaces := regexp.MustCompile(" +")
			parts := spaces.Split(genLine, -1)

			var pointer string
			var subsettedMethods, projectedTypes []string

			for _, s := range parts {
				if s == "*" {
					pointer = s
				}
				if x, found := parseTag("methods", genLine); found {
					subsettedMethods = x
				}
				if x, found := parseTag("projections", genLine); found {
					projectedTypes = x
				}
			}

			var standardMethods, projectionMethods []string

			if len(subsettedMethods) > 0 {
				// categorize subsetted methods as standard or projection
				for _, m := range subsettedMethods {
					if isStandardMethod(m) {
						standardMethods = append(standardMethods, m)
					}
					if isProjectionMethod(m) {
						projectionMethods = append(projectionMethods, m)
					}
					if !isStandardMethod(m) && !isProjectionMethod(m) {
						addError(fmt.Sprintf("method %s (subsetted on type %s) is unknown", m, docType.Name))
					}
				}

				if len(projectedTypes) > 0 && len(projectionMethods) == 0 {
					addError(fmt.Sprintf("you've included projection types without specifying projection methods on type %s", docType.Name))
				}

				if len(projectionMethods) > 0 && len(projectedTypes) == 0 {
					addError(fmt.Sprintf("you've included projection methods without specifying projection types on type %s", docType.Name))
				}
			} else {
				// default to all if not subsetted
				standardMethods = getStandardMethodKeys()
				if len(projectedTypes) > 0 {
					projectionMethods = getProjectionMethodKeys()
				}
			}

			typ := &Type{Package: pkg, Pointer: pointer, Name: docType.Name, StandardMethods: standardMethods}

			// assemble projections with type verification
			for _, s := range projectedTypes {
				numeric := false
				comparable := true // sensible default?
				ordered := false

				t, _, err := types.Eval(s, typesPkg, typesPkg.Scope())
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

					valid := (!pm.requiresNumeric || numeric) && (!pm.requiresComparable || comparable) && (!pm.requiresOrdered || ordered)

					if valid {
						typ.AddProjection(m, s)
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

func getGenLine(t *doc.Type) (result string, found bool) {
	lines := strings.Split(t.Doc, "\n")
	for _, line := range lines {
		if line = strings.TrimLeft(line, "/ "); strings.HasPrefix(line, "+gen") {
			found = true
			result = line
			return
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
	pattern := fmt.Sprintf(`%s:"(.+)"`, name)
	r := regexp.MustCompile(pattern)
	if matches := r.FindStringSubmatch(s); matches != nil && len(matches) > 1 {
		found = true
		result = strings.Split(matches[1], ",")
	}
	return
}

func determineImports(t *Type) {
	imports := make(map[string]bool)
	methodRequiresErrors := map[string]bool{
		"First":   true,
		"Single":  true,
		"Max":     true,
		"Min":     true,
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
		if strings.HasPrefix(m, "Sort") {
			return true
		}
	}
	return false
}
