---
layout: page
title: "home"
order: 1
---

# gen

### Type-driven code generation for Go

**gen** is an attempt to bring some generics-like functionality to Go. It uses type annotations to add "of &lt;T&gt;" functionality to your packages.

gen generates code for your types, at development time, using the command line. It is not an import; the generated source becomes part of your package and takes no external dependencies.

## Quick start

Of course, start by installing Go, [setting up paths](http://golang.org/doc/code.html), etc. Then:

	go get github.com/clipperhouse/gen

Create a new Go project, and `cd` into it. Create a `main.go` file and define a type.

Now, mark it up with a `+gen` annotation in an adjacent comment like so:

	// +gen slice:"Where,Count,GroupBy[string]"
	type MyType struct {}

And at the command line, simply type:

	gen

You should see a new file, named `mytype_slice.go`. Have a look around.

(The annotation syntax will look familiar to Go users, it is modeled after struct tags.)

## TypeWriters

gen is driven by "type writers" -- packages which are responsible for interpreting the annotated tags and turning them into generated code.

gen includes two built-in TypeWriters.

### `slice`

The `slice` typewriter generates functional convenience methods that will look familiar to users of C#'s LINQ or JavaScript's Array methods. It is intended to save you some loops, using a "pass a function" pattern. It offers grouping, filtering, ad-hoc sorts and projections. [Details and docs...](slice)

### `container`

The `container` typewriter implements strongly-typed versions of the List and Ring containers from the Go standard library, as well as a simple Set implementation. [Details and docs...](container)

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

You can create your own type writer by implementing the [TypeWriter interface](godoc).

As mentioned above, typewriters follow the pattern of codecs in the image package of Go's standard library. They are registered via an `init()` method.

The best thing to do is have a look at an existing implementation, [List](https://github.com/clipperhouse/typewriters/container) is straightforward.

The [typewriter package](https://github.com/clipperhouse/typewriter) handles tag parsing and type evaluation, and passes this information to your typewriter. It offers some conveniences for text templating, as well.

There aren't many third-party packages as of this writing, so don't hesistate to ask me (Matt) for help. We want to make it easy, and maybe even build an ecosystem.

We'd love to see a typewriter for strongly-typed json serialization, for example.

## FAQ

**Codegen, really?**

Yes! It felt a bit dirty to us at the beginning, too. But it turns out that a lot of actual generics implementations look a lot like code generation -- you just don't see it. (Compilers and JITs do it for you.)

Code generation removes mystery. It's just code, right there in your package. Read it. The history goes in your repo like everything else.

You get all the usual compiler checks and optimizations, of course, so gen won't introduce surprises in production.

gen is a *tool* that helps the developer produce code on their local workstation, alongside their text editor and utilities.

**Wait, is this `go generate`?**

No, [that's](https://docs.google.com/document/d/1V03LUfjSADDooDMhe-_K59EgpTEm3V8uvQRuNMAEnjg/edit) different (and very cool). `go generate` will run any command and is intended to obviate `make` files and such. `gen` is specifically about codegen for types.

The two tools are complementary.

**Can I run gen on the server?**

Like as part of the build? Sure, but that's not what it's designed around so we don't recommend it.

It's a local dev tool, not a platform or (shudder) a framework. Run it locally, test it, and commit the generated code to the repo.
