package typewriter

import "fmt"

// +gen methods:"Where"
type Tag struct {
	Name     string
	Items    []string
	Negated  bool
	original string
}

func (ts Tags) ByName(name string) (result Tag, found bool, err error) {
	tags := ts.Where(func(t Tag) bool {
		return t.Name == name
	})

	if len(tags) == 0 {
		return
	}

	if len(tags) == 1 {
		found = true
		result = tags[0]
		return
	}

	err = fmt.Errorf("more than one '%s' specified", name)

	return
}

func (t Tag) String() string {
	return t.original
}
