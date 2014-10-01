package slice

import "github.com/clipperhouse/gen/typewriter"

var distinct = &typewriter.Template{
	Text: `
// Distinct returns a new {{.SliceName}} whose elements are unique. See: http://clipperhouse.github.io/gen/#Distinct
func (rcv {{.SliceName}}) Distinct() (result {{.SliceName}}) {
	appended := make(map[{{.Type}}]bool)
	for _, v := range rcv {
		if !appended[v] {
			result = append(result, v)
			appended[v] = true
		}
	}
	return result
}
`,
	TypeConstraint: typewriter.Constraint{Comparable: true},
}
