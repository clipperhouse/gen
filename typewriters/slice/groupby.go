package slice

import "github.com/clipperhouse/gen/typewriter"

var groupByT = &typewriter.Template{
	Text: `
// GroupBy{{.TypeParameter.LongName}} groups elements into a map keyed by {{.TypeParameter}}. See: http://clipperhouse.github.io/gen/#GroupBy
func (rcv {{.SliceName}}) GroupBy{{.TypeParameter.LongName}}(fn func({{.Type}}) {{.TypeParameter}}) map[{{.TypeParameter}}]{{.SliceName}} {
	result := make(map[{{.TypeParameter}}]{{.SliceName}})
	for _, v := range rcv {
		key := fn(v)
		result[key] = append(result[key], v)
	}
	return result
}
`,
	TypeParameterConstraints: []typewriter.Constraint{
		// exactly one type parameter is required, and it must be comparable
		{Comparable: true},
	},
}
