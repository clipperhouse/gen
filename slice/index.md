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

Sums over all elements of the slice and divides by len. Returns an error on an empty slice. Comparable to LINQ’s Average.

Signature:

	func (ExampleSlice) Average() Example  // Example must be a numeric type

Example:

	// +gen slice:"Average"
	type Celsius float64

	temps := CelsiusSlice{15.1, -2, 3.6}

	temps.Average() // => 5.567, nil

### Average[T]

Returns the average **projected** value of a slice, where the projection is defined by a passed func. Returns an error on an empty slice.

Signature:

	func (ExampleSlice) AverageT(func(Example) T) (T, error)  // T must be a numeric type

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

	points := func(p Player) int {
		return p.Points
	}

	players.AverageInt(points) // => 250, nil

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

Distinct is supported only for types that support [equality](https://golang.org/doc/go1#equality). Bear in mind that pointers and values have different notions of equality, and therefore distinctness.

### DistinctBy

Returns a new slice representing unique elements, where equality is defined by a passed func.

Signature:

	func (ExampleSlice) DistinctBy(func(Example, Example) bool) ExampleSlice

Example:
	
	// +gen slice:"DistinctBy"
	type Hipster struct {
		FavoriteBand string
		Mustachioed  bool
	}

	hipsters := HipsterSlice {
		{"Neutral Milk Hotel", true},
		{"Death Cab for Cutie", true},
		{"You Probably Haven’t Heard of Them", true},
		{"Neutral Milk Hotel", false},
	}

	band := func(a Hipster, b Hipster) bool {
	    a.FavoriteBand == b.FavoriteBand
	}

	hipsters.DistinctBy(band) // => [{"Neutral Milk Hotel", true}, {"Death Cab for Cutie", true}, {"You Probably Haven’t Heard of Them", true}]

### First

Returns first element which returns true for passed func. Returns error if no elements satisfy the func. Comparable to LINQ’s First or underscore’s find.

Signature:

	func (ExampleSlice) First(func(Example) bool) (Example, error)

Example:

	// +gen slice:"First"
	type Customer struct {
		Name string
		Here bool
	}

	customers := CustomerSlice {
		{"Alice", false},
		{"Bob", true},
		{"Carly", true},
	}

	come := func(c Customer) bool {
	    return c.Here
	}

	served, err := customers.First(come) // => {"Bob", true}, nil

### GroupBy[T]

Groups elements into a map keyed by T. Comparable to LINQ’s GroupBy or underscore’s groupBy.

Signature:

	func (ExampleSlice) GroupByT (func(Example) T) map[T]ExampleSlice // => T must support equality

Example:

	// +gen slice:"GroupBy[int]"
	type Movie struct {
		Title string
		Year  int
	}

	movies := MovieSlice {
		{"Independence Day", 1996},
		{"Iron Man", 2008},
		{"Fargo", 1996},
		{"Django Unchained", 2012},
		{"WALL-E", 2008},
	}

	year := func(m Movie) int {
		return m.Year
	}

	movies.GroupByInt(year) // => { 1996: [{"Independence Day", 1996}, {"Fargo", 1996}], 2008: [{"Iron Man", 2008}, {"WALL-E", 2008}], 2012: [{"Django Unchained", 2012}] }

### Max

Returns the maximum value of a slice. Returns an error when invoked on an empty slice, an invalid operation. Comparable to LINQ’s Max.

Signature:

	func (ExampleSlice) Max() (Example, error) // => Example must be an ordered type

Example:

	// +gen slice:"Max"
	type Price float64

	prices := PriceSlice{12.34, 43.21, 23.45}

	prices.Max() // => 43.21

`Max` is only supported for ‘[ordered](http://godoc.org/code.google.com/p/go.tools/go/types#BasicInfo)’ types, i.e. those that support less than/greater than.

### Max[T]

Returns the maximum projected value of a slice, where the projection is defined by a passed func. Returns an error when invoked on an empty slice, an invalid operation. Comparable to LINQ’s Max.

Signature:

	func (ExampleSlice) MaxT(func(Example) T) (T, error) // => T must be an ordered type

Example:

	// +gen slice:"Max[Dollars]"
	type Movie struct {
		Title     string
		BoxOffice Dollars
	}

	type Dollars int

	movies := MovieSlice {
		{"Independence Day", 1000000},
		{"Iron Man", 5000000},
		{"Fargo", 3000000},
		{"Django Unchained", 9000000},
		{"WALL-E", 4000000},
	}

	box := func(e Employee) Dollars {
		return e.BoxOffice
	}

	movies.MaxDollars(box) // => 9000000

### MaxBy

Returns the element containing the maximum value, when compared to other elements using a passed func defining ‘less’. Returns an error when invoked on an empty slice, considered an invalid operation.

Signature:

	func (ExampleSlice) MaxBy(func(Example, Example) bool) (Example, error)

Example:

	// +gen slice:"MaxBy"
	type Rectangle struct {
		Width, Height int
	}

	func (r Rectangle) Area() int {
		return r.Width * r.Height
	}

	rectangles := RectangleSlice{
		{5, 4},
		{6, 7},
		{2, 3},
	}

	area := func(a, b Rectangle) bool {
	    return a.Area() < b.Area()
	}

	rectangles.MaxBy(area) // => {6, 7}

### Min

Returns the minimum value of a slice. Returns an error when invoked on an empty slice, an invalid operation. Comparable to LINQ’s Min.

Signature:

	func (ExampleSlice) Min() (Example, error) // => Example must be an ordered type

Example:

	// +gen slice:"Min"
	type Price float64

	prices := PriceSlice{12.34, 43.21, 23.45}

	prices.Min() // => 12.34

`Min` is only supported for ‘[ordered](http://godoc.org/code.google.com/p/go.tools/go/types#BasicInfo)’ types, i.e. those that support less than/greater than.

### Min[T]

Returns the minimum projected value of a slice, where the projection is defined by a passed func. Returns an error when invoked on an empty slice, an invalid operation. Comparable to LINQ’s Min.

Signature:

	func (ExampleSlice) MinT(func(Example) T) (T, error) // => T must be an ordered type

Example:

	// +gen slice:"Min[Dollars]"
	type Movie struct {
		Title     string
		BoxOffice Dollars
	}

	type Dollars int

	movies := MovieSlice {
		{"Independence Day", 1000000},
		{"Iron Man", 5000000},
		{"Fargo", 3000000},
		{"Django Unchained", 9000000},
		{"WALL-E", 4000000},
	}

	box := func(e Employee) Dollars {
		return e.BoxOffice
	}

	movies.MinDollars(box) // => 1000000

### MinBy

Returns the element containing the minimum value, when compared to other elements using a passed func defining ‘less’. Returns an error when invoked on an empty slice, considered an invalid operation.

Signature:

	func (ExampleSlice) MinBy(func(Example, Example) bool) (Example, error)

Example:

	// +gen slice:"MinBy"
	type Rectangle struct {
		Width, Height int
	}

	func (r Rectangle) Area() int {
		return r.Width * r.Height
	}

	rectangles := RectangleSlice{
		{5, 4},
		{6, 7},
		{2, 3},
	}

	area := func(a, b Rectangle) bool {
	    return a.Area() < b.Area()
	}

	rectangles.MinBy(area) // => {2, 3}

### Select[T]

Returns a projected slice given a func which maps Example to T. Comparable to LINQ’s Select or underscore’s map.

Signature:

	func (ExampleSlice) Select(func(Example) T) []T

Example:

	// +gen slice:"Select[int]"
	type Player struct {
		Name   string
		Points int
	}

	players := PlayerSlice {
		{"Alice", 450},
		{"Bob", 100},
		{"Carly", 200},
	}

	points := func(p Player) int {
		return p.Points
	}

	players.SelectInt(points) // => [450, 100, 200]

### Shuffle

Returns a new slice with the elements in a random order. Comparable to underscore’s shuffle.

Signature:

	func (ExampleSlice) Shuffle() ExampleSlice

Example:

	// +gen slice:"Shuffle"
	type Rating int

	ratings := RatingSlice{1, 2, 3, 4, 5, 6}

	ratings.Shuffle() // => {3, 6, 1, 2, 4, 5}

### Sort

Returns a new slice whose elements are sorted.

Signature:

	func (ExampleSlice) Sort() ExampleSlice

Example:

	// +gen slice:"Sort,SortDesc"
	type Rating int

	ratings := RatingSlice{5, 7, 2, 1, 9, 2}

	ratings.Sort() // => {1, 2, 2, 5, 7, 9}
	ratings.SortDesc() // => {9, 7, 5, 2, 2, 1}

`SortDesc` and `IsSorted(Desc)` are also available, and should be self-explanatory.

Sort uses Go’s sort package by implementing the interface required to use it. It is only supported for types that can be compared greater than or less than one another (‘ordered’ in Go terminology).

### SortBy

Returns a new slice whose elements are sorted based on a func defining ‘less’. The less func takes two elements, and returns true if the first element is less than the second element.

Signature:

	func (ExampleSlice) SortBy(func(Example, Example) bool) ExampleSlice

Example:

	// +gen slice:"SortBy"
	type Movie struct {
		Title string
		Year  int
	}

	movies := MovieSlice {
		{"Independence Day", 1996},
		{"Iron Man", 2008},
		{"Fargo", 1996},
		{"Django Unchained", 2012},
		{"WALL-E", 2008},
	}

	yearThenTitle := func(a, b Movie) bool {
		if a.Year == b.Year {
			return a.Title < b.Title
		}
		return a.Year < b.Year
	}

	movies.SortBy(yearThenTitle) // => [{"Fargo", 1996}, "Independence Day", 1996}, {"Iron Man", 2008}, {"WALL-E", 2008}, {"Django Unchained", 2012}]

`SortByDesc` and `IsSortedBy(Desc)` are also available, and should be self-explanatory.

### Where

Returns a new slice whose elements return true for passed func. Comparable to LINQ’s [Where](http://msdn.microsoft.com/en-us/library/bb534803(v=vs.110).aspx) and JavaScript’s [filter](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/filter).

	func (rcv ExampleSlice) Where(fn func(Example) bool) ExampleSlice

Example:

	shiny := func(p Product) bool {
		return p.Manufacturer == "Apple"
	}
	wishlist := products.Where(shiny)
