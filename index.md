---
layout: default
title: "gen"
path: "/"
order: 1
---

### Type-driven code generation for Go

**gen** is an attempt to bring some generics-like functionality to Go. It uses type annotations to add "of T" functionality to your packages.

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

gen is driven by “type writers” – packages which are responsible for interpreting the annotated tags and turning them into generated code. [Learn more...](typewriters)

## Usage

Type `gen help`:

	gen           Generate files for types marked with +gen.
	gen list      List available typewriters.
	gen add       Add a third-party typewriter to the current package.
	gen get       Download and install imported typewriters. 
	              Optional flags from go get: [-d] [-fix] [-t] [-u].
	gen watch     Watch the current directory for file changes, run gen
	              when detected. 
	gen help      Print usage.

## FAQ

**Why?**

The goal of gen is not just to offer conveniences, but stronger typing. By generating strongly-typed methods and containers, we avoid having to use `interface{}`s, type assertions and reflection.

We gain compile-time safety and (perhaps) performance.

**Codegen, really?**

Yes! It felt a bit dirty to us at the beginning, too. But it turns out that a lot of actual generics implementations look a lot like code generation -- you just don't see it. (Compilers and JITs do it for you.)

Code generation removes mystery. It's just code, right there in your package. Read it. The history goes in your repo like everything else.

You get all the usual compiler checks and optimizations, of course, so gen won't introduce surprises in production.

gen is a *tool* that helps the developer produce code on their local workstation, alongside their text editor and utilities.

**Is there a video?**

[Glad you asked, yes.](https://www.youtube.com/watch?v=KY8OXFi3CDU)

**Wait, is this `go generate`?**

No, [that's](https://docs.google.com/document/d/1V03LUfjSADDooDMhe-_K59EgpTEm3V8uvQRuNMAEnjg/edit) different (and very cool). `go generate` will run any command and is intended to obviate `make` files and such. `gen` is specifically about codegen for types.

The two tools are complementary.

**Can I run gen on the server?**

Like as part of the build? Sure, but that's not what it's designed around so we don't recommend it.

It's a local dev tool, not a platform or (shudder) a framework. Run it locally, test it, and commit the generated code to the repo.
