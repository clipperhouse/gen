package container

import (
	"testing"

	"github.com/clipperhouse/gen/typewriter"
)

func TestValidate(t *testing.T) {
	g := NewContainerWriter()

	pkg := typewriter.NewPackage("dummy", "SomePackage")

	typ := typewriter.Type{
		Package: pkg,
		Name:    "SomeType",
		Tags:    typewriter.Tags{},
	}

	write, err := g.Validate(typ)

	if write {
		t.Errorf("no 'containers' tag should not write")
	}

	if err != nil {
		t.Error(err)
	}

	typ2 := typewriter.Type{
		Package: pkg,
		Name:    "SomeType2",
		Tags: typewriter.Tags{
			typewriter.Tag{
				Name:  "containers",
				Items: []string{},
			},
		},
	}

	write2, err2 := g.Validate(typ2)

	if write2 {
		t.Errorf("empty 'containers' tag should not write")
	}

	if err2 != nil {
		t.Error(err)
	}

	typ3 := typewriter.Type{
		Package: pkg,
		Name:    "SomeType3",
		Tags: typewriter.Tags{
			typewriter.Tag{
				Name:  "containers",
				Items: []string{"List", "Foo"},
			},
		},
	}

	write3, err3 := g.Validate(typ3)

	if !write3 {
		t.Errorf("'containers' tag with List should write (and ignore others)")
	}

	if err3 != nil {
		t.Error(err)
	}
}
