# What’s this?

`gen` is an attempt to bring some generics-like functionality to Go, with some inspiration from C#’s Linq and JavaScript’s underscore libraries.

`gen` generates code for your types (at the command line), which implements the methods described below. This gives you static, compile-time assurances and perhaps some code-completion depending on your [editor](https://github.com/DisposaBoy/GoSublime).

The basic pattern is to pass func’s as you would pass lambdas to Linq or functions to underscore.

We’re starting with proof-of-concept basics like Any, Where, and Count, later intending to fill out the ‘family’ of map-reduce functions. This is an early prototype, caveat emptor and suggestions welcome.

<a href="http://clipperhouse.github.io/gen/">Nice docs.</a>

# Getting Started

Clone this repo, and cd into the directory. Then, `go install`, which will create a binary called `gen` that you can invoke from the command line (assuming you’ve [set up](http://golang.org/doc/install) your paths, etc).

Type `gen` to see usage.

`cd` into the /test directory and type `gen *models.Movie`. This should overwrite the [movie_gen.go](/clipperhouse/gen/blob/master/test/movie_gen.go) file that came with the repo. Have a look at the header comments.

Then, `go test`.

# Design goals

We want this library to be idiomatic, fast and as lightweight as possible. We are looking to bring a bit of terseness and clarity to operations that might otherwise require verbose loops. We also will refrain from implementing functionality that would compete for problems that Go already solves cleanly.

## Implemented so far:

- **Aggregate**: loop over a slice and operate each element, accumulating along the way; useful for e.g. concatenation or summing; comparable to underscore’s `reduce` or Linq’s `Aggregate`; implemented for `string` and `int` types
- **All**: determine if all elements of a slice return true for a passed func; comparable to underscore’s `every` or Linq’s `All`.
- **Any**: determine if one or more elements of a slice return true for a passed func; comparable to underscore’s `some` or Linq’s `Any`.
- **Count**: count elements of a slice that return true for a passed func; comparable to Linq’s `Count`.
- **Each**: apply a passed func to every element of a slice; comparable to underscore’s `each` or Linq’s `ForEach`.
- **GroupBy**: group elements into a map of slices based on a passed func; implemented for `string` and `int` as keys
- **Min/Max**: return *element* with greatest/least value based on passed comparer func; comparable to underscore’s `max/min`; not like Linq’s `Max/Min`, which return the least value (not the element with the value)
- **Sort** (+ IsSorted): reorder a slice based on passed comparer func; comparable to Linq’s `OrderBy`, with some of Go’s sort idiom; Desc versions are included as well
- **Where**: returns slice of elements that return true for a passed func. Comparable to underscore’s `filter` or Linq’s `Where`.
