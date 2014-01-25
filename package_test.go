package main

import (
	"testing"
)

func TestGenSpecParsing(t *testing.T) {
	dummy := "dummy"

	s1 := `// Here is a description of some type
// gen that may span lines`
	_, found1 := getGenSpec(s1, dummy)

	if found1 {
		t.Errorf("no gen spec should have been found")
	}

	s2 := `// Here is a description of some type
// +gen`
	spec2, found2 := getGenSpec(s2, dummy)

	if !found2 {
		t.Errorf("gen spec should have been found")
	}

	if spec2 == nil {
		t.Errorf("gen spec should not be nil")
	}

	if len(spec2.Pointer) > 0 {
		t.Errorf("gen spec should not be pointer by default")
	}

	if spec2.Methods != nil {
		t.Errorf("gen spec methods should be nil if unspecified")
	}

	if spec2.Projections != nil {
		t.Errorf("gen spec methods should be nil if unspecified")
	}

	s3 := `// Here is a description of some type
// +gen *`
	spec3, found3 := getGenSpec(s3, dummy)

	if !found3 {
		t.Errorf("gen spec should have been found")
	}

	if spec3 == nil {
		t.Errorf("gen spec should not be nil")
	}

	if spec3.Pointer != "*" {
		t.Errorf("gen spec should be pointer")
	}

	s4 := `// Here is a description of some type
// +gen * methods:"Any,All"`
	spec4, found4 := getGenSpec(s4, dummy)

	if !found4 {
		t.Errorf("gen spec should have been found")
	}

	if spec4 == nil {
		t.Errorf("gen spec should not be nil")
	}

	if spec4.Pointer != "*" {
		t.Errorf("gen spec should be pointer")
	}

	if len(spec4.Methods.Items) != 2 {
		t.Errorf("gen spec should have 2 methods")
	}

	if spec4.Projections != nil {
		t.Errorf("gen spec projections should be nil if unspecified")
	}

	s5 := `// Here is a description of some type
// +gen methods:"Any,All" projections:"GroupBy"`
	spec5, found5 := getGenSpec(s5, dummy)

	if !found5 {
		t.Errorf("gen spec should have been found")
	}

	if spec5 == nil {
		t.Errorf("gen spec should not be nil")
	}

	if len(spec5.Pointer) > 0 {
		t.Errorf("gen spec should not be pointer")
	}

	if len(spec5.Methods.Items) != 2 {
		t.Errorf("gen spec should have 2 subsetted methods")
	}

	if len(spec5.Projections.Items) != 1 {
		t.Errorf("gen spec should have 1 projected type")
	}

	s6 := `// Here is a description of some type
// +gen methods:"" projections:""`
	spec6, found6 := getGenSpec(s6, dummy)

	if !found6 {
		t.Errorf("gen spec should have been found")
	}

	if spec6 == nil {
		t.Errorf("gen spec should not be nil")
	}

	if len(spec6.Pointer) > 0 {
		t.Errorf("gen spec should not be pointer")
	}

	if spec6.Methods == nil {
		t.Errorf("gen spec methods should exist even if empty")
	}

	if len(spec6.Methods.Items) > 0 {
		t.Errorf("gen spec methods should be empty, instead got %v", len(spec6.Methods.Items))
	}

	if spec6.Projections == nil {
		t.Errorf("gen spec projections should exist even if empty")
	}

	if len(spec6.Projections.Items) > 0 {
		t.Errorf("gen spec projections should be empty, instead got %v", spec6.Projections.Items)
	}
}

