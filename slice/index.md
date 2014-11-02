---
layout: default
title: "slice"
path: "/slice"
order: 2
---

## The `slice` typewriter

The `slice` typewriter is built into [`gen`](../) by default. It generates functional convenience methods that will look familiar to users of C#'s LINQ or JavaScript's Array. It is intended to save you some loops, using a "pass a function" pattern. It offers easier ad-hoc sorts.

### Aggregate[T]

Iterates over a slice, operating on each element while maintaining ‘state’. Comparable to Linq’s Aggregate or underscore’s reduce.

Annotation:

	// +gen slice:"Aggregate[string]"
	type Foo struct{
		Name 	   string
		Department string
	}

Signature:

	func (rcv FooSlice) AggregateString(func(state string, value Foo) string) bool

Example:

	foos := FooSlice {
		{"Alice", "Accounting"},
		{"Bob", "Back Office"},
		{"Carly", "Containers"},
	}

	join := func(state string, value Foo) string {
	    if state != "" {
	        state += ", "
	    }
	    return state + value.Name
	}

	foos.AggregateString(join) // => "Alice, Bob, Carly"

### All

Returns true if every element returns true for passed func. Comparable to Linq’s All or underscore’s every.

Annotation:

	// +gen slice:"All"
	type Foo struct {
		Age int
	}

Signature:

	func (rcv FooSlice) All(fn func(Foo) bool) bool

Example:

	foos := FooSlice {
		{250},
		{101},
		{300},
	}

	over100 := func(value Foo) bool {
	    return value.Amount > 100
	}

	foos.All(over100) // => true

### Any

### Average

### Average[T]

### Count

### Distinct

### DistinctBy

### First

### GroupBy[T]

### Max

### Max[T]

### MaxBy

### Min

### Min[T]

### MinBy

### Select

### Sort

### SortBy

### Where

Returns a new slice whose elements return true for passed func. Comparable to LINQ’s [Where](http://msdn.microsoft.com/en-us/library/bb534803(v=vs.110).aspx) and JavaScript’s [filter](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/filter).

	func (rcv ExampleSlice) Where(fn func(Example) bool) ExampleSlice

Example:

	shiny := func(p Product) bool {
		return p.Manufacturer == "Apple"
	}
	wishlist := products.Where(shiny)
