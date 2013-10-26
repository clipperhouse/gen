package main

import (
	"bitbucket.org/pkg/inflect"
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
	Package   string
	Singular  string
	Plural    string
	Receiver  string
	Loop      string
	Pointer   string
	Generated string
	Command   string
	FileName  string
}

func (g genSpec) String() string {
	return fmt.Sprintf("%s.%s", g.Package, g.Plural)
}

type options struct {
	ExportedOnly bool
}

var genSpecs = make([]*genSpec, 0)

func main() {
	has_args := len(os.Args) > 1
	if !has_args {
		fmt.Println(usage)
		return
	}

	first := os.Args[1]
	if first == "-help" || first == "help" || first == "?" {
		fmt.Println(usage)
		return
	}

	for _, arg := range os.Args[1:] {
		valid := handleArg(arg)
		if !valid {
			fmt.Printf("Invalid argument: %s\n", arg)
			return
		}
	}

	t := getTemplate()
	writeFile(t, genSpecs)
}

func newGenSpec(ptr, pkg, typ string) *genSpec {
	typ = inflect.Singularize(typ)
	return &genSpec{
		Pointer:   ptr,
		Package:   pkg,
		Singular:  typ,
		Plural:    inflect.Pluralize(typ),
		Receiver:  "rcv",
		Loop:      "_item",
		Generated: time.Now().UTC().Format(time.RFC1123),
		Command:   fmt.Sprintf("%s %s%s.%s", "gen", ptr, pkg, typ),
		FileName:  strings.ToLower(typ) + "_gen.go",
	}
}

func handleArg(arg string) (valid bool) {
	genSpec, success := genSpecFromStructArg(arg)
	if success {
		genSpecs = append(genSpecs, genSpec)
		return true
	}

	all := allRegex.MatchString(arg)
	if all {
		genSpecs = append(genSpecs, genSpecsForAllStructs()...)
		return true
	}

	return false
}

var structRegex = regexp.MustCompile(`(\*?)(\p{L}+)\.(\p{L}+)`)

func genSpecFromStructArg(arg string) (*genSpec, bool) {
	matches := structRegex.FindStringSubmatch(arg)

	if matches == nil {
		return nil, false
	}

	ptr := matches[1]
	pkg := matches[2]
	typ := matches[3]

	return newGenSpec(ptr, pkg, typ), true
}

var allRegex = regexp.MustCompile(`-(\*?)a(ll)?`)

var goFiles = func(f os.FileInfo) bool {
	return strings.HasSuffix(f.Name(), ".go")
}

func genSpecsForAllStructs() (g []*genSpec) {
	fset := token.NewFileSet() // positions are relative to fset

	dir, err := parser.ParseDir(fset, "./", goFiles, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}

	for pkg, f := range dir {
		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.TypeSpec:
				switch y := x.Type.(type) {
				case *ast.StructType:
					_ = y
					typ := x.Name.String()
					g = append(g, newGenSpec("*", pkg, typ))
				}
			}
			return true
		})
	}
	return
}

func writeFile(t *template.Template, genSpecs []*genSpec) {
	for _, v := range genSpecs {
		f, err := os.Create(v.FileName)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		t.Execute(f, v)
		fmt.Printf("  generated %s, yay!\n", v)
	}
}

const usage = `Usage: gen [[*]package.TypeName] [-all]

  *package.TypeName    # generate funcs for specified struct type; use leading * to specify pointer type (recommended)
  -all                 # generates all structs in current directory; shortcut -a
`
