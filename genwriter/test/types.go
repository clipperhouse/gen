package main

// +test projections:"Other"
type Thing struct {
	Name   string
	Number Other
}

// methods where underlying type is ordered
// +test methods:"Max,Min,Sort,IsSorted,SortDesc,IsSortedDesc"
type Other Underlying

type Underlying int
