package typewriter

import (
	"fmt"
	"io"
	"sort"
	"text/template"
)

type TypeWriter interface {
	Validate(t Type) error
	Imports(t Type) []string
	Write(w io.Writer, t Type)
}

var TypeWriters = make(map[string]TypeWriter)

// Register allows template packages to make themselves known to a 'parent' package, usually in the init() func.
// Comparable to the approach taken by builtin image package for registration of image types (eg image/png).
// Your program will do something like:
//	import (
//		"github.com/clipperhouse/gen/templates"
//		_ "github.com/clipperhouse/gen/templates/projection"
//	)
func Register(name string, tw TypeWriter) {
	TypeWriters[name] = tw
}

// TODO: idea, write into a memory buffer by default, validate the ast before committing to passed writer (which is likely a file)
// It's a reasonable design principle not to break an existing codebase on disk, so validation before committment is desirable
// Possibly allow override?

// Write writes the generated code for all registered TypeWriters. Typically, this is a file.
func Write(w io.Writer, t Type) {
	for _, tw := range TypeWriters {
		err := tw.Validate(t)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	for _, tw := range TypeWriters {
		writeHeader(w, t)
		tw.Write(w, t)
	}
}

// writeHeader writes the package and import declarations for a type
func writeHeader(w io.Writer, t Type) {
	// unique imports across all TypeWriters
	uniq := make(map[string]struct{})
	for _, tw := range TypeWriters {
		for _, imp := range tw.Imports(t) {
			uniq[imp] = struct{}{}
		}
	}

	fmt.Println(uniq)

	// make it into a slice for the template
	imports := make([]string, 0)
	for imp := range uniq {
		imports = append(imports, imp)
	}

	// approximate the conventional ordering of imports :)
	sort.Strings(imports)

	headerTmpl.Execute(w, t.Package.Name)
	importsTmpl.Execute(w, imports)
}

var headerTmpl = template.Must(template.New("header").Parse(`// This file was auto-generated using github.com/clipperhouse/gen
// Modifying this file is not recommended as it will likely be overwritten in the future

package {{.}}
`))

var importsTmpl = template.Must(template.New("imports").Parse(`{{if gt (len .) 0}}
import ({{range .}}
	"{{.}}"{{end}}
)
{{end}}
`))
