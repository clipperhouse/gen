package models

import (
	"testing"
)

type test struct {
	Exec          func() (interface{}, error)
	Expected      interface{}
	ErrorExpected bool
}

func getTests() map[string][]test {
	// the basic pattern for tests is zero/many for 'many' slice & sanity checks on 'none' slice
	tests := make(map[string][]test)

	tests["AggregateInt"] = []test{
		test{
			func() (interface{}, error) {
				return many.AggregateInt(sum_theaters), nil
			},
			6 + 9 + 5 + 50 + 20,
			false,
		},
		test{
			func() (interface{}, error) {
				return none.AggregateInt(sum_theaters), nil
			},
			0,
			false,
		},
	}

	tests["AggregateString"] = []test{
		test{
			func() (interface{}, error) {
				return many.AggregateString(concat_title), nil
			},
			"first" + "second" + "third" + "fourth" + "fifth",
			false,
		},
		test{
			func() (interface{}, error) {
				return none.AggregateString(concat_title), nil
			},
			"",
			false,
		},
	}

	tests["All"] = []test{
		test{
			func() (interface{}, error) {
				return many.All(is_dummy), nil
			},
			false,
			false,
		},
		test{
			func() (interface{}, error) {
				return many.All(is_first), nil
			},
			false,
			false,
		},
		test{
			func() (interface{}, error) {
				return many.All(is_first_or_second_or_third_or_fourth_or_fifth), nil
			},
			true,
			false,
		},
		test{
			func() (interface{}, error) {
				return none.All(is_false), nil
			},
			true,
			false,
		},
		test{
			func() (interface{}, error) {
				return none.All(is_true), nil
			},
			true,
			false,
		},
	}

	tests["Any"] = []test{
		test{
			func() (interface{}, error) {
				return many.Any(is_dummy), nil
			},
			false,
			false,
		},
		test{
			func() (interface{}, error) {
				return many.Any(is_first), nil
			},
			true,
			false,
		},
		test{
			func() (interface{}, error) {
				return many.Any(is_first_or_third), nil
			},
			true,
			false,
		},
		test{
			func() (interface{}, error) {
				return none.Any(is_false), nil
			},
			false,
			false,
		},
		test{
			func() (interface{}, error) {
				return none.Any(is_true), nil
			},
			false,
			false,
		},
	}

	tests["Count"] = []test{
		test{
			func() (interface{}, error) {
				return many.Count(is_dummy), nil
			},
			0,
			false,
		},
		test{
			func() (interface{}, error) {
				return many.Count(is_first), nil
			},
			1,
			false,
		},
		test{
			func() (interface{}, error) {
				return many.Count(is_first_or_third), nil
			},
			2,
			false,
		},
		test{
			func() (interface{}, error) {
				return many.Count(is_true), nil
			},
			len(many),
			false,
		},
		test{
			func() (interface{}, error) {
				return none.Count(is_false), nil
			},
			0,
			false,
		},
		test{
			func() (interface{}, error) {
				return none.Count(is_true), nil
			},
			0,
			false,
		},
	}

	tests["GroupBy"] = []test{
		test{
			func() (interface{}, error) {
				return many.GroupByString(get_studio), nil
			},
			map[string]Movies{
				"Miramax":     Movies{_first, _fifth},
				"Universal":   Movies{_third, _fourth},
				"Warner Bros": Movies{_second},
			},
			false,
		},
		test{
			func() (interface{}, error) {
				return many.GroupByInt(get_box_office), nil
			},
			map[int]Movies{
				90:  Movies{_first, _fourth},
				100: Movies{_second, _third, _fifth},
			},
			false,
		},
	}

	tests["Min"] = []test{
		test{
			func() (interface{}, error) {
				return many.Min(by_theaters)
			},
			_third,
			false,
		},
		test{
			func() (interface{}, error) {
				return many.Min(by_title)
			},
			_fifth,
			false,
		},
		test{
			func() (interface{}, error) {
				return none.Min(by_theaters)
			},
			_nil,
			true,
		},
	}

	tests["Max"] = []test{
		test{
			func() (interface{}, error) {
				return many.Max(by_theaters)
			},
			_fourth,
			false,
		},
		test{
			func() (interface{}, error) {
				return many.Max(by_title)
			},
			_third,
			false,
		},
		test{
			func() (interface{}, error) {
				return none.Max(by_theaters)
			},
			_nil,
			true,
		},
	}

	tests["Sort"] = []test{
		test{
			func() (interface{}, error) {
				return many.Sort(by_title), nil
			},
			Movies{_fifth, _first, _fourth, _second, _third},
			false,
		},
		test{
			func() (interface{}, error) {
				return many.Sort(by_theaters), nil
			},
			Movies{_third, _first, _second, _fifth, _fourth},
			false,
		},
		test{
			func() (interface{}, error) {
				return many.IsSorted(by_title), nil
			},
			false,
			false,
		},
		test{
			func() (interface{}, error) {
				sorted := many.Sort(by_title)
				return sorted.IsSorted(by_title), nil
			},
			true,
			false,
		},
		test{
			func() (interface{}, error) {
				return many.SortDesc(by_title), nil
			},
			Movies{_third, _second, _fourth, _first, _fifth},
			false,
		},
		test{
			func() (interface{}, error) {
				return many.SortDesc(by_theaters), nil
			},
			Movies{_fourth, _fifth, _second, _first, _third},
			false,
		},
		test{
			func() (interface{}, error) {
				return many.IsSortedDesc(by_title), nil
			},
			false,
			false,
		},
		test{
			func() (interface{}, error) {
				sorted := many.SortDesc(by_title)
				return sorted.IsSortedDesc(by_title), nil
			},
			true,
			false,
		},
	}

	tests["SumInt"] = []test{
		test{
			func() (interface{}, error) {
				return many.SumInt(get_theaters), nil
			},
			6 + 9 + 5 + 50 + 20,
			false,
		},
		test{
			func() (interface{}, error) {
				return none.SumInt(get_theaters), nil
			},
			0,
			false,
		},
	}

	tests["Where"] = []test{
		test{
			func() (interface{}, error) {
				return many.Where(is_dummy), nil
			},
			Movies{},
			false,
		},
		test{
			func() (interface{}, error) {
				return many.Where(is_first), nil
			},
			Movies{_first},
			false,
		},
		test{
			func() (interface{}, error) {
				return many.Where(is_first_or_third), nil
			},
			Movies{_first, _third},
			false,
		},
		test{
			func() (interface{}, error) {
				return many.Where(is_true), nil
			},
			many,
			false,
		},
		test{
			func() (interface{}, error) {
				return none.Where(is_false), nil
			},
			Movies{},
			false,
		},
		test{
			func() (interface{}, error) {
				return none.Where(is_true), nil
			},
			Movies{},
			false,
		},
	}

	return tests
}

