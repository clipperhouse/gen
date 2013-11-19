package main

import (
	"bitbucket.org/clipperhouse/inflect"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"
)

type genSpec struct {
	Pointer    string
	Package    string
	Singular   string
	FieldSpecs []*fieldSpec
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
		Loop:      "_item",
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
	structTypes := getAllStructTypes()
	genSpecs := getGenSpecs(opts, structArgs, structTypes)

	if len(notes) > 0 {
		for _, n := range notes {
			fmt.Println("  note: " + n)
		}
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
	t := getTemplate()
	writeFile(t, genSpecs)
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
	regex := regexp.MustCompile(`^(\*?)(\p{L}+)\.(\p{L}+)$`)

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

func getAllStructTypes() map[string]*ast.StructType {
	goFiles := func(f os.FileInfo) bool {
		return strings.HasSuffix(f.Name(), ".go")
	}

	fset := token.NewFileSet()
	dir, err := parser.ParseDir(fset, "./", goFiles, parser.ParseComments)
	if err != nil {
		errs = append(errs, err)
		return nil
	}

	structTypes := make(map[string]*ast.StructType)
	for pkg, f := range dir {
		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.TypeSpec:
				switch y := x.Type.(type) {
				case *ast.StructType:
					key := joinName(pkg, x.Name.Name)
					structTypes[key] = y
				}
			}
			return true
		})
	}
	return structTypes
}

func getGenSpecs(opts *options, structArgs []*structArg, structTypes map[string]*ast.StructType) (genSpecs []*genSpec) {
	for _, structArg := range structArgs {
		key := joinName(structArg.Package, structArg.Name)
		typ, known := structTypes[key]
		if known {
			fieldSpecs := getFieldSpecs(typ)
			g := newGenSpec(structArg.Pointer, structArg.Package, structArg.Name)
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
		for key, typ := range structTypes {
			fieldSpecs := getFieldSpecs(typ)
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

func getFieldSpecs(typ *ast.StructType) (fieldSpecs []*fieldSpec) {
	genTag := regexp.MustCompile(`gen:"([A-Za-z,]+)"`)

	for _, fld := range typ.Fields.List {
		if fld.Tag != nil {
			parse := genTag.FindStringSubmatch(fld.Tag.Value)
			if parse != nil && len(parse) > 1 {
				methods := strings.Split(parse[1], ",")
				for _, name := range fld.Names {
					typeName := fmt.Sprintf("%v", fld.Type)
					fieldSpecs = append(fieldSpecs, &fieldSpec{Name: name.String(), Type: typeName, Methods: methods})
				}
			}
		}
	}
	return
}

func writeFile(t *template.Template, genSpecs []*genSpec) {
	for _, g := range genSpecs {
		file, err := os.Create(g.FileName)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		t.Execute(file, g)

		for _, f := range g.FieldSpecs {
			if f.Methods != nil {
				for _, m := range f.Methods {
					c := getCustomTemplate(m)
					c.Execute(file, f)
				}
			}
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
