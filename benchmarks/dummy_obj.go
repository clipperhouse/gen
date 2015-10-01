package benchmarks

// +gen * slice:"Any,Select[*dummyDestinationSelectObject],SortBy"
type dummyObject struct {
	Name string
	Num  int
}

type dummyDestinationSelectObject struct {
	Name string
}
