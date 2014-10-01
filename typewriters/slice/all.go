package slice

import "github.com/clipperhouse/gen/typewriter"

var all = &typewriter.Template{
	Text: `
// All verifies that all elements of {{.SliceName}} return true for the passed func. See: http://clipperhouse.github.io/gen/#All
func (rcv {{.SliceName}}) All(fn func({{.Type}}) bool) bool {
	for _, v := range rcv {
		if !fn(v) {
			return false
		}
	}
	return true
}
`}
