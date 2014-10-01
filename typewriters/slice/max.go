package slice

import "github.com/clipperhouse/gen/typewriter"

var max = &typewriter.Template{
	Text: `
	// Max returns the maximum value of {{.SliceName}}. In the case of multiple items being equally maximal, the first such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#Max
	func (rcv {{.SliceName}}) Max() (result {{.Type}}, err error) {
		l := len(rcv)
		if l == 0 {
			err = errors.New("cannot determine the Max of an empty slice")
			return
		}
		result = rcv[0]
		for _, v := range rcv {
			if v > result {
				result = v
			}
		}
		return
	}
	`,
	TypeConstraint: typewriter.Constraint{Ordered: true},
}

var maxT = &typewriter.Template{
	Text: `
// Max{{.TypeParameter.LongName}} selects the largest value of {{.TypeParameter}} in {{.SliceName}}. Returns error on {{.SliceName}} with no elements. See: http://clipperhouse.github.io/gen/#MaxCustom
func (rcv {{.SliceName}}) Max{{.TypeParameter.LongName}}(fn func({{.Type}}) {{.TypeParameter}}) (result {{.TypeParameter}}, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine Max of zero-length {{.SliceName}}")
		return
	}
	result = fn(rcv[0])
	if l > 1 {
		for _, v := range rcv[1:] {
			f := fn(v)
			if f > result {
				result = f
			}
		}
	}
	return
}
`,
	TypeParameterConstraints: []typewriter.Constraint{
		// exactly one type parameter is required, and it must be ordered
		{Ordered: true},
	},
}
