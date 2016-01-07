##Typewriters

This is a list of known open-source typewriters in alphabetical order. 
GoDoc links are to the original implementation that might use `interface{}` instead of a gen type.
Please add your own by making a pull request.

#### LinkedList [![GoDoc](https://godoc.org/container/list?status.svg)](https://godoc.org/container/list)
`gen add github.com/clipperhouse/linkedlist`

```go
// +gen list
type MyType struct{}
```

Implements a strongly-typed, doubly-linked list, based on [golang.org/pkg/container/list](https://golang.org/pkg/container/list). 


#### Ring [![GoDoc](https://godoc.org/container/ring?status.svg)](https://godoc.org/container/ring)
`gen add github.com/clipperhouse/ring`

```go
// +gen ring
type MyType struct{}
```

Implements strongly-typed operations on circular lists, based on [golang.org/pkg/container/ring](https://golang.org/pkg/container/ring). 


#### Set [![GoDoc](https://godoc.org/github.com/deckarep/golang-set?status.svg)](https://godoc.org/github.com/deckarep/golang-set)
`gen add github.com/clipperhouse/set`  

```go
// +gen set
type MyType struct{}
```
Implements a strongly-typed unordered set with unique values, based on [github.com/deckarep/golang-set](https://github.com/deckarep/golang-set).


#### Stringer [![GoDoc](https://godoc.org/golang.org/x/tools/cmd/stringer?status.svg)](https://godoc.org/golang.org/x/tools/cmd/stringer)
`github.com/clipperhouse/stringer` `built-in typewriter, no need to install`  

```go
// +gen stringer
type Pill int

const (
    Placebo Pill = iota
    Aspirin
    Ibuprofen
    Paracetamol
    Acetaminophen = Paracetamol
)
```
Generates an implementation of the Stringer interface from const variable names, for pretty-printing, based on  [golang.org/x/tools/cmd/stringer](https://golang.org/x/tools/cmd/stringer).

