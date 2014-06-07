package typewriter

import (
	// "fmt"
	// "os"
	"testing"
)

func TestGetTypes(t *testing.T) {
	// app and Package types are marked up with +test
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

	if _, found := m["Package"]; !found {
		t.Errorf("should have found the Package type")
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

	if len(m["Package"].Tags) != 1 {
		t.Errorf("typ should have 1 tag, found %v", len(m["Package"].Tags))
	}

	if len(m["Package"].Tags[0].Items) != 1 {
		t.Errorf("Tag should have 1 Item, found %v", len(m["Package"].Tags[0].Items))
	}

	typs2, err2 := getTypes("+dummy", nil)

	if len(typs2) != 0 {
		t.Errorf("should have no marked-up types for +dummy")
	}

	if err2 != nil {
		t.Error(err2)
	}

	// filter := func(f os.FileInfo) bool {
	// 	return f.Name() == "tag.go"
	// }

	// typs3, err3 := getTypes("+gen", filter)

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
