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

	tests["AggregateTheaters"] = []test{
		test{
			func() (interface{}, error) {
				return many.AggregateTheaters(sumInt), nil
			},
			6 + 9 + 5 + 50 + 20,
			false,
		},
		test{
			func() (interface{}, error) {
				return none.AggregateTheaters(sumInt), nil
			},
			0,
			false,
		},
	}

	tests["SumTheaters"] = []test{
		test{
			func() (interface{}, error) {
				return many.SumTheaters(), nil
			},
			6 + 9 + 5 + 50 + 20,
			false,
		},
		test{
			func() (interface{}, error) {
				return none.SumTheaters(), nil
			},
			0,
			false,
		},
	}

	tests["MaxTheaters"] = []test{
		test{
			func() (interface{}, error) {
				return many.MaxTheaters()
			},
			50,
			false,
		},
		test{
			func() (interface{}, error) {
				return Movies{second}.MaxTheaters()
			},
			second.Theaters,
			false,
		},
		test{
			func() (interface{}, error) {
				return none.MaxTheaters()
			},
			0,
			true,
		},
	}

	tests["All"] = []test{
		test{
			func() (interface{}, error) {
				return many.All(isDummy), nil
			},
			false,
			false,
		},
		test{
			func() (interface{}, error) {
				return many.All(isFirst), nil
			},
			false,
			false,
		},
		test{
			func() (interface{}, error) {
				return many.All(isFirstOrSecondOrThirdOrFourthOrFifth), nil
			},
			true,
			false,
		},
		test{
			func() (interface{}, error) {
				return none.All(isFalse), nil
			},
			true,
			false,
		},
		test{
			func() (interface{}, error) {
				return none.All(isTrue), nil
			},
			true,
			false,
		},
	}

	tests["Any"] = []test{
		test{
			func() (interface{}, error) {
				return many.Any(isDummy), nil
			},
			false,
			false,
		},
		test{
			func() (interface{}, error) {
				return many.Any(isFirst), nil
			},
			true,
			false,
		},
		test{
			func() (interface{}, error) {
				return many.Any(isFirstOrThird), nil
			},
			true,
			false,
		},
		test{
			func() (interface{}, error) {
				return none.Any(isFalse), nil
			},
			false,
			false,
		},
		test{
			func() (interface{}, error) {
				return none.Any(isTrue), nil
			},
			false,
			false,
		},
	}

	tests["Count"] = []test{
		test{
			func() (interface{}, error) {
				return many.Count(isDummy), nil
			},
			0,
			false,
		},
		test{
			func() (interface{}, error) {
				return many.Count(isFirst), nil
			},
			1,
			false,
		},
		test{
			func() (interface{}, error) {
				return many.Count(isFirstOrThird), nil
			},
			2,
			false,
		},
		test{
			func() (interface{}, error) {
				return many.Count(isTrue), nil
			},
			len(many),
			false,
		},
		test{
			func() (interface{}, error) {
				return none.Count(isFalse), nil
			},
			0,
			false,
		},
		test{
			func() (interface{}, error) {
				return none.Count(isTrue), nil
			},
			0,
			false,
		},
	}

	tests["Distinct"] = []test{
		test{
			func() (interface{}, error) {
				return Movies{first, second, third, fourth, fourth, fifth, third}.Distinct(), nil // TODO: value test
			},
			Movies{first, second, third, fourth, fifth},
			false,
		},
	}

	tests["DistinctBy"] = []test{
		test{
			func() (interface{}, error) {
				return Movies{first, third, fourth, fourth, fifth, third}.DistinctBy(sameTitle), nil
			},
			Movies{first, third, fourth, fifth},
			false,
		},
		test{
			func() (interface{}, error) {
				return Movies{first, third, fourth, fourth, fifth, third}.DistinctBy(sameMillions), nil
			},
			Movies{first, third},
			false,
		},
	}

	tests["SelectTitle"] = []test{
		test{
			func() (interface{}, error) {
				return many.SelectTitle(), nil
			},
			[]string{"first", "second", "third", "fourth", "fifth"},
			false,
		},
		test{
			func() (interface{}, error) {
				return none.SelectTitle(), nil
			},
			[]string{},
			false,
		},
	}

	tests["First"] = []test{
		test{
			func() (interface{}, error) {
				return many.First(isThird)
			},
			third,
			false,
		},
		test{
			func() (interface{}, error) {
				return many.First(isDummy)
			},
			_nil,
			true,
		},
		test{
			func() (interface{}, error) {
				return none.First(isFalse)
			},
			_nil,
			true,
		},
		test{
			func() (interface{}, error) {
				return none.First(isTrue)
			},
			_nil,
			true,
		},
	}

	tests["GroupByStudio"] = []test{
		test{
			func() (interface{}, error) {
				return many.GroupByStudio(), nil
			},
			map[string]Movies{
				"Miramax":     Movies{first, fifth},
				"Universal":   Movies{third, fourth},
				"Warner Bros": Movies{second},
			},
			false,
		},
	}

	tests["GroupByBoxOfficeMillions"] = []test{
		test{
			func() (interface{}, error) {
				return many.GroupByBoxOfficeMillions(), nil
			},
			map[int]Movies{
				90:  Movies{first, fourth},
				100: Movies{second, third, fifth},
			},
			false,
		},
	}

	tests["Min"] = []test{
		test{
			func() (interface{}, error) {
				return many.Min(byTheaters)
			},
			third,
			false,
		},
		test{
			func() (interface{}, error) {
				return many.Min(byTitle)
			},
			fifth,
			false,
		},
		test{
			func() (interface{}, error) {
				return none.Min(byTheaters)
			},
			_nil,
			true,
		},
	}

	tests["Max"] = []test{
		test{
			func() (interface{}, error) {
				return many.Max(byTheaters)
			},
			fourth,
			false,
		},
		test{
			func() (interface{}, error) {
				return many.Max(byTitle)
			},
			third,
			false,
		},
		test{
			func() (interface{}, error) {
				return none.Max(byTheaters)
			},
			_nil,
			true,
		},
	}

	tests["Single"] = []test{
		test{
			func() (interface{}, error) {
				return many.Single(isThird)
			},
			third,
			false,
		},
		test{
			func() (interface{}, error) {
				return many.Single(isDummy)
			},
			_nil,
			true,
		},
		test{
			func() (interface{}, error) {
				return Movies{third, fourth, fifth, third, first}.Single(isThird)
			},
			_nil,
			true,
		},
		test{
			func() (interface{}, error) {
				return none.First(isFalse)
			},
			_nil,
			true,
		},
		test{
			func() (interface{}, error) {
				return none.First(isTrue)
			},
			_nil,
			true,
		},
	}

	tests["Sort"] = []test{
		test{
			func() (interface{}, error) {
				return many.Sort(byTitle), nil
			},
			Movies{fifth, first, fourth, second, third},
			false,
		},
		test{
			func() (interface{}, error) {
				return many.Sort(byTheaters), nil
			},
			Movies{third, first, second, fifth, fourth},
			false,
		},
		test{
			func() (interface{}, error) {
				return many.IsSorted(byTitle), nil
			},
			false,
			false,
		},
		test{
			func() (interface{}, error) {
				sorted := many.Sort(byTitle)
				return sorted.IsSorted(byTitle), nil
			},
			true,
			false,
		},
		test{
			func() (interface{}, error) {
				return many.SortDesc(byTitle), nil
			},
			Movies{third, second, fourth, first, fifth},
			false,
		},
		test{
			func() (interface{}, error) {
				return many.SortDesc(byTheaters), nil
			},
			Movies{fourth, fifth, second, first, third},
			false,
		},
		test{
			func() (interface{}, error) {
				return many.IsSortedDesc(byTitle), nil
			},
			false,
			false,
		},
		test{
			func() (interface{}, error) {
				sorted := many.SortDesc(byTitle)
				return sorted.IsSortedDesc(byTitle), nil
			},
			true,
			false,
		},
	}

	tests["Where"] = []test{
		test{
			func() (interface{}, error) {
				return many.Where(isDummy), nil
			},
			Movies{},
			false,
		},
		test{
			func() (interface{}, error) {
				return many.Where(isFirst), nil
			},
			Movies{first},
			false,
		},
		test{
			func() (interface{}, error) {
				return many.Where(isFirstOrThird), nil
			},
			Movies{first, third},
			false,
		},
		test{
			func() (interface{}, error) {
				return many.Where(isTrue), nil
			},
			many,
			false,
		},
		test{
			func() (interface{}, error) {
				return none.Where(isFalse), nil
			},
			Movies{},
			false,
		},
		test{
			func() (interface{}, error) {
				return none.Where(isTrue), nil
			},
			Movies{},
			false,
		},
	}

	return tests
}

