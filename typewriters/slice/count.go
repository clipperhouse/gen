package slice

import "github.com/clipperhouse/gen/typewriter"

var count = &typewriter.Template{
	Text: `
// Count gives the number elements of {{.SliceName}} that return true for the passed func. See: http://clipperhouse.github.io/gen/#Count
func (rcv {{.SliceName}}) Count(fn func({{.Type}}) bool) (result int) {
	for _, v := range rcv {
		if fn(v) {
			result++
		}
	}
	return
}
`}
