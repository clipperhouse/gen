package typewriter

import (
	"testing"
)

func TestApplicableTo(t *testing.T) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			typ := Type{
				comparable: i == 0,
				numeric:    i == 1,
				ordered:    i == 2,
			}
			tmpl := Template{
				RequiresComparable: j == 0,
				RequiresNumeric:    j == 1,
				RequiresOrdered:    j == 2,
			}
			applicable := tmpl.ApplicableTo(typ)
			should := i == j
			if should != applicable {
				t.Errorf("ApplicableTo is incorrect when for Type %v on Template %v; should be %v", typ, tmpl, should)
			}
		}
	}
}
