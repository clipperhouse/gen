---
layout: default
title: "slice"
path: "/slice"
order: 2
---

## The `slice` typewriter

The `slice` typewriter is built into [`gen`](../) by default. It generates functional convenience methods that will look familiar to users of C#'s LINQ or JavaScript's Array methods. It is intended to save you some loops, using a "pass a function" pattern. It offers easier ad-hoc sorts.

The annotation looks like:

	// +gen slice:"Where,GroupBy[int],Any"
	type Example struct {}

(`Example` is used as a placeholder for your type.)

A new type, `ExampleSlice`, is generated and becomes the receiver for the methods below.

### Aggregate[T]

Iterates over a slice, operating on each element while maintaining ‘state’. Comparable to LINQ’s Aggregate or underscore’s reduce.

Signature:

	func (ExampleSlice) AggregateT(func(T, Example) T) T

Example:

	// +gen slice:"Aggregate[string]"
	type Employee struct{
		Name 	   string
		Department string
	}

	employees := EmployeeSlice {
		{"Alice", "Accounting"},
		{"Bob", "Back Office"},
		{"Carly", "Containers"},
	}

	join := func(state string, e Employee) string {
	    if state != "" {
	        state += ", "
	    }
	    return state + e.Name
	}

	employees.AggregateString(join) // => "Alice, Bob, Carly"

### All

Returns true if every element returns true for passed func. Comparable to LINQ’s All or underscore’s every.

Annotation:

Signature:

	func (ExampleSlice) All(func(Example) bool) bool

Example:

	// +gen slice:"All"
	type Person struct {
		Name    string
		Present bool
	}

	gang := PersonSlice {
		{"Alice", true},
		{"Bob", false},
		{"Carly", true},
	}

	here := func(p Person) bool {
	    return p.Present
	}

	gang.All(here) // => false, Bob didn't make it

### Any

Returns true if one or more elements returns true for passed func. Comparable to LINQ’s Any or underscore’s some.

Signature:

	func (ExampleSlice) Any(func(Example) bool) bool

Example:

	// +gen slice:"Any"
	type Person struct {
		Name string
	}

	people := PersonSlice {
		{"Bueller"},
		{"Spicoli"},
		{"Mr. Hand"},
	}

	bueller := func(p Person) bool {
	    return p.Name == "Bueller"
	}

	people.Any(bueller) // => true

### Average

Sums over all elements of the slice and divides by len. Comparable to LINQ’s Average.

Signature:

	func (ExampleSlice) Average() Example  // Example must be a numeric type

Example:

	// +gen slice:"Average"
	type Celsius float64

	temps := CelsiusSlice {
		Celsius(15.1),
		Celsius(-2),
		Celsius(3.6),
	}

	temps.Average() // => 5.567

### Average[T]

Sums over all projected values of a numeric type, and divides by len.

Signature:

	func (ExampleSlice) AverageT() T  // T must be a numeric type

Example:

	// +gen slice:"Average[int]"
	type Player struct {
		Name   string
		Points int
	}

	players := PlayerSlice {
		{"Alice", 450},
		{"Bob", 100},
		{"Carly", 200},
	}

	points := func(f Example) int {
		return f.Points
	}

	players.AverageInt(points) // => 250

### Count

Returns an int representing the number of elements which return true for passed func. Comparable to LINQ’s Count.

Signature:

	func (ExampleSlice) Count(func(Example) bool) int

Example:

	// +gen slice:"Count"
	type Monster struct {
		Name  string
		Furry bool
		Fangs int
	}

	monsters := MonsterSlice {
		{"Alice", false, 0},
		{"Bob", true, 4},
		{"Carly", true, 2},
		{"Dave", false, 2},
	}

	werewolf := func(m Monster) bool {
		return m.Fangs > 0 && m.Furry
	}

	monsters.Count(werewolf) // => 2 (Bob & Carly)

### Distinct

Returns a new slice representing unique elements. Comparable to LINQ’s Distinct or underscore’s uniq.

	func (ExampleSlice) Distinct() ExampleSlice

*Constraint: Example must support [equality](https://golang.org/doc/go1#equality).*

Example:

	// +gen slice:"Distinct"
	type Hipster struct {
		FavoriteBand string
		Mustachioed  bool
		Bepectacled  bool
	}

	hipsters := HipsterSlice {
		{"Neutral Milk Hotel", true, true},
		{"Neutral Milk Hotel", true, true},
		{"Neutral Milk Hotel", true, true},
		{"Neutral Milk Hotel", true, true},
	}

	hipsters.Distinct() // => [{"Neutral Milk Hotel", true, true}]

Keep in mind that pointers and values have different notions of equality, and therefore distinctness.

### DistinctBy

Returns a new slice (plural type) representing unique elements, where equality is defined by a passed func.

func (rcv Things) DistinctBy(func(*Thing, *Thing) bool) Things
Example:

hairstyle := func(a *Fashionista, b *Fashionista) bool {
    a.Hairstyle == b.Hairstyle
}
trendsetters := fashionistas.DistinctBy(hairstyle)

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
