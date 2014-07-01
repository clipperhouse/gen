###28 Jun 2014

To get the latest: `go get -u github.com/clipperhouse/gen`. Type `gen help` to see commands.

This release introduces the optional `_gen.go` file for importing custom typewriters.

Prior to this release, typewriters were simply part of the `gen` binary. Now, by creating a file of the above name in your package, third-party typewriters can be included at runtime.

`cd` into your package and type `gen custom`. You will see a new _gen.go file which looks like this:

```
package main

import (
	_ "github.com/clipperhouse/gen/typewriters/container"
	_ "github.com/clipperhouse/gen/typewriters/genwriter"
)
```

Change those import paths to other, third-party typewriters. Then call `gen get`.

Docs on how to create a typewriter are coming soon. In the meantime, have a look at the [container](https://github.com/clipperhouse/gen/tree/master/typewriters/container) typewriter for a decent example.

###12 Jun 2014

A new architecture based on ‘typewriters’. Ideally you should see little change, but…

This was a long-lived branch, around 6 weeks. The architecture should be a lot better, as well as the testing. But of course regressions are possible.

I will mark this with a release 2.0 tag following semver conventions, since there are some behavioral changes. If there are regressions for you, you can use gopkg.in to stick with 1.0.

A few important behavioral changes:

- Each typewriter now outputs a separate file, eg `*_gen.go` and `*_container.go`. This should not be a breaking change. Formerly, there was a single `*_gen.go` file
- Previously-gen’d files which have been ‘un-gen’d’ (i.e. removed a tag) will not be deleted as before. Do this manually if need be, but I hope it’s an edge case. Will consider adding it back in.
- The `-force` flag is gone. Officially, it was undefined behavior, intended for use if you get yourself into a jam.
- The contents of your gen’d files may be slightly different, but their behavior should be unchanged.

We are going to exploit this architecture to do some ~~evil~~ interesting things, stay tuned.

Any trouble, please let me know via a GitHub issue, or [Twitter](http://twitter.com/clipperhouse), or…

###9 Mar 2014

Preliminary support for containers, starting with Set, List and Ring. The use case is to generate a strongly-typed container, to avoid the overhead of casting (type assertions), and to add compile-time safety.

###1 Feb 2014

gen will now delete `*_gen.go` files in the case that a previously-gen’d type has been removed, per #34. Confirms with the user for safety. And of course you are using version control, right?

Added new `Max` and `Min` (alongside the existing respective *By). The difference is that these work on types that are known ordered such as type MyOrderable int. It requires no passed ‘less’ func because we already know what ‘less’ is for such types. #28.

MaxBy(less) and MinBy(less) are still the way to go for structs or ad-hoc ordering.

###30 Jan 2014

Added new `Sort` (alongside the existing `SortBy`). The difference is that Sort works on types that are known sortable such as `type MySortable int`. It requires no passed ‘less’ func because we already know what ‘less’ is for such types. :)

`SortBy(less)` is still the way to go for structs or any ad-hoc sorts.

The ‘integration tests’ (those which test the gen’d code) have been rewritten for clarity.

Generated code is now passed through the go/format package (gofmt) on output, so you don’t have to.

###26 Jan 2014

This release includes breaking changes. To update:

`go get -u github.com/clipperhouse/gen`

Command-line type specification has been deprecated, and replaced by markup per #23. It takes the form of:

```
// +gen
type MyType struct {...}
Where before the command-line would be gen package.Type, now it's simply gen, which will locate and process your marked-up types.
```
Here's a larger example:

```
// +gen * methods:"Count,Where" projections:"SomeType,int"
type MyType struct {...}
```

- The * is a directive to generate methods which take a pointer type instead of a value type. Optional but recommended.
- The methods tag is for subsetting methods; replaces gen:"...". Optional; omit it to generate all standard methods.
- The projections tag specifies types to be projected for methods such as Select and GroupBy. Optional. If the methods tag is omitted, all projection methods will be generated, appropriate to each type. (For example, Average will not be generated for non-numeric types.) You can subset projection methods using the methods tag above.
- The -all flag has been deprecated, it's no longer a valid use case, given the above. The -exported flag, which is a modifier of same, is gone too.

Custom methods, where specific member fields of a struct are marked up, have been deprecated. The rationale is that we prefer to project types, not fields.

Sort(func) has been renamed SortBy(func), and similarly Max → MaxBy, Min → MinBy. This is done in anticipation of methods of those names which will not take a func, see #28.
