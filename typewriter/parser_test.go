package typewriter

import (
	// "fmt"
	"os"
	"testing"
)

// +test foo:"bar"
type dummy int

func TestGetTypes(t *testing.T) {
	// app and dummy types are marked up with +test
	typs, err := getTypes("+test", nil)

	if err != nil {
		t.Error(err)
	}

	if len(typs) != 2 {
		t.Errorf("should have found the 2 marked-up types, found %v", len(typs))
	}

	// put 'em into a map for convenience
	m := typeSliceToMap(typs)

	if _, found := m["app"]; !found {
		t.Errorf("should have found the app type")
	}

	if _, found := m["dummy"]; !found {
		t.Errorf("should have found the dummy type")
	}

	dummy := m["dummy"]

	if !dummy.Comparable() {
		t.Errorf("dummy type should be comparable")
	}

	if !dummy.Ordered() {
		t.Errorf("dummy type should be ordered")
	}

	if !dummy.Numeric() {
		t.Errorf("dummy type should be numeric")
	}

	// check tag existence at a high level here, see also tag parsing tests

	if len(m["app"].Tags) != 2 {
		t.Errorf("typ should have 2 Tags, found %v", len(m["app"].Tags))
	}

	if len(m["app"].Tags[0].Items) != 1 {
		t.Errorf("Tag should have 1 Item, found %v", len(m["app"].Tags[0].Items))
	}

	if len(m["app"].Tags[1].Items) != 2 {
		t.Errorf("Tag should have 2 Items, found %v", len(m["app"].Tags[1].Items))
	}

	if len(m["dummy"].Tags) != 1 {
		t.Errorf("typ should have 1 tag, found %v", len(m["dummy"].Tags))
	}

	if len(m["dummy"].Tags[0].Items) != 1 {
		t.Errorf("Tag should have 1 Item, found %v", len(m["dummy"].Tags[0].Items))
	}

	// filtered types should not show up

	filter := func(f os.FileInfo) bool {
		return f.Name() != "app.go"
	}

	typs2, err2 := getTypes("+test", filter)

	if err2 != nil {
		t.Error(err2)
	}

	if len(typs2) != 1 {
		t.Errorf("should have found the 1 marked-up type when filtered, found %v", len(typs2))
	}

	m2 := typeSliceToMap(typs2)

	if _, found := m2["app"]; found {
		t.Errorf("should not have found the app type")
	}

	if _, found := m2["dummy"]; !found {
		t.Errorf("should have found the dummy type")
	}

	// no false positives
	typs3, err3 := getTypes("+dummy", nil)

	if len(typs3) != 0 {
		t.Errorf("should have no marked-up types for +dummy")
	}

	if err3 != nil {
		t.Error(err3)
	}
}

func typeSliceToMap(typs []Type) map[string]Type {
	result := make(map[string]Type)
	for _, v := range typs {
		result[v.Name] = v
	}
	return result
}

func tagSliceToMap(typs []Tag) map[string]Tag {
	result := make(map[string]Tag)
	for _, v := range typs {
		result[v.Name] = v
	}
	return result
}
