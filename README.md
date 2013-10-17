# What’s this?

`gen` is an attempt to bring some generics-like functionality to Go, with some inspiration from C#’s Linq and JavaScript’s array functions.

We’re starting with proof-of-concept basics like Any, Where, and Count, later intending to fill out the ‘family’ of map-reduce functions. This is an early prototype, caveat emptor and suggestions welcome.

# Getting Started

Clone this repo, and cd into the directory. Then, `go install`, which will create a binary called `gen` that you can invoke from the command line (assuming you’ve [set up](http://golang.org/doc/install) your paths, etc).

`cd` into the /test directory and type `gen models.*Movie`. This should overwrite the [movie_gen.go](/clipperhouse/gen/blob/master/test/movie_gen.go) file that came with the repo. Have a look at the header comments.

Then, `go test`.

# Design goals

We want this library to be idiomatic, fast and as lightweight as possible. We are looking to bring a bit of terseness and clarity to operations that might otherwise require verbose loops.

# Helping out

It’s probably a bit early for pull requests, we’re still designing the API. You can always find me @clipperhouse on GitHub and Twitter, your ideas are welcome.