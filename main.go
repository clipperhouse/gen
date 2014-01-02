package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
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

func getTypeArgs(args []string) (typeArgs []typeString) {
	regex := regexp.MustCompile(`^(\*?)([\p{L}\p{N}]+)\.([\p{L}\p{N}]+)$`)

	for _, s := range args {
		if regex.MatchString(s) {
			typeArgs = append(typeArgs, typeString(s))
		}
	}

	return
}

func getOptions(args []string) *options {
	opts := &options{}

	allOption := regexp.MustCompile(`^-(\*?)a(ll)?$`)
	exportedOption := regexp.MustCompile(`^-e(xported)?$`)
	forceOption := regexp.MustCompile(`^-f(orce)?$`)

	for _, a := range args {
		allMatches := allOption.FindStringSubmatch(a)
		if allMatches != nil {
			opts.All = true
			opts.AllPointer = allMatches[1]
		}
		if exportedOption.MatchString(a) {
			opts.ExportedOnly = true
		}
		if forceOption.MatchString(a) {
			opts.Force = true
		}
	}

	return opts
}

// Returns one type checker per package found in current directory
func getTypeCheckers() (result map[string]*typeChecker) {
	fset := token.NewFileSet()
	dir, err := parser.ParseDir(fset, "./", nil, parser.ParseComments)
	if err != nil {
		errs = append(errs, err)
		return
	}

	result = make(map[string]*typeChecker)

	for k, v := range dir {
		files := make([]*ast.File, 0)
		for _, f := range v.Files {
			files = append(files, f)
		}

		p, err := types.Check(k, fset, files)
		if err != nil {
			errs = append(errs, err)
		}

		d := doc.New(v, k, doc.AllDecls)
		typeDocs := make(map[string]string)
		for _, t := range d.Types {
			typeDocs[t.Name] = t.Doc
		}

		result[k] = &typeChecker{p, typeDocs}
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

func getSubsettedMethods(t *typeSpec) (result []string) {
	genTag := regexp.MustCompile(`gen:"([A-Za-z,]+)"`)
	parse := genTag.FindStringSubmatch(t.Doc)
	if parse != nil && len(parse) > 1 {
		result = strings.Split(parse[1], ",")
	}
	return
}

func getProjectionTypes(t *typeSpec) (result []string) {
	projectTag := regexp.MustCompile(`\bproject:"([A-Za-z,]+)"`)
	parse := projectTag.FindStringSubmatch(t.Doc)
	if parse != nil && len(parse) > 1 {
		result = strings.Split(parse[1], ",")
	}
	return
}

func getGenSpecs(opts *options, typeArgs []typeString, typeCheckers map[string]*typeChecker) (genSpecs []*genSpec) {
	if len(typeArgs) > 0 && opts.All {
		addError(fmt.Sprintf("you've specified a type as well as the -all option; please choose one or the other"))
	}

	typeSpecs := make([]*typeSpec, 0)

	for _, typeArg := range typeArgs {
		tc := typeCheckers[typeArg.Package()]
		t, err := tc.getTypeSpec(typeArg.LocalName()) // type checker is already package-specific
		if err != nil {
			errs = append(errs, err)
		} else {
			typeSpecs = append(typeSpecs, &t)
		}
	}

	if opts.All {
		for _, tc := range typeCheckers {
			for k := range tc.typeDocs {
				if !opts.ExportedOnly || ast.IsExported(k) {
					t, err := tc.getTypeSpec(opts.AllPointer + k)
					if err != nil {
						errs = append(errs, err)
					} else {
						typeSpecs = append(typeSpecs, &t)
					}
				}
			}
		}
	}

	for _, t := range typeSpecs {
		fmt.Println(t.Type.String() + " is:")
		switch x := t.Type.(type) {
		case *types.Pointer:
			fmt.Println("Pointer")
		case *types.Named:
			fmt.Println("Named")
		default:
			fmt.Println("dunno")
			fmt.Println(x)
		}
	}

	for _, g := range genSpecs {
		g.DetermineImports()
	}

	return
}

func populateGenSpec(g *genSpec, allTypes map[string]*typeSpec) {
	var subsettedMethods, standardMethods, projectionMethods, projectedTypes []string

	key := joinName("", g.Singular)
	typ, known := allTypes[key]

	if known {
		projectedTypes = getProjectionTypes(typ)
		subsettedMethods = getSubsettedMethods(typ)
		for _, m := range subsettedMethods {
			if isProjectionMethod(m) {
				projectionMethods = append(projectionMethods, m)
			} else {
				standardMethods = append(standardMethods, m)
			}
		}
	} else {
		addError(fmt.Sprintf("%s is not a known type", key))
	}

	if len(subsettedMethods) > 0 {
		g.Methods = standardMethods
		if len(projectedTypes) > 0 {
			if len(projectionMethods) == 0 { // TODO: reduce nesting
				addError(fmt.Sprintf("you've included projection types without specifying projection methods on %s", key))
			}
			g.Projections = getProjectionSpecs(g, projectionMethods, projectedTypes)
		}
	} else {
		g.Methods = getStandardMethodKeys()
		if len(projectedTypes) > 0 {
			g.Projections = getProjectionSpecs(g, getProjectionMethodKeys(), projectedTypes)
		}
	}
}

func getProjectionSpecs(g *genSpec, methods []string, types []string) (result []*projection) {
	for _, m := range methods {
		for _, t := range types {
			result = append(result, &projection{m, t, g})
		}
	}
	return
}

func joinName(pkg, name string) string {
	return fmt.Sprintf("%s.%s", pkg, name)
}

func splitName(s string) (string, string) {
	names := strings.Split(s, ".")
	return names[0], names[1]
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
