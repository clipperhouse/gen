package slice

import "github.com/clipperhouse/gen/typewriter"

var sort = &typewriter.Template{
	Text: `
// Sort returns a new ordered {{.SliceName}}. See: http://clipperhouse.github.io/gen/#Sort
func (rcv {{.SliceName}}) Sort() {{.SliceName}} {
	result := make({{.SliceName}}, len(rcv))
	copy(result, rcv)
	sort.Sort(result)
	return result
}
`,
	TypeConstraint: typewriter.Constraint{Ordered: true},
}

var isSorted = &typewriter.Template{
	Text: `
// IsSorted reports whether {{.SliceName}} is sorted. See: http://clipperhouse.github.io/gen/#Sort
func (rcv {{.SliceName}}) IsSorted() bool {
	return sort.IsSorted(rcv)
}
`,
	TypeConstraint: typewriter.Constraint{Ordered: true},
}

var sortDesc = &typewriter.Template{
	Text: `
// SortDesc returns a new reverse-ordered {{.SliceName}}. See: http://clipperhouse.github.io/gen/#Sort
func (rcv {{.SliceName}}) SortDesc() {{.SliceName}} {
	result := make({{.SliceName}}, len(rcv))
	copy(result, rcv)
	sort.Sort(sort.Reverse(result))
	return result
}
`,
	TypeConstraint: typewriter.Constraint{Ordered: true},
}

var isSortedDesc = &typewriter.Template{
	Text: `
// IsSortedDesc reports whether {{.SliceName}} is reverse-sorted. See: http://clipperhouse.github.io/gen/#Sort
func (rcv {{.SliceName}}) IsSortedDesc() bool {
	return sort.IsSorted(sort.Reverse(rcv))
}
`,
	TypeConstraint: typewriter.Constraint{Ordered: true},
}
