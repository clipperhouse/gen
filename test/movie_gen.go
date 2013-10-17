// gen *models.Movie
// this file was auto-generated using github.com/clipperhouse/gen
// Thu, 17 Oct 2013 20:05:11 UTC

package models

type Movies []*Movie

func (m Movies) AggregateInt(fn func(*Movie, int) int) (result int) {
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

func (m Movies) Count(fn func(*Movie) bool) (result int) {
	for _, _m := range m {
		if fn(_m) {
			result++
		}
	}
	return result
}

func (m Movies) Each(fn func(*Movie)) {
	for _, _m := range m {
		fn(_m)
	}
}

func (m Movies) SumInt(fn func(*Movie) int) (result int) {
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
