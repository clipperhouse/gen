## What’s this?

`gen` is an attempt to bring some generics-like functionality to Go, with some inspiration from C#’s Linq and JavaScript’s underscore libraries. It’s currently an early alpha.

####<a href="http://clipperhouse.github.io/gen/">Introduction and docs...</a>

## Contributing

It’s early days and the API is likely volatile, ideas and contributions are welcome. Have a look at the design principles below. Feel free to [open an issue](//github.com/clipperhouse/gen/issues), send a pull request, or ping Matt Sherman [@clipperhouse](http://twitter.com/clipperhouse).

## Design principles for contributors

This library exists to provide readability and terseness to users’ code. It’s intended to reduce the number of explicit loops, by instead passing func’s as you would with C#’s Linq, JavaScript’s Array methods, or the underscore library. If it feels like piping, that’s good.

It’s intended to fit well with idiomatic Go. Explicitness and compile-time safety are preferred. For this reason, we are not using interfaces or run-time reflection. (Though if a good case can be made, we’ll listen.)

The goal is to keep the API small. We aim to implement the **least number of orthogonal methods** which allow the desired range of function.

We avoid methods that feel like wrappers or aliases to existing methods, even if they are convenient. A good proxy is to imagine a user asking the question ‘which method should I use?’. If that’s a reasonable question, the library should be doing less.

We avoid methods that would require many boilerplate implementations by type. For example, we originally had things like GroupByString and SumInt. Clearly, this would lead to requests for SumFloat64 and GroupByRune or what have you. We solved this particular problem by using [custom methods](http://clipperhouse.github.io/gen/#Custom).

These guidelines are not entirely deterministic! There’s lots of room for judgment and taste, and we look forward to seeing how it evolves.