func TestAll(t *testing.T) {
	checkErr := func(_t test, err error) {
		if _t.ErrorExpected && err == nil {
			t.Errorf("Expected error but did not receive one")
		}
		if !_t.ErrorExpected && err != nil {
			t.Errorf("Did not expect error but received: %v", err)
		}
	}

	for name, tests := range getTests() {
		for i, test := range tests {
			switch test.Expected.(type) {
			default:
				got, err := test.Exec()

				checkErr(test, err)

				if got != test.Expected {
					t.Errorf("%s[%v]: Expected %v, got %v", name, i, test.Expected, got)
				}
			case []string:
				_got, err := test.Exec()

				checkErr(test, err)

				got := _got.([]string)
				exp := test.Expected.([]string)
				if len(got) != len(exp) {
					t.Errorf("%s[%v]: Expected %v strings, got %v", name, i, len(exp), len(got))
					break
				}
				for j := range got {
					if got[j] != exp[j] {
						t.Errorf("%s[%v]: Expected %v, got %v", name, i, exp[j], got[j])
						break
					}
				}
			case map[int]Movies:
				_got, err := test.Exec()

				checkErr(test, err)

				got := _got.(map[int]Movies)
				exp := test.Expected.(map[int]Movies)
				if len(got) != len(exp) {
					t.Errorf("%s[%v]: Expected %v groups, got %v", name, i, len(exp), len(got))
					break
				}
				for k, _ := range got {
					got2 := got[k]
					exp2 := exp[k]
					if len(got2) != len(exp2) {
						t.Errorf("%s[%v]: Expected %v Movies in %d element, got %v", name, i, len(exp2), k, len(got2))
						break
					}
					for i := range got2 {
						if got2[i] != exp2[i] {
							t.Errorf("%s[%v]: Expected %v, got %v", name, i, exp2[i], got2[i])
						}
					}
				}
			case map[string]Movies:
				_got, err := test.Exec()

				checkErr(test, err)

				got := _got.(map[string]Movies)
				exp := test.Expected.(map[string]Movies)
				if len(got) != len(exp) {
					t.Errorf("%s[%v]: Expected %v groups, got %v", name, i, len(exp), len(got))
					break
				}
				for k, _ := range got {
					got2 := got[k]
					exp2 := exp[k]
					if len(got2) != len(exp2) {
						t.Errorf("%s[%v]: Expected %v Movies in %s element, got %v", name, i, len(exp2), k, len(got2))
						break
					}
					for i := range got2 {
						if got2[i] != exp2[i] {
							t.Errorf("%s[%v]: Expected %v, got %v", name, i, exp2[i], got2[i])
						}
					}
				}
			case Movies:
				_got, err := test.Exec()

				checkErr(test, err)

				got := _got.(Movies)
				exp := test.Expected.(Movies)
				if len(got) != len(exp) {
					t.Errorf("%s[%v]: Expected %v Movies, got %v", name, i, len(exp), len(got))
					break
				}
				for i := range got {
					if got[i] != exp[i] {
						t.Errorf("%s[%v]: Expected %v, got %v", name, i, exp[i], got[i])
					}
				}
			}
		}
	}
}

