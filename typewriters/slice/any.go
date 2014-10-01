package slice

import "github.com/clipperhouse/gen/typewriter"

var any = &typewriter.Template{
	Text: `
// Any verifies that one or more elements of {{.SliceName}} return true for the passed func. See: http://clipperhouse.github.io/gen/#Any
func (rcv {{.SliceName}}) Any(fn func({{.Type}}) bool) bool {
	for _, v := range rcv {
		if fn(v) {
			return true
		}
	}
	return false
}
`}
