package main

import (
	"errors"
	"fmt"
	"go/ast"
	"os"
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

	typeArgs, opts, errs := parseArgs(args)

	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		fmt.Println("type 'gen' to see usage")
		return // command-line errors are fatal, other errors can be forced
	}

	packages := getPackages()
	genSpecs := getGenSpecs(opts, typeArgs, packages)

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

func getGenSpecs(opts *options, typeArgs []*typeArg, packages map[string]*Package) (genSpecs []*genSpec) {
	if len(typeArgs) > 0 && opts.All {
		addError(fmt.Sprintf("you've specified a type as well as the -all option; please choose one or the other"))
	}

	// 1. gather up info on types to be gen'd; strictly parsing, no validation
	typs := make([]*Type, 0)

	for _, t := range typeArgs {
		p, ok := packages[t.Package]

		if ok {
			typ, err := p.GetType(t)
			if err != nil {
				errs = append(errs, err)
			}
			typs = append(typs, typ)
		} else {
			addError(fmt.Sprintf("%s is not a known package", t.Package))
			typ := newType(t)
			typs = append(typs, typ)
		}
	}

	if opts.All {
		for k, p := range packages {
			for s := range p.TypeNamesAndDocs {
				if !opts.ExportedOnly || ast.IsExported(s) {
					t := &typeArg{opts.AllPointer, k, s}
					typ, err := p.GetType(t)
					if err != nil {
						errs = append(errs, err)
					}
					typs = append(typs, typ)
				}
			}
		}
	}

	// 2. create specs including type validation
	for _, typ := range typs {
		g := newGenSpec(typ)

		var stdMethods, prjMethods []string

		if len(typ.SubsettedMethods) > 0 {
			for _, m := range typ.SubsettedMethods {
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

			if len(typ.ProjectedTypes) > 0 && len(prjMethods) == 0 {
				addError(fmt.Sprintf("you've included projection types without specifying projection methods on type %s", g.Type()))
			}

			if len(prjMethods) > 0 && len(typ.ProjectedTypes) == 0 {
				addError(fmt.Sprintf("you've included projection methods without specifying projection types on type %s", g.Type()))
			}
		} else {
			stdMethods = getStandardMethodKeys()
			if len(typ.ProjectedTypes) > 0 {
				prjMethods = getProjectionMethodKeys()
			}
		}

		g.Methods = stdMethods

		for _, s := range typ.ProjectedTypes {
			isNumeric := false
			isComparable := true
			isOrdered := true

			p := packages[typ.Package]
			t, err := p.Eval(s)
			knownType := err == nil

			if err != nil {
				addError(fmt.Sprintf("unable to identify type %s, projected on %s (%s)", s, typ, err))
			} else {
				switch x := t.(type) {
				case *types.Slice:
					isComparable = false
					isOrdered = false
				case *types.Array:
					isComparable = false
					isOrdered = false
				case *types.Chan:
					isComparable = false
					isOrdered = false
				case *types.Map:
					isComparable = false
					isOrdered = false
				case *types.Struct:
					isComparable = true
					isOrdered = false
				default:
					switch u := x.Underlying().(type) {
					case *types.Basic:
						isNumeric = u.Info()|types.IsNumeric == types.IsNumeric
						isOrdered = u.Info()|types.IsOrdered == types.IsOrdered
					}
				}
			}

			for _, m := range prjMethods {
				pm, ok := projectionMethods[m]

				if !ok {
					addError(fmt.Sprintf("unknown projection method %v", m))
					continue
				}

				valid := (knownType || opts.Force) && (!pm.requiresNumeric || isNumeric) && (!pm.requiresComparable || isComparable) && (!pm.requiresOrdered || isOrdered)

				if valid {
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
