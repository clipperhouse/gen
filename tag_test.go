package typewriter

import (
	// "fmt"
	"testing"
)

func TestParseTag(t *testing.T) {
	s := `foo:"bar,Baz"`
	tag, found := parseTag(s)

	if !found {
		t.Errorf("tag foo should have been found in %s", tag)
	}

	should := []string{"bar", "Baz"}

	if !sliceEqual(tag.Items, should) {
		t.Errorf("Items should be %v, got %v", should, tag.Items)
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
