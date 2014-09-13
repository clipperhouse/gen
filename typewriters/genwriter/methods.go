package genwriter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/clipperhouse/gen/typewriter"
)

// This business exists because I overload the methods tag to specify both standard and projection methods.
// Kind of a mess, but for the end user, arguably simpler. And arguably not.
func evaluateTags(t typewriter.Type) (standardMethods, projectionMethods []string, err error) {
	var nilMethods, nilProjections bool

	methods, found, methodsErr := t.Tags.ByName("methods")

	if methodsErr != nil {
		err = methodsErr
		return
	}

	nilMethods = !found // non-existent methods tag is different than empty

	_, found, projectionsErr := t.Tags.ByName("projections")

	if projectionsErr != nil {
		err = projectionsErr
		return
	}

	nilProjections = !found

	if nilMethods || methods.Negated {
		// default to all
		standardMethods = standardTemplates.GetAllKeys()
		if !nilProjections {
			projectionMethods = projectionTemplates.GetAllKeys()
		}
	}

	if !nilMethods {
		// categorize subsetted methods as standard or projection
		std := make([]string, 0)
		prj := make([]string, 0)

		// collect unknowns for err later
		unknown := make([]string, 0)

		for _, m := range methods.Values {
			isStd := standardTemplates.Contains(m.Name)
			if isStd {
				std = append(std, m.Name)
			}

			// only consider projection methods in presence of projected types
			isPrj := !nilProjections && projectionTemplates.Contains(m.Name)
			if isPrj {
				prj = append(prj, m.Name)
			}

			if !isStd && !isPrj {
				unknown = append(unknown, m.Name)
			}
		}

		if methods.Negated {
			standardMethods = remove(standardMethods, std...)
			projectionMethods = remove(projectionMethods, prj...)
		} else {
			standardMethods = std
			projectionMethods = prj
		}

		if len(unknown) > 0 {
			err = fmt.Errorf("method(s) %v on type %s are unknown", unknown, t.String())
			return
		}
	}

	return
}

func includeSortSupport(standardMethods []string) bool {
	for _, m := range standardMethods {
		if strings.HasPrefix(m, "SortBy") {
			return true
		}
	}
	return false
}

func includeSortInterface(standardMethods []string) bool {
	reg := regexp.MustCompile(`^Sort(Desc)?$`)
	for _, m := range standardMethods {
		if reg.MatchString(m) {
			return true
		}
	}
	return false
}
