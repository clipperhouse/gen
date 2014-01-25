## Changelog

**Date TBD, on projection branch as of this writing. Update on merge with master.**

Command-line type specification has been deprecated, and replaced by markup per https://github.com/clipperhouse/gen/issues/23. It takes the form of:

```
// +gen
type MyType struct {...}
```
Where before the command-line would be `$ gen package.Type`, now it's simply `$ gen`, which will locate and process your marked-up types.

Here's a larger example:
```
// +gen * methods:"Count,Where" projections:"SomeType,int"
type MyType struct {...}
```
- The `*` is a directive to generate methods which take a pointer type instead of a value type. Optional but recommended.
- The `methods` tag is for [subsetting](http://clipperhouse.github.io/gen/#Subsetting) methods; replaces `gen:"..."`. Optional, omit it to simply generate all standard methods.
- The `projections` tag is for specifying types to be projected for methods such as `Select` and `GroupBy`. Optional. If the `methods` tag is omitted, all projection methods will be generated, appropriate to each type. (For example, Average will not be generated for non-numeric types.) Otherwise, subset projection methods using the `methods` tag.

The `-all` flag has been deprecated, it's no longer a valid use case, given the above. The `-exported` flag, which is a modifier of same, is gone too.

Custom methods, where specific member fields of struct are marked up, have been deprecated. See the above-referenced issue for [rationale](https://github.com/clipperhouse/gen/issues/23).

`Sort(func)` has been renamed `SortBy(func)`, and similarly Max → MaxBy, Min → MinBy. This is done in anticipation of methods of those names which will *not* take a func, see [this issue](https://github.com/clipperhouse/gen/issues/28).