package main

import (
	"bitbucket.org/clipperhouse/inflect"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

type genSpec struct {
	Pointer    string
	Package    string
	Singular   string
	FieldSpecs []*fieldSpec
	Methods    []string
	Plural     string
	Receiver   string
	Loop       string
	Generated  string
	Command    string
	FileName   string
}

func newGenSpec(ptr, pkg, name string) *genSpec {
	name = inflect.Singularize(name)
	return &genSpec{
		Pointer:   ptr,
		Package:   pkg,
		Singular:  name,
		Plural:    inflect.Pluralize(name),
		Receiver:  "rcv",
		Loop:      "v",
		Generated: time.Now().UTC().Format(time.RFC1123),
		Command:   fmt.Sprintf("%s %s%s.%s", "gen", ptr, pkg, name),
		FileName:  strings.ToLower(name) + "_gen.go",
	}
}

func (g *genSpec) AddFieldSpecs(fieldSpecs []*fieldSpec) {
	for _, f := range fieldSpecs {
		f.Parent = g
	}
	g.FieldSpecs = fieldSpecs
}

func (g genSpec) String() string {
	return joinName(g.Package, g.Plural)
}

func (g genSpec) RequiresSortSupport() bool {
	for _, m := range g.Methods {
		if strings.HasPrefix(m, "Sort") {
			return true
		}
	}
	return false
}

type structArg struct {
	Pointer string
	Package string
	Name    string
}

type fieldSpec struct {
	Name    string
	Type    string
	Methods []string
	Parent  *genSpec
}

type options struct {
	All          bool
	AllPointer   string
	ExportedOnly bool
	Force        bool
}

var errs = make([]error, 0)
var notes = make([]string, 0)

func addError(text string) {
	errs = append(errs, errors.New(text))
}

func main() {
	has_args := len(os.Args) > 1
	if !has_args {
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
	structArgs := getStructArgs(args)
	genSpecs := getGenSpecs(opts, structArgs)

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

func getStructArgs(args []string) (structArgs []*structArg) {
	regex := regexp.MustCompile(`^(\*?)([\p{L}\p{N}]+)\.([\p{L}\p{N}]+)$`)

	for _, s := range args {
		matches := regex.FindStringSubmatch(s)

		if matches == nil {
			continue
		}

		ptr := matches[1]
		pkg := matches[2]
		typ := matches[3]

		structArgs = append(structArgs, &structArg{ptr, pkg, typ})
	}

	return
}

func getAllStructTypes(fset *token.FileSet) (types map[string]*ast.StructType) {
	goFiles := func(f os.FileInfo) bool {
		return strings.HasSuffix(f.Name(), ".go")
	}
	dir, err := parser.ParseDir(fset, "./", goFiles, parser.ParseComments)
	if err != nil {
		errs = append(errs, err)
		return
	}
	types = make(map[string]*ast.StructType)
	for pkg, f := range dir {
		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.TypeSpec:
				switch y := x.Type.(type) {
				case *ast.StructType:
					key := joinName(pkg, x.Name.Name)
					types[key] = y
				}
			}
			return true
		})
	}
	return
}

var genTag = regexp.MustCompile(`gen:"([A-Za-z,]+)"`)

func getMethods(typ *ast.StructType) (result []string) {
	// look for comments of the form gen:"Method,Method", like struct (field) tags but at type level
	tagged := false
	include := make(map[string]bool)
	ast.Inspect(typ, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Comment:
			c := strings.Trim(x.Text, " /")
			parse := genTag.FindStringSubmatch(c)
			if parse != nil && len(parse) > 1 {
				tagged = true
				methods := strings.Split(parse[1], ",")
				if len(methods) > 0 {
					for _, m := range methods {
						_, err := getStandardTemplate(m)
						if err != nil {
							errs = append(errs, err)
						} else {
							include[m] = true
						}
					}
				}
			}
		}
		return !tagged // stop inspecting after found
	})

	if !tagged {
		result = getStandardMethodKeys()
	}

	// dependency
	if include["SortDesc"] {
		include["Sort"] = true
	}

	for k := range include {
		result = append(result, k)
	}
	sort.Strings(result) // order of keys not guaranteed: http://blog.golang.org/go-maps-in-action#TOC_7.

	return
}

func getGenSpecs(opts *options, structArgs []*structArg) (genSpecs []*genSpec) {
	fset := token.NewFileSet()
	types := getAllStructTypes(fset)

	for _, structArg := range structArgs {
		key := joinName(structArg.Package, structArg.Name)
		typ, known := types[key]
		if known {
			fieldSpecs := getFieldSpecs(typ, fset, opts)
			g := newGenSpec(structArg.Pointer, structArg.Package, structArg.Name)
			g.Methods = getMethods(typ)
			g.AddFieldSpecs(fieldSpecs)
			genSpecs = append(genSpecs, g)
		} else {
			addError(fmt.Sprintf("%s is not a known struct type", key))
			genSpecs = append(genSpecs, newGenSpec(structArg.Pointer, structArg.Package, structArg.Name))
		}
		if opts.ExportedOnly {
			if ast.IsExported(structArg.Name) {
				notes = append(notes, fmt.Sprintf("the %s type is already exported; the -e[xported] flag is redundant (ignored)", structArg.Name))
			} else {
				addError(fmt.Sprintf("the %s type is not exported; the -e[xported] flag conflicts", structArg.Name))
			}
		}
	}
	if opts.All {
		for key, typ := range types {
			fieldSpecs := getFieldSpecs(typ, fset, opts)
			pkg, name := splitName(key)
			if !opts.ExportedOnly || ast.IsExported(name) {
				g := newGenSpec(opts.AllPointer, pkg, name)
				g.AddFieldSpecs(fieldSpecs)
				genSpecs = append(genSpecs, g)
			}
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

func getFieldSpecs(typ *ast.StructType, fset *token.FileSet, opts *options) (fieldSpecs []*fieldSpec) {
	for _, fld := range typ.Fields.List {
		if fld.Tag == nil {
			continue
		}

		parse := genTag.FindStringSubmatch(fld.Tag.Value)
		if parse == nil || len(parse) < 2 {
			continue
		}

		methods := strings.Split(parse[1], ",")
		for _, m := range methods {
			_, err := getCustomTemplate(m)
			if err != nil {
				errs = append(errs, err)
			}
		}
		for _, name := range fld.Names {
			t := getSourceString(fld.Type, fset)
			fieldSpecs = append(fieldSpecs, &fieldSpec{Name: name.String(), Type: t, Methods: methods})
		}
	}
	return
}

func getSourceString(node ast.Node, fset *token.FileSet) string {
	p1 := fset.Position(node.Pos())
	p2 := fset.Position(node.End())

	b := getFileBytes(p1.Filename)
	return string(b[p1.Offset:p2.Offset])
}

var filebytes = make(map[string][]byte) // cache
func getFileBytes(filename string) []byte {
	b, exists := filebytes[filename]
	if exists {
		return b
	}

	b2, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	filebytes[filename] = b2
	return b2
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
				panic(err) // shouldn't get here, should have been caught in getMethods
			}
		}

		for _, f := range g.FieldSpecs {
			for _, m := range f.Methods {
				c, err := getCustomTemplate(m)
				if err == nil {
					c.Execute(file, f)
				} else if opts.Force {
					fmt.Printf("  skipping %v custom method\n", m)
				} else {
					panic(err) // shouldn't get here, should have been caught in getFieldSpecs
				}
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
