package slice

import (
	"github.com/clipperhouse/gen/typewriter"
)

// a convenience for passing values into templates; in MVC it'd be called a view model
type model struct {
	Type      typewriter.Type
	SliceName string
	// these templates only ever happen to use one type parameter
	TypeParameter typewriter.Type
	typewriter.TagValue
}

var templates = typewriter.TemplateSet{

	"slice": slice,

	"Aggregate[T]": aggregateT,
	"All":          all,
	"Any":          any,
	"Average":      average,
	"Average[T]":   averageT,
	"Count":        count,
	"Distinct":     distinct,
	"DistinctBy":   distinctBy,
	"Each":         each,
	"First":        first,
	"GroupBy[T]":   groupByT,
	"Max":          max,
	"Max[T]":       maxT,
	"MaxBy":        maxBy,
	"Min":          min,
	"Min[T]":       minT,
	"MinBy":        minBy,
	"Select[T]":    selectT,
	"Single":       single,
	"Sum[T]":       sumT,
	"Where":        where,

	"Sort":         sort,
	"IsSorted":     isSorted,
	"SortDesc":     sortDesc,
	"IsSortedDesc": isSortedDesc,

	"SortBy":         sortBy,
	"IsSortedBy":     isSortedBy,
	"SortByDesc":     sortByDesc,
	"IsSortedByDesc": isSortedByDesc,

	"sortImplementation": sortImplementation,
	"sortInterface":      sortInterface,
}
