package models

// Any amount of docs might be here
// project:"int"
type Movie struct {
	Title             string
	Theaters          int
	Studio            string
	BoxOfficeMillions int
}

type Thing float64

type Thing2 Thing
