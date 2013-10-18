package models

import (
	"log"
	"testing"
)

type genTest struct {
	in    Movies
	fn    func(*Movie) bool
	all   bool
	any   bool
	count int
	where Movies
}

var genTests = []genTest{
	genTest{
		some,
		is_first,
		false,
		true,
		1,
		Movies{_first},
	},
	genTest{
		some,
		is_first_or_second_or_third,
		true,
		true,
		3,
		Movies{_first, _second, _third},
	},
	genTest{
		some,
		is_first_or_third,
		false,
		true,
		2,
		Movies{_first, _third},
	},
	genTest{
		some,
		is_dummy,
		false,
		false,
		0,
		nil,
	},
	genTest{
		none,
		is_true,
		true,
		false,
		0,
		nil,
	},
	genTest{
		none,
		is_false,
		true,
		false,
		0,
		nil,
	},
	genTest{
		none,
		is_first,
		true,
		false,
		0,
		nil,
	},
}

func TestGen(t *testing.T) {
	for _, v := range genTests {
		if o := v.in.All(v.fn); o != v.all {
			t.Errorf("all error: expected %v, got %v", v.all, o)
		}
		if o := v.in.Any(v.fn); o != v.any {
			t.Errorf("any error: expected %v, got %v", v.any, o)
		}
		if o := v.in.Count(v.fn); o != v.count {
			t.Errorf("count error: expected %v, got %v", v.count, o)
		}
		if o := v.in.Where(v.fn); !same(o, v.where) {
			t.Errorf("where error: expected %v, got %v", v.where, o)
		}
	}
}

func same(a, b Movies) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if *v != *b[i] {
			return false
		}
	}
	return true
}

var _first = &Movie{Title: "first", Theaters: 6}
var _second = &Movie{Title: "second", Theaters: 9}
var _third = &Movie{Title: "third", Theaters: 5}

var some = Movies{
	_first,
	_second,
	_third,
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
