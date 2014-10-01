package slice

import "github.com/clipperhouse/gen/typewriter"

var each = &typewriter.Template{
	Text: `
// Each iterates over {{.SliceName}} and executes the passed func against each element. See: http://clipperhouse.github.io/gen/#Each
func (rcv {{.SliceName}}) Each(fn func({{.Type}})) {
	for _, v := range rcv {
		fn(v)
	}
}
`}
