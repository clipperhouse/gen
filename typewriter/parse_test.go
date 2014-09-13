package typewriter

import (
	"go/ast"
	"os"
	"strings"
	"testing"
)

type findDirectiveTest struct {
	text  string
	found bool
}

// dummy type evaluator
type eval struct{}

func (eval) Eval(name string) (Type, error) {
	return Type{
		Name:    name,
		Package: NewPackage("dummy", "dummy"),
	}, nil
}

func TestFindDirective(t *testing.T) {
	tests := []findDirectiveTest{
		{`// +test`, true},
		{`// +test foo:"bar,Baz"`, true},
		{`// there's nothing here`, false},
		{`//+test foo:"bar,Baz"`, true},
		{`//+test foo:"bar,Baz"`, true},
		{`// +test * foo:"bar,Baz"`, true},
		{`// +test foo:"bar,Baz" qux:"thing"`, true},
		{`// +tested`, false},
	}

	for i, test := range tests {
		g := &ast.CommentGroup{
			List: []*ast.Comment{{Text: test.text}},
		}
		c := findDirective(g, "+test")
		found := c != nil
		if found != test.found {
			t.Errorf("[test %v] found should have been %v for:\n%s", i, test.found, test.text)
		}
	}
}

type parseTest struct {
	comment string
	pointer Pointer
	tags    Tags
	valid   bool
}

func TestParse(t *testing.T) {
	tests := []parseTest{
		{`// +test foo:"bar,Baz"`, false, Tags{
			{"foo", []TagValue{
				{"bar", nil},
				{"Baz", nil},
			}, false},
		}, true},
		{`// +test * foo:"bar,Baz"`, true, Tags{
			{"foo", []TagValue{
				{"bar", nil},
				{"Baz", nil},
			}, false},
		}, true},
		{`// +test foo:"bar,Baz" qux:"stuff"`, false, Tags{
			{"foo", []TagValue{
				{"bar", nil},
				{"Baz", nil},
			}, false},
			{"qux", []TagValue{
				{"stuff", nil},
			}, false},
		}, true},
		{`// +test foo:"-bar,Baz"`, false, Tags{
			{"foo", []TagValue{
				{"bar", nil},
				{"Baz", nil},
			}, true},
		}, true},
		{`// +test foo:"bar  ,Baz "  `, false, Tags{
			{"foo", []TagValue{
				{"bar", nil},
				{"Baz", nil},
			}, false},
		}, true},
		{`// +test foo:"bar,Baz[qaz], qux"`, false, Tags{
			{"foo", []TagValue{
				{"bar", nil},
				{"Baz", []Type{
					Type{Name: "qaz"},
				}},
				{"qux", nil},
			}, false},
		}, true},
		{`// +test foo:"bar,Baz[qaz,hey]" qux:"stuff"`, false, Tags{
			{"foo", []TagValue{
				{"bar", nil},
				{"Baz", []Type{
					Type{Name: "qaz"},
					Type{Name: "hey"},
				}},
			}, false},
			{"qux", []TagValue{
				{"stuff", nil},
			}, false},
		}, true},
		{`// +test foo:"Baz[qaz],yo[dude]" qux:"stuff[things]"`, false, Tags{
			{"foo", []TagValue{
				{"Baz", []Type{
					Type{Name: "qaz"},
				}},
			}, false},
			{"qux", []TagValue{
				{"Baz", []Type{
					Type{Name: "qaz"},
				}},
			}, false},
		}, true},
		{`// +test foo:"bar,-Baz"`, false, nil, false},
		{`// +test foo:"bar,Baz-"`, false, nil, false},
		{`// +test foo:bar,Baz" qux:"stuff"`, false, nil, false},
		{`// +test foo"bar,Baz" qux:"stuff"`, false, nil, false},
		{`// +test foo:"bar,Baz" junk qux:"stuff"`, false, nil, false},
		{`// +test foo:"bar,Baz" 8qux:"stuff"`, false, nil, false},
		{`// +test fo^o:"bar,Baz" qux:"stuff"`, false, nil, false},
		{`// +test foo:"bar,Ba|z" qux:"stuff"`, false, nil, false},
		{`// +test foo:"bar,Baz" qux:"stuff`, false, nil, false},
		{`// +test *foo:"bar,Baz" qux:"stuff"`, false, nil, false},
		{`// +test foo:"bar,Baz" * qux:"stuff"`, false, nil, false},
		{`// +test * foo:"bar,Baz" * qux:"stuff"`, false, nil, false},
	}

	for i, test := range tests {
		c := &ast.Comment{
			Text: test.comment,
		}
		pointer, tags, err := parse(c.Text, "+test", eval{})

		if test.valid != (err == nil) {
			t.Errorf("[test %v] valid should have been %v for: %s\n%s", i, test.valid, test.comment, err)
		}

		if pointer != test.pointer {
			t.Errorf("[test %v] pointer should have been %v for: \n%s", i, bool(test.pointer), test.comment)
		}

		if !tagsEqual(tags, test.tags) {
			t.Errorf("[test %v] tags should have been %v, got %v", i, test.tags, tags)
		}
	}
}

func tagsEqual(tags, other Tags) bool {
	if len(tags) != len(other) {
		return false
	}

	for i, _ := range tags {
		t := tags[i]
		o := other[i]

		if t.Name != o.Name {
			return false
		}

		if t.Negated != o.Negated {
			return false
		}

		if len(t.Values) != len(o.Values) {
			return false
		}

		for j, _ := range t.Values {
			tv := t.Values[j]
			ov := o.Values[j]

			if tv.Name != ov.Name {
				return false
			}

			if len(tv.TypeParameters) != len(ov.TypeParameters) {
				return false
			}

			for k, _ := range tv.TypeParameters {
				if tv.TypeParameters[k].String() != ov.TypeParameters[k].String() {
					return false
				}
			}
		}
	}

	return true
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

	if t.Failed() {
		return // get out instead of nil panics below
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

	if len(dummy.Tags) != 1 {
		t.Errorf("typ should have 1 tag, found %v", len(m["dummy"].Tags))
	}

	if len(dummy.Tags[0].Values) != 1 {
		t.Errorf("Tag should have 1 Item, found %v", len(m["dummy"].Tags[0].Values))
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

	app := m["app"]

	if len(app.Tags) != 2 {
		t.Errorf("typ should have 2 Tags, found %v", len(app.Tags))
	}

	if len(app.Tags[0].Values) != 1 {
		t.Errorf("Tag should have 1 Item, found %v", len(app.Tags[0].Values))
	}

	if len(app.Tags[1].Values) != 2 {
		t.Errorf("Tag should have 2 Values, found %v", len(app.Tags[1].Values))
	}

	if len(app.Tags[1].Values[0].TypeParameters) != 1 {
		t.Errorf("TagValue should have 1 TypeParameter, found %v", len(app.Tags[1].Values[0].TypeParameters))
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
