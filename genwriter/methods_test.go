package genwriter

import (
	"github.com/clipperhouse/typewriter"
	"testing"
)

func TestEvaluateTags(t *testing.T) {
	typ := typewriter.Type{
		Name:    "TestType",
		Package: typewriter.NewPackage("dummy", "TestPackage"),
	}

	typ.Tags = typewriter.Tags{}

	standardMethods1, projectionMethods1, err1 := evaluateTags(typ)

	if err1 != nil {
		t.Errorf("empty methods should be ok, instead got '%v'", err1)
	}

	if len(standardMethods1) != len(standardTemplates.GetAllKeys()) {
		t.Errorf("standard methods should default to all")
	}

	if len(projectionMethods1) != 0 {
		t.Errorf("projection methods without projected type should be none, instead got %v", projectionMethods1)
	}

	typ.Tags = typewriter.Tags{
		{
			Name:  "methods",
			Items: []string{"Count", "Where"},
		},
	}

	standardMethods2, projectionMethods2, err2 := evaluateTags(typ)

	if err2 != nil {
		t.Errorf("empty methods should be ok, instead got %v", err2)
	}

	if len(standardMethods2) != 2 {
		t.Errorf("standard methods should be parsed")
	}

	if len(projectionMethods2) != 0 {
		t.Errorf("projection methods without projected typs should be none")
	}

	typ.Tags = typewriter.Tags{
		{
			Name:  "methods",
			Items: []string{"Count", "Unknown"},
		},
	}

	standardMethods3, projectionMethods3, err3 := evaluateTags(typ)

	if err3 == nil {
		t.Errorf("unknown method should be error")
	}

	if len(standardMethods3) != 1 {
		t.Errorf("standard methods (except unknown) should be 1, got %v", len(standardMethods3))
	}

	if len(projectionMethods3) != 0 {
		t.Errorf("projection methods without projected types should be none")
	}

	typ.Tags = typewriter.Tags{
		{
			Name:  "projections",
			Items: []string{"SomeType"},
		},
	}

	standardMethods4, projectionMethods4, err4 := evaluateTags(typ)

	if err4 != nil {
		t.Errorf("projected types without subsetted methods should be ok, instead got: '%v'", err4)
	}

	if len(standardMethods4) != len(standardTemplates.GetAllKeys()) {
		t.Errorf("standard methods should default to all")
	}

	if len(projectionMethods4) != len(projectionTemplates.GetAllKeys()) {
		t.Errorf("projection methods should default to all in presence of projected types")
	}

	typ.Tags = typewriter.Tags{
		{
			Name:  "methods",
			Items: []string{"GroupBy"},
		},
		{
			Name:  "projections",
			Items: []string{"SomeType"},
		},
	}

	standardMethods5, projectionMethods5, err5 := evaluateTags(typ)

	if err5 != nil {
		t.Errorf("projected types with subsetted methods should be ok, instead got: '%v'", err5)
	}

	if len(standardMethods5) != 0 {
		t.Errorf("standard methods should be none")
	}

	if len(projectionMethods5) != 1 {
		t.Errorf("projection methods should be subsetted")
	}

	typ.Tags = typewriter.Tags{
		{
			Name: "methods",
		},
	}

	standardMethods6, projectionMethods6, err6 := evaluateTags(typ)

	if err6 != nil {
		t.Errorf("empty subsetted methods should be ok, instead got: '%v'", err6)
	}

	if len(standardMethods6) != 0 {
		t.Errorf("standard methods should be empty when the tag is empty")
	}

	if len(projectionMethods6) != 0 {
		t.Errorf("projection methods should be none")
	}

	typ.Tags = typewriter.Tags{
		{
			Name:    "methods",
			Items:   []string{"Sort", "Any"},
			Negated: true,
		},
	}

	standardMethods7, projectionMethods7, err7 := evaluateTags(typ)

	if err7 != nil {
		t.Errorf("subsetted methods should be ok, instead got: '%v'", err7)
	}

	expected7 := []string{"All", "Count", "Distinct", "DistinctBy", "Each", "First", "IsSorted", "IsSortedBy", "IsSortedByDesc", "IsSortedDesc", "Max", "MaxBy", "Min", "MinBy", "Single", "SortBy", "SortByDesc", "SortDesc", "Where"}
	if !sliceEqual(standardMethods7, expected7) {
		t.Errorf("standard methods should be negatively subsetted, expected %v, got %v", expected7, standardMethods7)
	}

	if len(projectionMethods7) != 0 {
		t.Errorf("projection methods should be none")
	}

	typ.Tags = typewriter.Tags{
		{
			Name:    "methods",
			Items:   []string{"Sort", "Where", "GroupBy"},
			Negated: true,
		},
		{
			Name:  "projections",
			Items: []string{"int"},
		},
	}

	standardMethods8, projectionMethods8, err8 := evaluateTags(typ)

	if err8 != nil {
		t.Errorf("subsetted methods should be ok, instead got: '%v'", err8)
	}

	expectedStd8 := []string{"All", "Any", "Count", "Distinct", "DistinctBy", "Each", "First", "IsSorted", "IsSortedBy", "IsSortedByDesc", "IsSortedDesc", "Max", "MaxBy", "Min", "MinBy", "Single", "SortBy", "SortByDesc", "SortDesc"}
	if !sliceEqual(standardMethods8, expectedStd8) {
		t.Errorf("standard methods should be negatively subsetted, expected %v, got %v", expectedStd8, standardMethods8)
	}

	expectedPrj8 := []string{"Aggregate", "Average", "Max", "Min", "Select", "Sum"}
	if !sliceEqual(projectionMethods8, expectedPrj8) {
		t.Errorf("projection methods should be negatively subsetted, expected %v, got %v", expectedPrj8, projectionMethods8)
	}
}

func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
