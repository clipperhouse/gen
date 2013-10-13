// gen models.Movie
// this file was auto-generated using github.com/clipperhouse/gen
// Sun, 13 Oct 2013 23:54:35 UTC

package models

type Movies []*Movie

func (movies Movies) All(fn func(movie *Movie) bool) bool {
	for _, m := range movies {
		if !fn(m) {
			return false
		}
	}
	return true
}
func (movies Movies) Any(fn func(movie *Movie) bool) bool {
	for _, m := range movies {
		if fn(m) {
			return true
		}
	}
	return false
}
func (movies Movies) Count(fn func(movie *Movie) bool) (result int) {
	for _, m := range movies {
		if fn(m) {
			result++
		}
	}
	return result
}
func (movies Movies) Where(fn func(movie *Movie) bool) (result Movies) {
	for _, m := range movies {
		if fn(m) {
			result = append(result, m)
		}
	}
	return result
}
