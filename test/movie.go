package models

type Movie struct {
	Title             string `gen:"Select"`
	Theaters          int    `gen:"Aggregate,Sum"`
	Studio            string `gen:"GroupBy"`
	BoxOfficeMillions int    `gen:"GroupBy"`
}
