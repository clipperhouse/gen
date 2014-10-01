package slice

import "github.com/clipperhouse/gen/typewriter"

var where = &typewriter.Template{
	Text: `
// Where returns a new {{.SliceName}} whose elements return true for func. See: http://clipperhouse.github.io/gen/#Where
func (rcv {{.SliceName}}) Where(fn func({{.Type}}) bool) (result {{.SliceName}}) {
	for _, v := range rcv {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
`}
