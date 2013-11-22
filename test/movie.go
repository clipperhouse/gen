package models

type Movie struct {
	Title             string `gen:"Select"`
	Theaters          int    `gen:"SortBy"`
	Studio            string `gen:"DistinctBy,SortBy"`
	BoxOfficeMillions int
}
