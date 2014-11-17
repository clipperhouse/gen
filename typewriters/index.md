---
layout: default
title: "typewriter"
path: "/typewriter"
order: 2
---

## TypeWriters

gen is driven by "type writers" -- packages which are responsible for interpreting the annotated tags and turning them into generated code.

gen includes two built-in TypeWriters.

#### `slice`

The `slice` typewriter generates functional convenience methods that will look familiar to users of C#'s LINQ or JavaScript's Array methods. It is intended to save you some loops, using a "pass a function" pattern. It offers grouping, filtering, ad-hoc sorts and projections. [Details and docs...](/slice)

#### `container`

The `container` typewriter implements strongly-typed versions of the List and Ring containers from the Go standard library, as well as a simple Set implementation. [Details and docs...](/container)

### Listing typewriters

To view the currently-available typewriters, `cd` into your package and type:

	gen list

### Using third-party TypeWriters

TypeWriters can be implemented by third-parties and used at "gen time". To use a third-party typewriter, `cd` into the root of your package and type:

	gen custom

This will create a `_gen.go` file. Have a look at it -- it contains the two built-in typewriters (above) as imports.

You'll note that typewriters are imported similarly to [codecs](http://golang.org/pkg/image/png/) in the image package, or [drivers](http://golang.org/pkg/database/sql/driver/) in the sql package. If you know the import path of a third-party typewriter, add the import:

	. "gitplace.com/kimye/bling"

Now, make sure to type:

	gen get

Add the third-party tag to your type annotation:

	// +gen slice:"Where,Count,GroupBy[string]" bling:"Diamonds,Bentleys"
	type MyType struct {}

And gen it again:

	gen

You should have a new file `mytype_bling.go` (assuming the author followed naming conventions). Refer to the third-party's documentation to understand what their tags do, of course.

### Implementing TypeWriters

You can create your own typewriter by implementing the [TypeWriter interface](http://godoc.org/github.com/clipperhouse/typewriter#TypeWriter).

Typewriters follow the pattern of formats in the [image package](http://blog.golang.org/go-image-package#TOC_5.) of Go's standard library. They are registered via an `init()` method.

The best thing to do is have a look at an existing implementation, [List](https://github.com/clipperhouse/containerwriter) is straightforward.

The [typewriter package](https://github.com/clipperhouse/typewriter) handles tag parsing and type evaluation, and passes this information to your typewriter. It offers some conveniences for text templating, as well.

There aren't many third-party packages as of this writing, so don't hesistate to ask me (Matt) for help. We want to make it easy, and maybe even build an ecosystem.

We'd love to see a typewriter for strongly-typed json serialization, for example.

