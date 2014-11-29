---
layout: default
title: "typewriter"
path: "/typewriter"
order: 2
---

## TypeWriters

gen is driven by "type writers" -- packages which are responsible for interpreting the annotated tags and turning them into generated code.

gen includes one built-in TypeWriter:

#### `slice`

The `slice` typewriter generates functional convenience methods that will look familiar to users of C#'s LINQ or JavaScript's Array methods. It is intended to save you some loops, using a "pass a function" pattern. It offers grouping, filtering, ad-hoc sorts and projections. [Details and docs...](/slice)

### Listing typewriters

To view the currently-available typewriters, `cd` into your package and type:

	gen list

### Adding third-party TypeWriters

TypeWriters can be implemented by third-parties and used at "gen time". To use a third-party typewriter, `cd` into the root of your package and type (for example):

	gen add github.com/clipperhouse/setwriter

This will create a `_gen.go` file. Have a look at it -- it should contain imports for the built-in slicewriter (above) and your new typewriter.

Confirm it by typing:

	gen list

To ensure you've got your third-party packages locally, type:

	gen get

Now use the third-party tag (`set` in this case) to your type annotation:

	// +gen set slice:"Where,Count,GroupBy[string]"
	type MyType struct {}

And run gen on your package:

	gen

You should have a new file `mytype_set.go`. Third-parties are responsible for the quality and documentation of their typewriters, of course.

To **remove** a typewriter, simply remove the import from the `_gen.go` file.

### Implementing TypeWriters

You can create your own typewriter by implementing the [TypeWriter interface](http://godoc.org/github.com/clipperhouse/typewriter#TypeWriter).

	type TypeWriter interface {
		Name() string
		// Imports is a slice of imports required for the type; each will be written into the imports declaration.
		Imports(t Type) []ImportSpec
		// Write writes to the body of the generated code, following package declaration and imports.
		Write(w io.Writer, t Type) error
	}

Typewriters follow the pattern of formats in the [image package](http://blog.golang.org/go-image-package#TOC_5.) of Go's standard library. They are registered via an `init()` method.

The best thing to do is have a look at an existing implementation, [Set](https://github.com/clipperhouse/setwriter) is straightforward.

By convention, the name of the typewriter is the name of the tag that it uses. Type `gen list` to see the name of the typewriters available in your current package.

The [typewriter package](https://github.com/clipperhouse/typewriter) handles tag parsing and type evaluation, and passes this information to your typewriter. It offers some conveniences for text templating, as well.
