package main

import (
	"errors"
	"fmt"
	"go/ast"
	"os"
	"regexp"
	"strings"

	_ "code.google.com/p/go.tools/go/gcimporter"
	"code.google.com/p/go.tools/go/types"
)

var errs = make([]error, 0)
var notes = make([]string, 0)

func main() {
	if len(os.Args) <= 1 { // no args, only "gen"
		fmt.Println(usage)
		return
	}
	args := os.Args[1:]

	first := args[0]
	if first == "-help" || first == "help" || first == "?" {
		fmt.Println(usage)
		return
	}

	opts := getOptions(args)
	typeArgs := getTypeArgs(args)
	tc := getTypeCheckers()
	genSpecs := getGenSpecs(opts, typeArgs, tc)

	for _, n := range notes {
		fmt.Println("  note: " + n)
	}

	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Printf("  error: %v\n", e)
		}
		if opts.Force {
			fmt.Println("  forced...")
		} else {
			fmt.Println("  operation canceled")
			fmt.Println("  use the -f flag if you wish to force generation (i.e., ignore errors)")
			return
		}
	}

	writeFile(genSpecs, opts)
}

func getTypeArgs(args []string) (typeArgs []string) {
	regex := regexp.MustCompile(`^(\*?)([\p{L}\p{N}]+)\.([\p{L}\p{N}]+)$`)

	for _, s := range args {
		if regex.MatchString(s) {
			typeArgs = append(typeArgs, s)
		}
	}

	return
}

func (g *genSpec) DetermineImports() {
	imports := make(map[string]bool)
	methodRequiresErrors := map[string]bool{
		"First":   true,
		"Single":  true,
		"Max":     true,
		"Min":     true,
		"Average": true,
	}

	for _, m := range g.Methods {
		if methodRequiresErrors[m] {
			imports["errors"] = true
		}
	}

	for _, f := range g.Projections {
		if methodRequiresErrors[f.Method] {
			imports["errors"] = true
		}
	}

	for k := range imports {
		g.Imports = append(g.Imports, k)
	}
}

func (g genSpec) RequiresSortSupport() bool {
	for _, m := range g.Methods {
		if strings.HasPrefix(m, "Sort") {
			return true
		}
	}
	return false
}

func addError(text string) {
	errs = append(errs, errors.New(text))
}

func getGenSpecs(opts *options, typeArgs []string, typeCheckers map[string]*typeChecker) (genSpecs []*genSpec) {
	if len(typeArgs) > 0 && opts.All {
		addError(fmt.Sprintf("you've specified a type as well as the -all option; please choose one or the other"))
	}

	// 1. gather up info on types to be gen'd; strictly parsing, no validation
	typeSpecs := make([]*typeSpec, 0)

	for _, typeArg := range typeArgs {
		p := typeString(typeArg).Package()
		tc, ok := typeCheckers[p]

		if !ok {
			addError(fmt.Sprintf("no typeChecker found for package %s", p))
		}

		typeSpecs = append(typeSpecs, tc.getTypeSpec(typeArg))
	}

	if opts.All {
		for p, tc := range typeCheckers {
			for k := range tc.typeDocs {
				if !opts.ExportedOnly || ast.IsExported(k) {
					typeSpecs = append(typeSpecs, tc.getTypeSpec(opts.AllPointer+p+"."+k))
				}
			}
		}
	}

	// 2. create specs including type validation
	for _, t := range typeSpecs {
		g := newGenSpec(t.Pointer, t.Package, t.Name)

		var stdMethods, prjMethods []string

		if len(t.SubsettedMethods) > 0 {
			for _, m := range t.SubsettedMethods {
				if isProjectionMethod(m) {
					prjMethods = append(prjMethods, m)
				}
				if isStandardMethod(m) {
					stdMethods = append(stdMethods, m)
				}
				if !isProjectionMethod(m) && !isStandardMethod(m) {
					addError(fmt.Sprintf("method %s (subsetted on type %s) is unknown", m, g.Type()))
				}
			}

			if len(t.ProjectedTypes) > 0 && len(prjMethods) == 0 {
				addError(fmt.Sprintf("you've included projection types without specifying projection methods on type %s", g.Type()))
			}

			if len(prjMethods) > 0 && len(t.ProjectedTypes) == 0 {
				addError(fmt.Sprintf("you've included projection methods without specifying projection types on type %s", g.Type()))
			}
		} else {
			stdMethods = getStandardMethodKeys()
			if len(t.ProjectedTypes) > 0 {
				prjMethods = getProjectionMethodKeys()
			}
		}

		g.Methods = stdMethods

		for _, s := range t.ProjectedTypes {
			tc := typeCheckers[t.Package]
			isNumeric := false

			typ, err := tc.eval(s)
			if err != nil {
				errs = append(errs, err)
			} else {
				switch u := typ.Underlying().(type) {
				case *types.Basic:
					isNumeric = u.Info()|types.IsNumeric == types.IsNumeric
				}
			}

			for _, m := range prjMethods {
				pm, ok := projectionMethods[m]

				if !ok {
					addError(fmt.Sprintf("unknown projection method %v", m))
					continue
				}

				if !pm.requiresNumeric || isNumeric || opts.Force {
					g.Projections = append(g.Projections, &projection{m, s, g})
				}
			}
		}

		g.DetermineImports()

		genSpecs = append(genSpecs, g)
	}

	return
}

func writeFile(genSpecs []*genSpec, opts *options) {
	for _, g := range genSpecs {
		file, err := os.Create(g.FileName)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		h := getHeaderTemplate()
		h.Execute(file, g)

		for _, m := range g.Methods {
			t, err := getStandardTemplate(m)
			if err == nil {
				t.Execute(file, g)
			} else if opts.Force {
				fmt.Printf("  skipping %v method\n", m)
			} else {
				panic(err) // shouldn't get here, should have been caught in getSubsettedMethods
			}
		}

		for _, f := range g.Projections {
			c, err := getProjectionTemplate(f.Method)
			if err == nil {
				c.Execute(file, f)
			} else if opts.Force {
				fmt.Printf("  skipping %v projection method\n", f.Method)
			} else {
				panic(err) // shouldn't get here, should have been caught in getProjectionSpecs
			}
		}

		if g.RequiresSortSupport() {
			s := getSortSupportTemplate()
			s.Execute(file, g)
		}

		fmt.Printf("  generated %s, yay!\n", g)
	}
}

const usage = `Usage: gen [[*]package.TypeName] [-[*]all] [-exported]

  *package.TypeName    # generate funcs for specified struct type; use leading * to specify pointer type (recommended)
  -all                 # generate all structs in current directory; use leading * to specify pointer type (recommended); shortcut -a or -*a
  -exported            # only generate exported structs; shortcut -e
  -force               # force generate, overriding errors; shortcut -f

Examples:
  gen -*models.Movie   # generates funcs for Movie type in the models package; generated Movies type is []*Movie
  gen -models.Movie    # generates funcs for Movie type; generated Movies type is []Movie
  gen -*all            # generates funcs for all struct types in current directory, as pointers
  gen -all             # generates funcs for all struct types in current directory, as values
  gen -*a -e           # generates funcs for all exported struct types in current directory, as pointers
`
