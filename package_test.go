package main

import (
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
	type1, errs := pkg.GetType(typeArg1)

	if len(errs) > 0 || type1 == nil {
		t.Errorf("should have been found %s without error", typeArg1)
	}

	if len(type1.SubsettedMethods) > 0 {
		t.Errorf("type %s should have no subsetted methods", type1.String())
	}

	if len(type1.ProjectedTypes) > 0 {
		t.Errorf("type %s should have no projected types", type1.String())
	}

	typeArg2 := &typeArg{"", "main", "test2"}
	type2, errs := pkg.GetType(typeArg2)

	if len(errs) > 0 || type2 == nil {
		t.Errorf("should have been found %s without error", typeArg2)
	}

	if len(type2.SubsettedMethods) == 0 {
		t.Errorf("type %s should have subsetted methods", type2.String())
	}

	if len(type2.ProjectedTypes) > 0 {
		t.Errorf("type %s should have no projected types but has %v", type2.String(), type2.ProjectedTypes)
	}

	typeArg3 := &typeArg{"", "main", "test3"}
	type3, errs := pkg.GetType(typeArg3)

	if len(errs) > 0 || type3 == nil {
		t.Errorf("should have been found %s without error", typeArg3)
	}

	if len(type3.SubsettedMethods) > 0 {
		t.Errorf("type %s should have no subsetted methods but has %v", type3.String(), type3.SubsettedMethods)
	}

	if len(type3.ProjectedTypes) == 0 {
		t.Errorf("type %s should have projected types", type3.String())
	}

	typeArg4 := &typeArg{"", "main", "test4"}
	type4, errs := pkg.GetType(typeArg4)

	if len(errs) > 0 || type4 == nil {
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

	typeArg5 := &typeArg{"", "main", "test5"}
	type5, errs := pkg.GetType(typeArg5)

	if len(errs) != 2 {
		t.Errorf("should have been found 2 errors for duplicates on %s", type5)
	}

	typeArg6 := &typeArg{"", "main", "dummy"}
	type6, errs := pkg.GetType(typeArg6)

	if len(errs) != 1 {
		t.Errorf("should have returned 1 error for an unknown type %s", type6)
	}
}

func TestEval(t *testing.T) {
	real := packages["main"]
	typ, err := real.Eval("test")

	if err != nil {
		t.Errorf("valid type %s should Eval", "test")
	}

	if typ == nil {
		t.Errorf("valid type %s should not be nil")
	}

	typ2, err := real.Eval("dummy")

	if err == nil {
		t.Errorf("invalid type %s should fail to Eval", "test")
	}

	if typ2 != nil {
		t.Errorf("invalid type %s should Eval to nil")
	}

	fake := &Package{}

	typ3, err := fake.Eval("test")

	if err == nil {
		t.Errorf("valid named type %s should fail to Eval for invalid package", "test")
	}

	if typ3 != nil {
		t.Errorf("named valid type %s should fail to Eval for invalid package", "test")
	}

	typ4, err := fake.Eval("float64")

	if err != nil {
		t.Errorf("valid builtin type %s should Eval (Universe scope) for invalid package", "int")
	}

	if typ4 == nil {
		t.Errorf("valid builtin type %s should Eval (Universe scope) for invalid package", "int")
	}

}
