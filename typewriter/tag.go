package typewriter

import "fmt"

// +gen methods:"Where"
type Tag struct {
	Name    string
	Items   []string
	Negated bool
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

func (tags Tags) Equal(other Tags) bool {
	if len(tags) != len(other) {
		return false
	}

	for i, _ := range tags {
		if tags[i].Name != other[i].Name {
			return false
		}

		if tags[i].Negated != other[i].Negated {
			return false
		}

		if len(tags[i].Items) != len(other[i].Items) {
			return false
		}

		for j, _ := range tags[i].Items {
			if tags[i].Items[j] != other[i].Items[j] {
				return false
			}
		}
	}

	return true
}