var _nil *Movie

var first = &Movie{Title: "first", Theaters: 6, Studio: "Miramax", BoxOfficeMillions: 90}
var second = &Movie{Title: "second", Theaters: 9, Studio: "Warner Bros", BoxOfficeMillions: 100}
var third = &Movie{Title: "third", Theaters: 5, Studio: "Universal", BoxOfficeMillions: 100}
var fourth = &Movie{Title: "fourth", Theaters: 50, Studio: "Universal", BoxOfficeMillions: 90}
var fifth = &Movie{Title: "fifth", Theaters: 20, Studio: "Miramax", BoxOfficeMillions: 100}

var none = Movies{}

var many = Movies{
	first,
	second,
	third,
	fourth,
	fifth,
}

var isFirst = func(movie *Movie) bool {
	return movie.Title == "first"
}
var isThird = func(movie *Movie) bool {
	return movie.Title == "third"
}
var isFirstOrSecondOrThirdOrFourthOrFifth = func(movie *Movie) bool {
	return movie.Title == "first" || movie.Title == "second" || movie.Title == "third" || movie.Title == "fourth" || movie.Title == "fifth"
}
var isFirstOrThird = func(movie *Movie) bool {
	return movie.Title == "first" || movie.Title == "third"
}
var isDummy = func(movie *Movie) bool {
	return movie.Title == "dummy"
}
var isTrue = func(movie *Movie) bool {
	return true
}
var isFalse = func(movie *Movie) bool {
	return false
}
var sumInt = func(state int, value int) int {
	return state + value
}
var sameTitle = func(a *Movie, b *Movie) bool {
	return a.Title == b.Title
}
var sameMillions = func(a *Movie, b *Movie) bool {
	return a.BoxOfficeMillions == b.BoxOfficeMillions
}
var byTitle = func(a, b *Movie) bool {
	return a.Title < b.Title
}
var byTheaters = func(a, b *Movie) bool {
	return a.Theaters < b.Theaters
}
