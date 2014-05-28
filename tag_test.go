package typewriter

import (
	// "fmt"
	"testing"
)

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
