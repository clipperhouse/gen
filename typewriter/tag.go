package typewriter

import (
	"fmt"
	"strings"
)

// +gen methods:"Where"
type Tag struct {
	Name    string
	Values  []TagValue
	Negated bool
}

type TagValue struct {
	Name           string
	TypeParameters []Type
}

func (v TagValue) String() string {
	if len(v.TypeParameters) == 0 {
		return v.Name
	}

	var a []string
	for i := 0; i < len(v.TypeParameters); i++ {
		a = append(a, v.TypeParameters[i].Name)
	}

	return v.Name + "[" + strings.Join(a, ",") + "]"
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
