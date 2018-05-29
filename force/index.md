---
layout: default
title: "force flag"
path: "/force"
order: 1
redirect_to: "http://clipperhouse.com/gen/force/"
---

### Tolerating type errors

`gen` operates by parsing and type-checking your source code. The correctness of your code will affect gen's ability to operate.

There are scenarios where one can get a bit stuck, however. For example, if you generate incorrect code, your codebase is now incorrect -- which will cause gen to refuse to generate new code! (Most users won't experience this, but if you are authoring your own typewriter, it will sound familiar.)

There are classes of errors that gen can tolerate while maintaining correctness. For example:

	// +gen slice:"Where"
	type Foo int

	var Bar DoesntExist

Here, the code is syntactically valid and the `Foo` type is fully understood. However, `DoesntExist` is an unknown type, so the code won't build.

Under normal circumstances, gen will stop in the presence of any type error such as the one above. One can imagine that gen *should* logically be able to operate here, since it only cares about `Foo`.

The solution to this is the `-f` flag, which means "force". In this mode, type-check errors such as the one above **will be ignored** by gen. gen will do its best to continue in their presence.

Officially gen considers "force" mode to be unspecified behavior. There are many classes of errors and we can't guarantee the correct behavior under all circumstances.

However, for certain classes of type-check errors, `-f` is a good way to get yourself out of a jam.

*Here's a bit of [background discussion](https://github.com/clipperhouse/gen/issues/77).*
