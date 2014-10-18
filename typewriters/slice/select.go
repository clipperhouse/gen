package slice

import "github.com/clipperhouse/gen/typewriter"

var selectT = &typewriter.Template{
	Text: `
// Select{{.TypeParameter.LongName}} projects a slice of {{.TypeParameter}} from {{.SliceName}}, typically called a map in other frameworks. See: http://clipperhouse.github.io/gen/#Select
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
