package main

// +test slice:"Any,All,Count,Distinct,DistinctBy,Each,First,MaxBy,MinBy,Single,Where,SortBy,SortByDesc,IsSortedBy,IsSortedByDesc,Aggregate[Other],Average[Other],GroupBy[Other],Max[Other],Min[Other],Select[Other],Sum[Other]"
type Thing struct {
	Name   string
	Number Other
}
