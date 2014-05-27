package genwriter

import (
	"bytes"
	"code.google.com/p/go.tools/go/types"
	"fmt"
	"github.com/clipperhouse/typewriter"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	g := GenWriter{}
	pkg := &typewriter.Package{
		types.NewPackage("dummy", "SomePackage"),
	}

	typ := typewriter.Type{
		Package: pkg,
		Name:    "SomeType",
		Tags:    typewriter.Tags{},
	}

	if validated[typ.String()] {
		t.Errorf("type should not show having been validated yet")
	}

	if err := ensureValidation(typ); err == nil {
		t.Errorf("ensure validation should return err prior to validation")
	}

	valid, err := g.Validate(typ)

	if !valid || err != nil {
		t.Errorf("type should be valid")
	}

	if !validated[typ.String()] {
		t.Errorf("type should show having been validated")
	}

	if err := ensureValidation(typ); err != nil {
		t.Errorf("ensure validation should not return err after validation")
	}

	if _, ok := models[typ.String()]; !ok {
		t.Errorf("type should appear in models")
	}

	if m := models[typ.String()]; len(m.methods) == 0 {
		t.Errorf("model without tags should have methods")
	}

	if m := models[typ.String()]; len(m.projections) != 0 {
		t.Errorf("model without tags should have no projections")
	}

	typ2 := typewriter.Type{
		Package: pkg,
		Name:    "SomeType2",
		Tags: typewriter.Tags{
			typewriter.Tag{
				Name:  "projections",
				Items: []string{"int", "string"},
			},
		},
	}

	valid2, err2 := g.Validate(typ2)

	if !valid2 || err2 != nil {
		t.Errorf("type should be valid")
	}

	if m := models[typ2.String()]; len(m.projections) == 0 {
		t.Errorf("model with projections tag should have projections")
	}

	typ3 := typewriter.Type{
		Package: pkg,
		Name:    "SomeType3",
		Tags: typewriter.Tags{
			typewriter.Tag{
				Name:  "projections",
				Items: []string{"int", "Foo"},
			},
		},
	}

	valid3, err3 := g.Validate(typ3)

	if valid3 {
		t.Errorf("type with unknown projection type should be invalid")
	}

	if err3 == nil {
		t.Errorf("type with unknown projection should return error")
	}

	if !strings.Contains(err3.Error(), "Foo") {
		t.Errorf("type with unknown projection type should mention the unknown projection type; got %v", err3)
	}

	if !strings.Contains(err3.Error(), typ3.Name) {
		t.Errorf("type with unknown projection type should mention the type on which it was declared; got %v", err3)
	}

	if m := models[typ2.String()]; len(m.projections) == 0 {
		t.Errorf("model with projections tag should have projections")
	}
}

func TestWriteHeader(t *testing.T) {
	var b bytes.Buffer

	g := GenWriter{}
	pkg := &typewriter.Package{
		types.NewPackage("dummy", "SomePackage"),
	}

	typ := typewriter.Type{
		Package: pkg,
		Name:    "SomeType",
		Tags: typewriter.Tags{
			{
				Name:  "methods",
				Items: []string{"All"},
			}, // subset to ensure no Sort
		},
	}

	g.Validate(typ)
	g.WriteHeader(&b, typ)

	if strings.Contains(b.String(), "Copyright") {
		t.Errorf("should not contain license info if no sort; got:\n%s", b.String())
	}

	var b2 bytes.Buffer

	typ2 := typewriter.Type{
		Package: pkg,
		Name:    "SomeType",
		Tags:    typewriter.Tags{}, // default includes sort
	}

	g.Validate(typ2)
	g.WriteHeader(&b2, typ2)

	if !strings.Contains(b2.String(), "Copyright") {
		t.Errorf("should contain license info if sort; got:\n%s", b2.String())
	}
}

func TestImports(t *testing.T) {
	g := GenWriter{}
	pkg := &typewriter.Package{
		types.NewPackage("dummy", "SomePackage"),
	}

	typ := typewriter.Type{
		Package: pkg,
		Name:    "SomeType",
		Tags:    typewriter.Tags{},
	}

	g.Validate(typ)
	imports := g.Imports(typ)

	if !sliceContains(imports, "errors") {
		t.Errorf("imports should include 'errors'")
	}

	if len(imports) > 1 {
		t.Errorf("imports should only include 'errors'")
	}

	// not easy to devise a test for Sort, since it requires the unexported typewriter.Type.ordered
}

func TestWrite(t *testing.T) {
	g := GenWriter{}

	pkg := &typewriter.Package{
		types.NewPackage("dummy", "SomePackage"),
	}

	typs := []typewriter.Type{
		{
			Package: pkg,
			Name:    "FirstType",
			Tags:    typewriter.Tags{}, // simple default
		},
		{
			Package: pkg,
			Name:    "SecondType",
			Tags: typewriter.Tags{
				{
					Name:  "projections",
					Items: []string{"int", "string"}, // with projections
				},
			},
		},
		{
			Package: pkg,
			Name:    "ThirdType",
			Tags: typewriter.Tags{
				{
					Name:  "methods",
					Items: []string{"All", "Any"}, // subsetted
				},
			},
		},
		{
			Package: pkg,
			Name:    "FourthType",
			Tags: typewriter.Tags{
				{
					Name:  "methods",
					Items: []string{"All", "Any"}, // subsetted
				},
				{
					Name:  "projections",
					Items: []string{"int", "string"}, // and projections
				},
			},
		},
		{
			Package: pkg,
			Name:    "FifthType",
			Tags: typewriter.Tags{
				{
					Name:    "methods",
					Items:   []string{"Count", "Where"},
					Negated: true,
				},
			},
		},
	}

	for _, typ := range typs {
		var b bytes.Buffer

		g.Validate(typ)

		b.WriteString(fmt.Sprintf("package %s\n", pkg.Name()))

		g.Write(&b, typ)

		src := b.String()

		fset := token.NewFileSet()

		_, err := parser.ParseFile(fset, "testwrite.go", src, 0)

		if err != nil {
			t.Errorf(err.Error())
		}
	}
}

func sliceContains(a []string, s string) bool {
	for _, v := range a {
		if v == s {
			return true
		}
	}
	return false
}
