package models

import (
	"log"
	"testing"
)

var some = Movies{
	&Movie{Title: "first", Theaters: 6},
	&Movie{Title: "second", Theaters: 9},
	&Movie{Title: "third", Theaters: 5},
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
var is_false = func(movie *Movie) bool {
	return false
}
var get_theaters = func(movie *Movie) int {
	return movie.Theaters
}
var sum_theaters = func(movie *Movie, acc int) int {
	return acc + movie.Theaters
}
var get_title = func(movie *Movie) string {
	return movie.Title
}
var concat_title = func(movie *Movie, acc string) string {
	return acc + movie.Title
}

func TestAggregateInt(t *testing.T) {
	expected := 6 + 9 + 5

	if expected != some.AggregateInt(sum_theaters) {
		t.Error(some.AggregateInt(sum_theaters))
	}
}

func TestAggregateString(t *testing.T) {
	expected := "first" + "second" + "third"

	if expected != some.AggregateString(concat_title) {
		t.Error(some.AggregateString(concat_title))
	}
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

	// All is always true on empty collection
	if true != none.All(is_true) {
		log.Println(some.All(is_true))
		t.Fail()
	}

	if true != none.All(is_false) {
		log.Println(some.All(is_false))
		t.Fail()
	}
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

func TestJoinString(t *testing.T) {
	expected := "first, second, third"

	if expected != some.JoinString(get_title, ", ") {
		t.Error(some.JoinString(get_title, ", "))
	}
}

func TestSumInt(t *testing.T) {
	expected := 6 + 9 + 5

	if expected != some.SumInt(get_theaters) {
		t.Error(some.SumInt(get_theaters))
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
