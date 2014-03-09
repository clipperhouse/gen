package models

// +gen projections:"int,Thing2,string"
type Thing1 struct {
	Name   string
	Number int
}

// +gen methods:"Max,Min,Sort,IsSorted,SortDesc,IsSortedDesc" containers:"List,Ring"
type Thing2 Thing3

type Thing3 float64
