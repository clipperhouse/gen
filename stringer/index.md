---
layout: default
title: "stringer"
path: "/stringer"
order: 2
---

## The `stringer` typewriter

The `stringer` typewriter is a fork of Rob Pike’s [tool](https://godoc.org/golang.org/x/tools/cmd/stringer) of the same name, which generates readable strings for consts. It is built into gen by default.

To use, mark up a type for which you have consts, using Rob’s example:

	// +gen stringer
	type Pill int

	const (
		Placebo Pill = iota
		Aspirin
		Ibuprofen
		Paracetamol
		Acetaminophen = Paracetamol
	)

Then simply `gen` your package, as normal. You should see a new file, `pill_stringer.go`, containing the String method for your type.
