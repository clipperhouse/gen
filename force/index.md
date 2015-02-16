---
layout: default
title: "force flag"
path: "/force"
order: 1
---

### Tolerating type errors

`gen` operates by parsing and type-checking your source code. The correctness of your code will affect gen's ability to operate.

There are scenarios where one can get a bit stuck, however. For example, if you generate incorrect code, your codebase is now incorrect -- which will cause gen to refuse to generate new code! (Most users won't experience this, but if you are writing your own typewriter, it will sound familiar.)

There are classes of errors that gen can tolerate while maintaining correctness. For example:

	// +gen slice:"Where"
	type Foo int

	var Bar DoesntExist

Here, the code is syntactally correct and the `Foo` type is fully valid. However, `DoesntExist` is an unknown type, and the code won't build.

Under normal circumstances, gen will stop in the presence of any type error such as the one above. However, one can imagine that gen *should* logically be able to operate here, since it only cares about `Foo`.

The solution to this is the `-f` flag, which means "force". In this mode, type-check errors such as the one above **will be ignored** by gen. gen will do its best to continue in their presence.

Officially gen considers "force" mode to be undefined behavior. There are many classes of errors and we can't guarantee the correct behavior for all of them.

However, for certain classes of type-check errors, `-f` is a good way to get yourself out of a jam.

*Here's a bit of [background discussion](https://github.com/clipperhouse/gen/issues/77).*
