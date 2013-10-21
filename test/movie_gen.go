// gen *models.Movie
// this file was auto-generated using github.com/clipperhouse/gen
// Sun, 20 Oct 2013 03:28:37 UTC

package models

type Movies []*Movie

func (m Movies) AggregateInt(fn func(*Movie, int) int) (result int) {
	for _, _m := range m {
		result = fn(_m, result)
	}
	return result
}

func (m Movies) AggregateString(fn func(*Movie, string) string) (result string) {
	for _, _m := range m {
		result = fn(_m, result)
	}
	return result
}

func (m Movies) All(fn func(*Movie) bool) bool {
	for _, _m := range m {
		if !fn(_m) {
			return false
		}
	}
	return true
}

func (m Movies) Any(fn func(*Movie) bool) bool {
	for _, _m := range m {
		if fn(_m) {
			return true
		}
	}
	return false
}

func (m Movies) Count(fn func(*Movie) bool) int {
	var count = func(_m *Movie, acc int) int {
		if fn(_m) {
			acc++
		}
		return acc
	}
	return m.AggregateInt(count)
}

func (m Movies) Each(fn func(*Movie)) {
	for _, _m := range m {
		fn(_m)
	}
}

func (m Movies) JoinString(fn func(*Movie) string, delimiter string) string {
	var join = func(_m *Movie, acc string) string {
		if _m != m[0] {
			acc += delimiter
		}
		return acc + fn(_m)
	}
	return m.AggregateString(join)
}

func (m Movies) Skip(n int) Movies {
	if len(m) > n {
		return m[n:]
	}
	return Movies{}
}

func (m Movies) SumInt(fn func(*Movie) int) int {
	var sum = func(_m *Movie, acc int) int {
		return acc + fn(_m)
	}
	return m.AggregateInt(sum)
}

func (m Movies) Where(fn func(*Movie) bool) (result Movies) {
	for _, _m := range m {
		if fn(_m) {
			result = append(result, _m)
		}
	}
	return result
}
