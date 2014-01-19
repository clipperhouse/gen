package models

// +gen * projections:"int,Thing2,string"
type Movie struct {
	Title             string
	Theaters          int
	Studio            string
	BoxOfficeMillions int
}

type Thing float64

type Thing2 Thing
