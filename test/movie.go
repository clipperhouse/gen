package models

type Movie struct {
	Title             string
	Theaters          int    `gen:"SortBy"`
	Studio            string `gen:"DistinctBy,SortBy"`
	BoxOfficeMillions int
}
