package slice

import "github.com/clipperhouse/gen/typewriter"

var selectT = &typewriter.Template{
	Text: `
// Select{{.TypeParameter.LongName}} returns a slice of {{.Type}} in {{.SliceName}}, projected by passed func. See: http://clipperhouse.github.io/gen/#Select
func (rcv {{.SliceName}}) Select{{.TypeParameter.LongName}}(fn func({{.Type}}) {{.TypeParameter}}) (result []{{.TypeParameter}}) {
	for _, v := range rcv {
		result = append(result, fn(v))
	}
	return
}
`,
	TypeParameterConstraints: []typewriter.Constraint{
		// exactly one type parameter is required, but no constraints on that type
		{},
	},
}
