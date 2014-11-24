## What’s this?

`gen` is a code-generation tool for Go. It’s intended to offer generics-like functionality on your types. Out of the box, it offers offers LINQ/underscore-inspired methods.

It also offers third-party, runtime extensibility via [typewriters](https://github.com/clipperhouse/typewriter).

####[Introduction and docs…](http://clipperhouse.github.io/gen/)

[Changelog](https://github.com/clipperhouse/gen/blob/master/CHANGELOG.md)

###Contributing

There are three big parts of `gen`.

####gen

This repository. The gen package is primarily the command-line interface. Most of the work is done by the typewriter package, and individual typwriters.

####typewriter

The [typewriter package](https://github.com/clipperhouse/typewriter) is where most of the parsing, type evaluation and code generation architecture lives.

####typewriters

Typewriters are where templates and logic live for generating code. Here’s [setwriter](https://github.com/clipperhouse/setwriter), which will make a lovely Set container for your type. Here’s [slicewriter](https://github.com/clipperhouse/slicewriter), which provides the built-in LINQ-like functionality.

Third-party typewriters are added easily by the end user. You publish them as Go packages for import. [Learn more].

We’d love to see typewriter packages for things like strongly-typed JSON serialization, `Queue`s, `Pool`s or other containers. Anything “of T” is a candidate for a typewriter.