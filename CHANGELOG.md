## Changelog

**26 Jan 2014**

This release includes breaking changes.

Command-line type specification has been deprecated, and replaced by markup per [#23](https://github.com/clipperhouse/gen/issues/23). It takes the form of:

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
- The `methods` tag is for [subsetting](http://clipperhouse.github.io/gen/#Subsetting) methods; replaces `gen:"..."`. Optional; omit it to generate all standard methods.
- The `projections` tag specifies types to be [projected](http://clipperhouse.github.io/gen/#Projections) for methods such as `Select` and `GroupBy`. Optional. If the `methods` tag is omitted, all projection methods will be generated, appropriate to each type. (For example, Average will not be generated for non-numeric types.) You can subset projection methods using the `methods` tag above.

The `-all` flag has been deprecated, it's no longer a valid use case, given the above. The `-exported` flag, which is a modifier of same, is gone too.

Custom methods, where specific member fields of a struct are marked up, have been deprecated. The [rationale](https://github.com/clipperhouse/gen/issues/23) is that we prefer to project types, not fields.

`Sort(func)` has been renamed `SortBy(func)`, and similarly Max → MaxBy, Min → MinBy. This is done in anticipation of methods of those names which will *not* take a func, see [#28](https://github.com/clipperhouse/gen/issues/28).