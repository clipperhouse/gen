package main

import (
	"testing"
)

func TestGenSpecParsing(t *testing.T) {
	s1 := `// Here is a description of some type
// gen that may span lines`
	_, found1 := getGenSpec(s1)

	if found1 {
		t.Errorf("no gen spec should have been found")
	}

	s2 := `// Here is a description of some type
// +gen`
	spec2, found2 := getGenSpec(s2)

	if !found2 {
		t.Errorf("gen spec should have been found")
	}

	if spec2 == nil {
		t.Errorf("gen spec should not be nil")
	}

	if len(spec2.Pointer) > 0 {
		t.Errorf("gen spec should not be pointer by default")
	}

	s3 := `// Here is a description of some type
// +gen *`
	spec3, found3 := getGenSpec(s3)

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
	spec4, found4 := getGenSpec(s4)

	if !found4 {
		t.Errorf("gen spec should have been found")
	}

	if spec4 == nil {
		t.Errorf("gen spec should not be nil")
	}

	if spec4.Pointer != "*" {
		t.Errorf("gen spec should be pointer")
	}

	if len(spec4.SubsettedMethods) != 2 {
		t.Errorf("gen spec should have 2 subsetted methods")
	}

	if len(spec4.ProjectedTypes) > 0 {
		t.Errorf("gen spec should no projected types")
	}

	s5 := `// Here is a description of some type
// +gen methods:"Any,All" projections:"GroupBy"`
	spec5, found5 := getGenSpec(s5)

	if !found5 {
		t.Errorf("gen spec should have been found")
	}

	if spec5 == nil {
		t.Errorf("gen spec should not be nil")
	}

	if len(spec5.Pointer) > 0 {
		t.Errorf("gen spec should not be pointer")
	}

	if len(spec5.SubsettedMethods) != 2 {
		t.Errorf("gen spec should have 2 subsetted methods")
	}

	if len(spec5.ProjectedTypes) != 1 {
		t.Errorf("gen spec should have 1 projected type")
	}
}
