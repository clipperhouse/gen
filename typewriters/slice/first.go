package slice

import "github.com/clipperhouse/gen/typewriter"

var first = &typewriter.Template{
	Text: `
// First returns the first element that returns true for the passed func. Returns error if no elements return true. See: http://clipperhouse.github.io/gen/#First
func (rcv {{.SliceName}}) First(fn func({{.Type}}) bool) (result {{.Type}}, err error) {
	for _, v := range rcv {
		if fn(v) {
			result = v
			return
		}
	}
	err = errors.New("no {{.SliceName}} elements return true for passed func")
	return
}
`}
