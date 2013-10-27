package models

import (
	"testing"
)

type test struct {
	Exec     func() interface{}
	Expected interface{}
}

func getTests() map[string][]test {
	// the basic pattern for tests is zero/one/many for 'some' slice & sanity checks on 'none' slice
	tests := make(map[string][]test)

	tests["AggregateInt"] = []test{
		test{
			func() interface{} {
				return some.AggregateInt(sum_theaters)
			},
			6 + 9 + 5,
		},
		test{
			func() interface{} {
				return none.AggregateInt(sum_theaters)
			},
			0,
		},
	}

	tests["AggregateString"] = []test{
		test{
			func() interface{} {
				return some.AggregateString(concat_title)
			},
			"first" + "second" + "third",
		},
		test{
			func() interface{} {
				return none.AggregateString(concat_title)
			},
			"",
		},
	}

	tests["All"] = []test{
		test{
			func() interface{} {
				return some.All(is_dummy)
			},
			false,
		},
		test{
			func() interface{} {
				return some.All(is_first)
			},
			false,
		},
		test{
			func() interface{} {
				return some.All(is_first_or_second_or_third)
			},
			true,
		},
		test{
			func() interface{} {
				return none.All(is_false)
			},
			true,
		},
		test{
			func() interface{} {
				return none.All(is_true)
			},
			true,
		},
	}

	tests["Any"] = []test{
		test{
			func() interface{} {
				return some.Any(is_dummy)
			},
			false,
		},
		test{
			func() interface{} {
				return some.Any(is_first)
			},
			true,
		},
		test{
			func() interface{} {
				return some.Any(is_first_or_third)
			},
			true,
		},
		test{
			func() interface{} {
				return none.Any(is_false)
			},
			false,
		},
		test{
			func() interface{} {
				return none.Any(is_true)
			},
			false,
		},
	}

	tests["Count"] = []test{
		test{
			func() interface{} {
				return some.Count(is_dummy)
			},
			0,
		},
		test{
			func() interface{} {
				return some.Count(is_first)
			},
			1,
		},
		test{
			func() interface{} {
				return some.Count(is_first_or_third)
			},
			2,
		},
		test{
			func() interface{} {
				return some.Count(is_true)
			},
			len(some),
		},
		test{
			func() interface{} {
				return none.Count(is_false)
			},
			0,
		},
		test{
			func() interface{} {
				return none.Count(is_true)
			},
			0,
		},
	}

	tests["GroupBy"] = []test{
		test{
			func() interface{} {
				return many.GroupByString(get_studio)
			},
			map[string]Movies{
				"Miramax":     Movies{_first, _fifth},
				"Universal":   Movies{_third, _fourth},
				"Warner Bros": Movies{_second},
			},
		},
	}

	tests["Sort"] = []test{
		test{
			func() interface{} {
				return many.Sort(by_title)
			},
			Movies{_fifth, _first, _fourth, _second, _third},
		},
		test{
			func() interface{} {
				return many.Sort(by_theaters)
			},
			Movies{_third, _first, _second, _fifth, _fourth},
		},
	}

	tests["SumInt"] = []test{
		test{
			func() interface{} {
				return some.SumInt(get_theaters)
			},
			6 + 9 + 5,
		},
		test{
			func() interface{} {
				return none.SumInt(get_theaters)
			},
			0,
		},
	}

	tests["Where"] = []test{
		test{
			func() interface{} {
				return some.Where(is_dummy)
			},
			Movies{},
		},
		test{
			func() interface{} {
				return some.Where(is_first)
			},
			Movies{_first},
		},
		test{
			func() interface{} {
				return some.Where(is_first_or_third)
			},
			Movies{_first, _third},
		},
		test{
			func() interface{} {
				return some.Where(is_true)
			},
			some,
		},
		test{
			func() interface{} {
				return none.Where(is_false)
			},
			Movies{},
		},
		test{
			func() interface{} {
				return none.Where(is_true)
			},
			Movies{},
		},
	}

	return tests
}

func TestAll(t *testing.T) {
	for _, tests := range getTests() {
		for _, test := range tests {
			switch test.Expected.(type) {
			default:
				got := test.Exec()
				if got != test.Expected {
					t.Errorf("Expected %v, got %v", test.Expected, got)
				}
			case map[string]Movies:
				got := test.Exec().(map[string]Movies)
				exp := test.Expected.(map[string]Movies)
				if len(got) != len(exp) {
					t.Errorf("Expected %v groups, got %v", len(exp), len(got))
					break
				}
				for k, _ := range got {
					got2 := got[k]
					exp2 := exp[k]
					if len(got2) != len(exp2) {
						t.Errorf("Expected %v Movies in %s element, got %v", len(exp2), k, len(got2))
						break
					}
					for i := range got2 {
						if got2[i] != exp2[i] {
							t.Errorf("Expected %v, got %v", exp2[i], got2[i])
						}
					}
				}
			case Movies:
				got := test.Exec().(Movies)
				exp := test.Expected.(Movies)
				if len(got) != len(exp) {
					t.Errorf("Expected %v Movies, got %v", len(exp), len(got))
					break
				}
				for i := range got {
					if got[i] != exp[i] {
						t.Errorf("Expected %v, got %v", exp[i], got[i])
					}
				}
			}
		}
	}
}

var _first = &Movie{Title: "first", Theaters: 6, Studio: "Miramax"}
var _second = &Movie{Title: "second", Theaters: 9, Studio: "Warner Bros"}
var _third = &Movie{Title: "third", Theaters: 5, Studio: "Universal"}
var _fourth = &Movie{Title: "fourth", Theaters: 50, Studio: "Universal"}
var _fifth = &Movie{Title: "fifth", Theaters: 20, Studio: "Miramax"}

var some = Movies{
	_first,
	_second,
	_third,
}

var none = Movies{}

var many = Movies{
	_first,
	_second,
	_third,
	_fourth,
	_fifth,
}

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
var get_studio = func(movie *Movie) string {
	return movie.Studio
}
var by_title = func(movies Movies, a, b int) bool {
	return movies[a].Title < movies[b].Title
}
var by_theaters = func(movies Movies, a, b int) bool {
	return movies[a].Theaters < movies[b].Theaters
}
var concat_title = func(movie *Movie, acc string) string {
	return acc + movie.Title
}
