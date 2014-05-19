package genwriter

import (
	"code.google.com/p/go.tools/go/types"
	//	"fmt"
	"github.com/clipperhouse/typewriter"
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
		t.Errorf("plain type should be valid")
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
}
