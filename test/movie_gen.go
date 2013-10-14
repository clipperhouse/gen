// gen models.Movie
// this file was auto-generated using github.com/clipperhouse/gen
// Mon, 14 Oct 2013 02:46:18 UTC

package models

type Movies []*Movie

func (movies Movies) All(fn func(movie *Movie) bool) bool {
	if fn == nil {
		return true
	}
	for _, m := range movies {
		if !fn(m) {
			return false
		}
	}
	return true
}
func (movies Movies) Any(fn func(movie *Movie) bool) bool {
	if fn == nil {
		return true
	}
	for _, m := range movies {
		if fn(m) {
			return true
		}
	}
	return false
}
func (movies Movies) Count(fn func(movie *Movie) bool) (result int) {
	if fn == nil {
		return len(movies)
	}
	for _, m := range movies {
		if fn(m) {
			result++
		}
	}
	return result
}
func (movies Movies) Where(fn func(movie *Movie) bool) (result Movies) {
	for _, m := range movies {
		if fn == nil || fn(m) {
			result = append(result, m)
		}
	}
	return result
}
