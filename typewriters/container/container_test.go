package container

import (
	"bytes"
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

	var b bytes.Buffer
	err := g.WriteHeader(&b, typ)

	if b.Len() > 0 {
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
				Name:   "containers",
				Values: []typewriter.TagValue{},
			},
		},
	}

	var b2 bytes.Buffer
	err := g.WriteBody(&b2, typ2)

	if b2.Len() > 0 {
		t.Errorf("empty 'containers' tag should not write")
	}

	if err != nil {
		t.Error(err)
	}

	typ3 := typewriter.Type{
		Package: pkg,
		Name:    "SomeType3",
		Tags: typewriter.Tags{
			typewriter.Tag{
				Name: "containers",
				Values: []typewriter.TagValue{
					{"List", nil},
				},
			},
		},
	}

	var b3 bytes.Buffer
	err := g.WriteBody(&b3, typ3)

	if b3.Len() == 0 {
		t.Errorf("'containers' tag with List should write (and ignore others)")
	}

	if err != nil {
		t.Error(err)
	}
}
