package slice

import "github.com/clipperhouse/gen/typewriter"

var single = &typewriter.Template{
	Text: `
// Single returns exactly one element of {{.SliceName}} that returns true for the passed func. Returns error if no or multiple elements return true. See: http://clipperhouse.github.io/gen/#Single
func (rcv {{.SliceName}}) Single(fn func({{.Type}}) bool) (result {{.Type}}, err error) {
	var candidate {{.Type}}
	found := false
	for _, v := range rcv {
		if fn(v) {
			if found {
				err = errors.New("multiple {{.SliceName}} elements return true for passed func")
				return
			}
			candidate = v
			found = true
		}
	}
	if found {
		result = candidate
	} else {
		err = errors.New("no {{.SliceName}} elements return true for passed func")
	}
	return
}
`}