func TestMethodDetermination(t *testing.T) {
	dummy := "dummy"

	spec1 := &GenSpec{"", dummy, nil, nil}

	standardMethods1, projectionMethods1, err1 := determineMethods(spec1)

	if err1 != nil {
		t.Errorf("empty methods should be ok, instead got '%v'", err1)
	}

	if len(standardMethods1) != len(getStandardMethodKeys()) {
		t.Errorf("standard methods should default to all")
	}

	if len(projectionMethods1) != 0 {
		t.Errorf("projection methods without projected type should be none, instead got %v", projectionMethods1)
	}

	spec2 := &GenSpec{"", dummy, &GenTag{[]string{"Count", "Where"}}, nil}

	standardMethods2, projectionMethods2, err2 := determineMethods(spec2)

	if err2 != nil {
		t.Errorf("empty methods should be ok, instead got %v", err2)
	}

	if len(standardMethods2) != 2 {
		t.Errorf("standard methods should be parsed")
	}

	if len(projectionMethods2) != 0 {
		t.Errorf("projection methods without projected types should be none")
	}

	spec3 := &GenSpec{"", dummy, &GenTag{[]string{"Count", "Unknown"}}, &GenTag{[]string{}}}

	standardMethods3, projectionMethods3, err3 := determineMethods(spec3)

	if err3 == nil {
		t.Errorf("unknown type should be error")
	}

	if len(standardMethods3) != 1 {
		t.Errorf("standard methods should be parsed, minus unknown")
	}

	if len(projectionMethods3) != 0 {
		t.Errorf("projection methods without projected types should be none")
	}

	spec4 := &GenSpec{"", dummy, nil, &GenTag{[]string{"SomeType"}}}

	standardMethods4, projectionMethods4, err4 := determineMethods(spec4)

	if err4 != nil {
		t.Errorf("projected types without subsetted methods should be ok, instead got: '%v'", err4)
	}

	if len(standardMethods4) != len(getStandardMethodKeys()) {
		t.Errorf("standard methods should default to all")
	}

	if len(projectionMethods4) != len(getProjectionMethodKeys()) {
		t.Errorf("projection methods should default to all in presence of projected types")
	}

	spec5 := &GenSpec{"", dummy, &GenTag{[]string{"GroupBy"}}, &GenTag{[]string{"SomeType"}}}

	standardMethods5, projectionMethods5, err5 := determineMethods(spec5)

	if err5 != nil {
		t.Errorf("projected types with subsetted methods should be ok, instead got: '%v'", err5)
	}

	if len(standardMethods5) != 0 {
		t.Errorf("standard methods should be none")
	}

	if len(projectionMethods5) != 1 {
		t.Errorf("projection methods should be subsetted")
	}

	spec6 := &GenSpec{"", dummy, &GenTag{[]string{}}, nil}

	standardMethods6, projectionMethods6, err6 := determineMethods(spec6)

	if err6 != nil {
		t.Errorf("empty subsetted methods should be ok, instead got: '%v'", err6)
	}

	if len(standardMethods6) != 0 {
		t.Errorf("standard methods should be empty when the tag is empty")
	}

	if len(projectionMethods6) != 0 {
		t.Errorf("projection methods should be none")
	}
}

// +gen
type Thing1 int

type Thing2 Thing1

// +gen * methods:"Any,Where"
type Thing3 float64

// +gen projections:"int,Thing2"
type Thing4 struct{}

// +gen methods:"Count,GroupBy,Select,Aggregate" projections:"string,Thing4"
type Thing5 Thing4

func TestStandardMethods(t *testing.T) {
	packages := getPackages()

	if len(packages) != 1 {
		t.Errorf("should find one package")
	}

	// put into a map for convenience
	types := make(map[string]*Type)
	for _, typ := range packages[0].Types {
		types[typ.Name] = typ
	}

	thing1, ok1 := types["Thing1"]

	if !ok1 || thing1 == nil {
		t.Errorf("Thing1 should have been identified as a gen Type")
	}

	if len(thing1.Pointer) != 0 {
		t.Errorf("Thing1 should not generate pointers")
	}

	if len(thing1.StandardMethods) != len(StandardTemplates) {
		t.Errorf("Thing1 should have all standard methods")
	}

	if len(thing1.Projections) != 0 {
		t.Errorf("Thing1 should have no projections")
	}

	thing2, ok2 := types["Thing2"]

	if ok2 || thing2 != nil {
		t.Errorf("Thing2 should not have been identified as a gen Type")
	}

	thing3 := types["Thing3"]

	if thing3.Pointer != "*" {
		t.Errorf("Thing3 should generate pointers")
	}

	if len(thing3.StandardMethods) != 2 {
		t.Errorf("Thing3 should have subsetted Any and Where, but has: %v", thing3.StandardMethods)
	}

	if len(thing1.Projections) != 0 {
		t.Errorf("Thing3 should have no projections, but has: %v", thing3.Projections)
	}

	thing4 := types["Thing4"]

	if len(thing4.Projections) != 2*len(ProjectionTemplates) {
		t.Errorf("Thing4 should have all projection methods for 2 types, but has: %v", thing4.Projections)
	}

	thing5 := types["Thing5"]

	if len(thing5.StandardMethods) != 1 {
		t.Errorf("Thing5 should have 1 subsetted standard method, but has: %v", thing5.StandardMethods)
	}

	if len(thing5.Projections) != 2*3 {
		t.Errorf("Thing4 should have 3 subsetted projection methods for 2 types, but has: %v", thing5.Projections)
	}
}
