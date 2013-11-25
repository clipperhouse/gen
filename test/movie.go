package models

type Movie struct {
	Title             string `gen:"Select"`
	Theaters          int    `gen:"SortBy,Aggregate,Sum"`
	Studio            string `gen:"DistinctBy,SortBy,GroupBy"`
	BoxOfficeMillions int    `gen:"GroupBy"`
}
