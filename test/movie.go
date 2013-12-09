package models

type Movie struct {
	Title             string `gen:"Select"`
	Theaters          int    `gen:"Aggregate,Sum,Max"`
	Studio            string `gen:"GroupBy"`
	BoxOfficeMillions int    `gen:"GroupBy,Min,Average"`
}

type Sub struct {
	// gen:"Where,Sort"
	Name string
}
