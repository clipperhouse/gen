## What’s this?

`gen` is a code-generation tool for Go. It’s intended to offer generics-like functionality on your types. Out of the box, it offers offers LINQ/underscore-inspired methods.

It also offers third-party, runtime extensibility via [typewriters](https://github.com/clipperhouse/typewriter).

#### [Introduction and docs…](http://clipperhouse.github.io/gen/)

[Changelog](https://github.com/clipperhouse/gen/blob/master/CHANGELOG.md)

[Hey, a video](https://www.youtube.com/watch?v=KY8OXFi3CDU)

### Typewriters
There is a list of open-source typewriters in [TYPEWRITERS.md](https://github.com/clipperhouse/gen/blob/master/TYPEWRITERS.md). Please add your own.

### Contributing

There are three big parts of `gen`.

#### gen

This repository. The gen package is primarily the command-line interface. Most of the work is done by the typewriter package, and individual typewriters.

#### typewriter

The [typewriter package](https://github.com/clipperhouse/typewriter) is where most of the parsing, type evaluation and code generation architecture lives.

#### typewriters

Typewriters are where templates and logic live for generating code. Here’s [set](https://github.com/clipperhouse/set), which will make a lovely Set container for your type. Here’s [slice](https://github.com/clipperhouse/slice), which provides the built-in LINQ-like functionality. Here’s [stringer](https://github.com/clipperhouse/stringer), a fork of Rob Pike’s [tool](https://godoc.org/golang.org/x/tools/cmd/stringer).

Third-party typewriters are added easily by the end user. You publish them as Go packages for import. [Learn more...](https://clipperhouse.github.io/gen/typewriters/)

We’d love to see typewriter packages for things like strongly-typed JSON serialization, `Queue`s, `Pool`s or other containers. Anything “of T” is a candidate for a typewriter.
