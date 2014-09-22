package slice

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"testing"

	"github.com/clipperhouse/gen/typewriter"
)

var typs []typewriter.Type

func init() {
	pkg := typewriter.NewPackage("dummy", "SomePackage")

	t1, err := pkg.Eval("int")

	if err != nil {
		panic(err)
	}

	t2, err := pkg.Eval("*rune")

	if err != nil {
		panic(err)
	}

	t1.Tags = typewriter.Tags{
		typewriter.Tag{
			Name: "slice",
			Values: []typewriter.TagValue{
				{"GroupBy", []typewriter.Type{t2}},
				{"Where", nil},
			},
		},
	}

	typs = append(typs, t1)
}

func TestValidate(t *testing.T) {
	sw := NewSliceWriter()

	for _, typ := range typs {
		if sw.validated[typ.String()] {
			t.Errorf("type should not show having been validated yet")
		}

		write, err := sw.Validate(typ)

		if !write {
			t.Errorf("type should write")
		}

		if err != nil {
			t.Errorf("type should be valid; got error '%s'", err)
		}

		if !sw.validated[typ.String()] {
			t.Errorf("type should show having been validated")
		}

		if _, ok := sw.caches[typ.String()]; !ok {
			t.Errorf("type should appear in sw.caches")
		}
	}
}

func TestWrite(t *testing.T) {
	for _, typ := range typs {
		var b bytes.Buffer

		sw := NewSliceWriter()

		write, err := sw.Validate(typ)

		if err != nil {
			t.Error(err)
		}

		if !write {
			t.Errorf("should write %s", typ)
		}

		b.WriteString(fmt.Sprintf("package %s\n", typ.Package.Name()))
		sw.WriteBody(&b, typ)

		src := b.String()

		fmt.Println(src)

		fset := token.NewFileSet()
		if _, err := parser.ParseFile(fset, "testwrite.go", src, 0); err != nil {
			t.Error(err)
		}
	}
}