func TestAll(t *testing.T) {
	for _, tests := range getTests() {
		for _, test := range tests {
			switch test.Expected.(type) {
			default:
				got, err := test.Exec()
				if test.ErrorExpected && err == nil {
					t.Errorf("Expected error but did not receive one")
				}
				if !test.ErrorExpected && err != nil {
					t.Errorf("Did not expect error but received: %v", err)
				}
				if got != test.Expected {
					t.Errorf("Expected %v, got %v", test.Expected, got)
				}
			case map[int]Movies:
				_got, err := test.Exec()
				if test.ErrorExpected && err == nil {
					t.Errorf("Expected error but did not receive one")
				}
				if !test.ErrorExpected && err != nil {
					t.Errorf("Did not expect error but received: %v", err)
				}
				got := _got.(map[int]Movies)
				exp := test.Expected.(map[int]Movies)
				if len(got) != len(exp) {
					t.Errorf("Expected %v groups, got %v", len(exp), len(got))
					break
				}
				for k, _ := range got {
					got2 := got[k]
					exp2 := exp[k]
					if len(got2) != len(exp2) {
						t.Errorf("Expected %v Movies in %d element, got %v", len(exp2), k, len(got2))
						break
					}
					for i := range got2 {
						if got2[i] != exp2[i] {
							t.Errorf("Expected %v, got %v", exp2[i], got2[i])
						}
					}
				}
			case map[string]Movies:
				_got, err := test.Exec()
				if test.ErrorExpected && err == nil {
					t.Errorf("Expected error but did not receive one")
				}
				if !test.ErrorExpected && err != nil {
					t.Errorf("Did not expect error but received: %v", err)
				}
				got := _got.(map[string]Movies)
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
				_got, err := test.Exec()
				if test.ErrorExpected && err == nil {
					t.Errorf("Expected error but did not receive one")
				}
				if !test.ErrorExpected && err != nil {
					t.Errorf("Did not expect error but received: %v", err)
				}
				got := _got.(Movies)
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

var _nil *Movie

var _first = &Movie{Title: "first", Theaters: 6, Studio: "Miramax", BoxOfficeMillions: 90}
var _second = &Movie{Title: "second", Theaters: 9, Studio: "Warner Bros", BoxOfficeMillions: 100}
var _third = &Movie{Title: "third", Theaters: 5, Studio: "Universal", BoxOfficeMillions: 100}
var _fourth = &Movie{Title: "fourth", Theaters: 50, Studio: "Universal", BoxOfficeMillions: 90}
var _fifth = &Movie{Title: "fifth", Theaters: 20, Studio: "Miramax", BoxOfficeMillions: 100}

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
var is_first_or_second_or_third_or_fourth_or_fifth = func(movie *Movie) bool {
	return movie.Title == "first" || movie.Title == "second" || movie.Title == "third" || movie.Title == "fourth" || movie.Title == "fifth"
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
var get_box_office = func(movie *Movie) int {
	return movie.BoxOfficeMillions
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
