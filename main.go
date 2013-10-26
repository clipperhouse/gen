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

type Values struct {
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

func (v Values) String() string {
	return fmt.Sprintf("%s.%s", v.Package, v.Plural)
}

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

	var values []*Values

	for _, arg := range os.Args[1:] {
		v, valid := getValuesFromArg(arg)
		if valid {
			values = append(values, v...)
		} else {
			fmt.Printf("Invalid argument: %s\n", arg)
			return
		}
	}

	t := getTemplate()
	writeFile(t, values)
}

func newValues(ptr, pkg, typ string) *Values {
	typ = inflect.Singularize(typ)
	return &Values{
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

func getValuesFromArg(arg string) (values []*Values, valid bool) {
	value, success := getStructFromArg(arg)
	if success {
		return append(values, value), true
	}

	all := allRegex.MatchString(arg)
	if all {
		return append(values, getAllStructs()...), true
	}

	return nil, false
}

var structRegex = regexp.MustCompile(`(\*?)(\p{L}+)\.(\p{L}+)`)

func getStructFromArg(arg string) (*Values, bool) {
	matches := structRegex.FindStringSubmatch(arg)

	if matches == nil {
		return nil, false
	}

	ptr := matches[1]
	pkg := matches[2]
	typ := matches[3]

	return newValues(ptr, pkg, typ), true
}

var allRegex = regexp.MustCompile(`-(\*?)a(ll)?`)

var goFiles = func(f os.FileInfo) bool {
	return strings.HasSuffix(f.Name(), ".go")
}

func getAllStructs() (v []*Values) {
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
					v = append(v, newValues("*", pkg, typ))
				}
			}
			return true
		})
	}
	return
}

func writeFile(t *template.Template, values []*Values) {
	for _, v := range values {
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
