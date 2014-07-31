package typewriter

import (
	"os"
	"strings"
	"testing"
)

func TestParseTags(t *testing.T) {
	doc := `// some stuff that's actually a comment
+test foo:"bar,Baz"
`
	pointer, tags, found, err := parseTags("+test", doc)

	if err != nil {
		t.Error(err)
	}

	if !found {
		t.Errorf("tag should have been found in %s", doc)
	}

	if len(tags) != 1 {
		t.Errorf("one tag should have been found in %s", doc)
	}

	if pointer {
		t.Errorf("pointer should not have been found in %s", doc)
	}

	doc2 := `// some stuff that's actually a comment
+test foo:"bar,Baz" thing:"Stuff, yay"
`
	pointer2, tags2, found2, err2 := parseTags("+test", doc2)

	if err2 != nil {
		t.Error(err2)
	}

	if !found2 {
		t.Errorf("tags should have been found in %s", doc2)
	}

	if len(tags2) != 2 {
		t.Errorf("two tags should have been found in %s; found %v", tags2)
	}

	if tags2[0].Name != "foo" || tags2[1].Name != "thing" {
		t.Errorf("'foo' and 'thing' should have been found in %s; found %v", tags2)
	}

	if len(tags2[0].Items) != 2 || len(tags2[1].Items) != 2 {
		t.Errorf("each tag should have 2 Items")
	}

	if pointer2 {
		t.Errorf("pointer should not have been found in %s", doc2)
	}

	doc3 := `// some stuff that's actually a comment
+test * foo:"bar,Baz" thing:"Stuff, yay" more:"stuff"
`
	pointer3, tags3, found3, err3 := parseTags("+test", doc3)

	if err3 != nil {
		t.Error(err3)
	}

	if !found3 {
		t.Errorf("tags should have been found in %s", doc3)
	}

	if !pointer3 {
		t.Errorf("pointer should have been found in %s", doc3)
	}

	if len(tags3) != 3 {
		t.Errorf("3 tags should have been found in %s; found %v", doc3, len(tags3))
	}

	doc4 := `// some stuff that's actually a comment
+test * foo:"bar,Baz" thing:"Stuff, yay" crap more:"stuff" garbage
`
	_, _, _, err4 := parseTags("+test", doc4)

	if err4 == nil {
		t.Error("invalid syntax should return error")
	}

	doc5 := `// some stuff that's actually a comment
+testfoo:"bar,Baz" thing:"Stuff, yay" crap more:"stuff" garbage
`
	_, _, _, err5 := parseTags("+test", doc5)

	if err5 == nil {
		t.Error("invalid syntax should return error")
	}

	doc6 := `// some stuff that's actually a comment
+test *foo:"bar,Baz" thing:"Stuff, yay" crap more:"stuff" garbage
`
	_, _, _, err6 := parseTags("+test", doc6)

	if err6 == nil {
		t.Error("invalid syntax should return error")
	}
}

func TestGetTypes(t *testing.T) {
	// app and dummy types are marked up with +test
	typs, err := getTypes("+test", nil)

	if err != nil {
		t.Error(err)
	}

	if len(typs) != 4 {
		t.Errorf("should have found the 4 marked-up types, found %v", len(typs))
	}

	// put 'em into a map for convenience
	m := typeSliceToMap(typs)

	if _, found := m["app"]; !found {
		t.Errorf("should have found the app type")
	}

	if _, found := m["dummy"]; !found {
		t.Errorf("should have found the dummy type")
	}

	if _, found := m["dummy2"]; !found {
		t.Errorf("should have found the dummy2 type")
	}

	if _, found := m["dummy3"]; !found {
		t.Errorf("should have found the dummy3 type")
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

	dummy2 := m["dummy2"]

	if dummy2.Comparable() {
		t.Errorf("dummy2 type should not be comparable")
	}

	if dummy2.Ordered() {
		t.Errorf("dummy2 type should not be ordered")
	}

	if dummy2.Numeric() {
		t.Errorf("dummy2 type should not be numeric")
	}

	dummy3 := m["dummy3"]

	if !dummy3.Comparable() {
		t.Errorf("dummy3 type should be comparable")
	}

	if !dummy3.Ordered() {
		t.Errorf("dummy3 type should be ordered")
	}

	if dummy3.Numeric() {
		t.Errorf("dummy3 type should not be numeric")
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
		return !strings.HasPrefix(f.Name(), "dummy")
	}

	typs2, err2 := getTypes("+test", filter)

	if err2 != nil {
		t.Error(err2)
	}

	if len(typs2) != 1 {
		t.Errorf("should have found the 1 marked-up type when filtered, found %v", len(typs2))
	}

	m2 := typeSliceToMap(typs2)

	if _, found := m2["dummy"]; found {
		t.Errorf("should not have found the dummy type")
	}

	if _, found := m2["app"]; !found {
		t.Errorf("should have found the app type")
	}

	// no false positives
	typs3, err3 := getTypes("+notreal", nil)

	if len(typs3) != 0 {
		t.Errorf("should have no marked-up types for +notreal")
	}

	if err3 != nil {
		t.Error(err3)
	}

	// should fail if types can't be evaluated
	// package.go by itself can't compile since it depends on other types

	filter4 := func(f os.FileInfo) bool {
		return f.Name() == "package.go"
	}

	_, err4 := getTypes("+test", filter4)

	if err4 == nil {
		t.Error("should have been unable to evaluate types of incomplete package")
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
