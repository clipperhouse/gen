package typewriter

import "testing"

func TestTryType(t *testing.T) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			typ := Type{
				comparable: i == 0,
				numeric:    i == 1,
				ordered:    i == 2,
			}
			c := Constraint{
				Comparable: j == 0,
				Numeric:    j == 1,
				Ordered:    j == 2,
			}

			err := c.tryType(typ)
			should := i == j

			if should != (err == nil) {
				t.Errorf("tryType is incorrect when for Type %v on Constraint %v; should be %v", typ, c, should)
			}
		}
	}
}
