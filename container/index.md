---
layout: default
title: "container"
path: "/container"
order: 3
---

## The `container` typewriter

The `container` typewriter is built into [`gen`](../) by default. It implements strongly-typed versions of the List and Ring containers from the Go standard library, as well as a simple Set implementation.

Containers are specified using the containers tag.

	// +gen containers:"Set,List,Ring"
	type Thing struct{}

### List

	type ThingList struct

Implements a strongly-typed, doubly-linked list, based on [golang.org/pkg/container/list](https://golang.org/pkg/container/list). API documentation is available at that link. Parameters and return values that would be interface{} in the golang implementation will instead use your type in the gen implementation.

### Ring

	type ThingRing struct

Implements strongly-typed operations on circular lists, based on [golang.org/pkg/container/ring](https://golang.org/pkg/container/ring). API documentation is available at that link. Parameters and return values that would be interface{} in the golang implementation will instead use your type in the gen implementation.

### Set

	type ExampleSet map[Example]struct{}

Implements a strongly-typed set with common [operations](http://godoc.org/github.com/deckarep/golang-set) (Union, Difference, etc). Items stored within it are unordered and unique.

The implementation is based on [github.com/deckarep/golang-set](https://github.com/deckarep/golang-set), with permission. API documentation is available here. Parameters and return values that would be `interface{}` in the @deckarep implementation will instead use your type in the gen implementation.