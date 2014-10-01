package slice

import "github.com/clipperhouse/gen/typewriter"

var sumT = &typewriter.Template{
	Text: `
// Sum{{.TypeParameter.LongName}} sums {{.Type}} over elements in {{.SliceName}}. See: http://clipperhouse.github.io/gen/#Sum
func (rcv {{.SliceName}}) Sum{{.TypeParameter.LongName}}(fn func({{.Type}}) {{.TypeParameter}}) (result {{.TypeParameter}}) {
	for _, v := range rcv {
		result += fn(v)
	}
	return
}
`,
	TypeParameterConstraints: []typewriter.Constraint{
		// exactly one type parameter is required, and it must be numeric
		{Numeric: true},
	},
}
