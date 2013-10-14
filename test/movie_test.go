package models

import (
	"log"
	"testing"
)

func TestMovies(t *testing.T) {
	var some = Movies{
		&Movie{Title: "first"},
		&Movie{Title: "second"},
		&Movie{Title: "third"},
	}

	var none = Movies{}

	var is_first = func(movie *Movie) bool {
		return movie.Title == "first"
	}

	if true != some.Any(is_first) {
		log.Println(some.Any(is_first))
		t.Fail()
	}

	var is_dummy = func(movie *Movie) bool {
		return movie.Title == "dummy"
	}

	if false != some.Any(is_dummy) {
		log.Println(some.Any(is_dummy))
		t.Fail()
	}

	var where_first = some.Where(is_first)
	if len(where_first) != 1 {
		log.Println(len(where_first))
		t.Fail()
	}

	var is_first_or_third = func(movie *Movie) bool {
		return movie.Title == "first" || movie.Title == "third"
	}

	var where_first_or_third = some.Where(is_first_or_third)
	if len(where_first_or_third) != 2 {
		log.Println(len(where_first_or_third))
		t.Fail()
	}

	var count_first_or_third = some.Count(is_first_or_third)
	if count_first_or_third != 2 {
		log.Println(count_first_or_third)
		t.Fail()
	}

	var count_none = none.Count(is_first_or_third)
	if count_none != 0 {
		log.Println(count_none)
		t.Fail()
	}

	var is_true = func(movie *Movie) bool {
		return true
	}
	if false != none.Any(is_true) {
		log.Println(none.Any(is_true))
		t.Fail()
	}
}
