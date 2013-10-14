package models

import (
	"log"
	"testing"
)

var some = Movies{
	&Movie{Title: "first"},
	&Movie{Title: "second"},
	&Movie{Title: "third"},
}

var none = Movies{}

var is_first = func(movie *Movie) bool {
	return movie.Title == "first"
}
var is_first_or_second_or_third = func(movie *Movie) bool {
	return movie.Title == "first" || movie.Title == "second" || movie.Title == "third"
}
var is_first_or_third = func(movie *Movie) bool {
	return movie.Title == "first" || movie.Title == "third"
}
var is_dummy = func(movie *Movie) bool {
	return movie.Title == "dummy"
}
var is_true = func(movie *Movie) bool {
	return true
}

func TestAll(t *testing.T) {
	if false != some.All(is_first) {
		log.Println(some.All(is_first))
		t.Fail()
	}

	if false != some.All(is_first_or_third) {
		log.Println(some.All(is_first))
		t.Fail()
	}

	if true != some.All(is_first_or_second_or_third) {
		log.Println(some.All(is_first_or_second_or_third))
		t.Fail()
	}

	// TODO: what's the right behavior on empty set?
}

func TestAny(t *testing.T) {
	if true != some.Any(is_first) {
		log.Println(some.Any(is_first))
		t.Fail()
	}

	if true != some.Any(is_first_or_third) {
		log.Println(some.Any(is_first))
		t.Fail()
	}

	if false != some.Any(is_dummy) {
		log.Println(some.Any(is_dummy))
		t.Fail()
	}

	if false != none.Any(is_first) {
		log.Println(some.Any(is_first))
		t.Fail()
	}

	if false != none.Any(is_true) {
		log.Println(some.Any(is_true))
		t.Fail()
	}
}

func TestCount(t *testing.T) {
	if some.Count(is_first) != 1 {
		log.Println(some.Count(is_first))
		t.Fail()
	}

	if some.Count(is_first_or_third) != 2 {
		log.Println(some.Count(is_first_or_third))
		t.Fail()
	}

	if some.Count(is_dummy) != 0 {
		log.Println(some.Count(is_dummy))
		t.Fail()
	}

	if none.Count(is_first) != 0 {
		log.Println(none.Count(is_first))
		t.Fail()
	}

	if none.Count(is_true) != 0 {
		log.Println(none.Count(is_true))
		t.Fail()
	}
}

func TestWhere(t *testing.T) {
	var where_first = some.Where(is_first)
	if len(where_first) != 1 || where_first[0].Title != "first" {
		log.Println(len(where_first))
		t.Fail()
	}

	var where_first_or_third = some.Where(is_first_or_third)
	if len(where_first_or_third) != 2 || where_first_or_third[1].Title != "third" {
		log.Println(len(where_first_or_third))
		t.Fail()
	}

	var where_dummy = some.Where(is_dummy)
	if len(where_dummy) != 0 {
		log.Println(len(where_dummy))
		t.Fail()
	}

	var where_none = none.Where(is_true)
	if len(where_none) != 0 {
		log.Println(len(where_none))
		t.Fail()
	}
}

func TestNil(t *testing.T) {
	if some.Count(nil) != len(some) {
		t.Error("count")
	}

	var where = some.Where(nil)
	if len(where) != len(some) {
		t.Error("where")
	}

	if !some.All(nil) {
		t.Error("all")
	}

	if !some.Any(nil) {
		t.Errorf("any")
	}
}
