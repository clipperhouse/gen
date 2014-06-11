package typewriter

import (
	// "fmt"
	"testing"
)

func TestParseTags(t *testing.T) {
	doc := `// some stuff that's actually a comment
+test foo:"bar,Baz"
`
	pointer, tags, found := parseTags("+test", doc)

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
	pointer2, tags2, found2 := parseTags("+test", doc2)

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
	pointer3, tags3, found3 := parseTags("+test", doc3)

	if !found3 {
		t.Errorf("tags should have been found in %s", doc3)
	}

	if !pointer3 {
		t.Errorf("pointer should have been found in %s", doc3)
	}

	if len(tags3) != 3 {
		t.Errorf("3 tags should have been found in %s; found %v", doc3, len(tags3))
	}
}

func TestParseTag(t *testing.T) {
	s := `foo:"bar,Baz"`
	tag, found := parseTag(s)

	if !found {
		t.Errorf("tag foo should have been found in %s", s)
	}

	should := []string{"bar", "Baz"}

	if !sliceEqual(tag.Items, should) {
		t.Errorf("Items should be %v, got %v", should, tag.Items)
	}

	s2 := ""
	tag2, found2 := parseTag(s2)

	if found2 {
		t.Errorf("empty tags should not have been found in %s", tag2)
	}

	s3 := `foo:"-Baz,qaz"`
	tag3, found3 := parseTag(s3)

	if !found3 {
		t.Errorf("tag foo should have been found in %s", s3)
	}

	if !tag3.Negated {
		t.Errorf("tag %s should be negated", tag2)
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
