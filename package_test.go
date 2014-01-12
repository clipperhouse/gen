package main

import (
	"strings"
	"testing"
)

type test struct{}

// gen:"Where,Single"
type test2 int

// project:"string,test"
type test3 test2

// gen:"Count,Distinct" project:"test2,float64"
type test4 test

// gen:"Count,Distinct,Count" project:"test2,test2,float64"
type test5 string

var packages map[string]*Package

func TestSetup(t *testing.T) {
	packages = getPackages() // use gen itself :)
}

func TestTagParsing(t *testing.T) {
	pkg := packages["main"]

	typeArg1 := &typeArg{"", "main", "test"}
	type1, err1 := pkg.GetType(typeArg1)

	if err1 != nil || type1 == nil {
		t.Errorf("should have been found %s without error", typeArg1)
	}

	if len(type1.SubsettedMethods) > 0 {
		t.Errorf("type %s should have no subsetted methods", type1.String())
	}

	if len(type1.ProjectedTypes) > 0 {
		t.Errorf("type %s should have no projected types", type1.String())
	}

	typeArg2 := &typeArg{"", "main", "test2"}
	type2, err2 := pkg.GetType(typeArg2)

	if err2 != nil || type2 == nil {
		t.Errorf("should have been found %s without error", typeArg2)
	}

	if len(type2.SubsettedMethods) == 0 {
		t.Errorf("type %s should have subsetted methods", type2.String())
	}

	if len(type2.ProjectedTypes) > 0 {
		t.Errorf("type %s should have no projected types but has %v", type2.String(), type2.ProjectedTypes)
	}

	typeArg3 := &typeArg{"", "main", "test3"}
	type3, err3 := pkg.GetType(typeArg3)

	if err3 != nil || type3 == nil {
		t.Errorf("should have been found %s without error", typeArg3)
	}

	if len(type3.SubsettedMethods) > 0 {
		t.Errorf("type %s should have no subsetted methods but has %v", type3.String(), type3.SubsettedMethods)
	}

	if len(type3.ProjectedTypes) == 0 {
		t.Errorf("type %s should have projected types", type3.String())
	}

	typeArg4 := &typeArg{"", "main", "test4"}
	type4, err4 := pkg.GetType(typeArg4)

	if err4 != nil || type4 == nil {
		t.Errorf("should have been found %s without error", typeArg4)
	}

	if len(type4.SubsettedMethods) == 0 {
		t.Errorf("type %s should have subsetted methods", type4.String())
	}

	if len(type4.ProjectedTypes) == 0 {
		t.Errorf("type %s should have projected types", type4.String())
	}
}

func TestTagParsingErrors(t *testing.T) {
	pkg := packages["main"]

	typeArg1 := &typeArg{"", "main", "dummy"}
	type1, err1 := pkg.GetType(typeArg1)

	if err1 == nil {
		t.Errorf("should have returned 1 error for an unknown type %s", type1)
	}
}

func TestEval(t *testing.T) {
	real := packages["main"]
	typ, err := real.Eval("test")

	if err != nil {
		t.Errorf("valid type %s should Eval", "test")
	}

	if typ == nil {
		t.Errorf("valid type %v should not be nil", typ)
	}

	typ2, err := real.Eval("dummy")

	if err == nil {
		t.Errorf("invalid type %s should fail to Eval", "test")
	}

	if typ2 != nil {
		t.Errorf("invalid type %v should Eval to nil", typ2)
	}

	typ3, err := real.Eval("*test")

	if err != nil {
		t.Errorf("pointer type %s should Eval", "*test")
	}

	if typ3 == nil {
		t.Errorf("pointer type %s should Eval", "*test")
	}

	if !strings.HasPrefix(typ3.String(), "*") {
		t.Errorf("type %s should Eval to pointer type", "*test")
	}

	fake := &Package{}

	typ4, err := fake.Eval("test")

	if err == nil {
		t.Errorf("valid named type %s should fail to Eval for invalid package", "test")
	}

	if typ4 != nil {
		t.Errorf("named valid type %s should fail to Eval for invalid package", "test")
	}

	typ5, err := fake.Eval("float64")

	if err != nil {
		t.Errorf("valid builtin type %s should Eval (Universe scope) for invalid package", "int")
	}

	if typ5 == nil {
		t.Errorf("valid builtin type %s should Eval (Universe scope) for invalid package", "int")
	}
}
