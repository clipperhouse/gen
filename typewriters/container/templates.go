package container

import (
	"github.com/clipperhouse/gen/typewriter"
)

var templates = typewriter.TemplateSet{

	"List": &typewriter.Template{
		Text: `
// {{.Name}}Element is an element of a linked list.
type {{.Name}}Element struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *{{.Name}}Element

	// The list to which this element belongs.
	list *{{.Name}}List

	// The value stored with this element.
	Value {{.Pointer}}{{.Name}}
}

// Next returns the next list element or nil.
func (e *{{.Name}}Element) Next() *{{.Name}}Element {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list element or nil.
func (e *{{.Name}}Element) Prev() *{{.Name}}Element {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// {{.Name}}List represents a doubly linked list.
// The zero value for {{.Name}}List is an empty list ready to use.
type {{.Name}}List struct {
	root {{.Name}}Element // sentinel list element, only &root, root.prev, and root.next are used
	len  int     // current list length excluding (this) sentinel element
}

// Init initializes or clears list l.
func (l *{{.Name}}List) Init() *{{.Name}}List {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// New returns an initialized list.
func New{{.Name}}List() *{{.Name}}List { return new({{.Name}}List).Init() }

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *{{.Name}}List) Len() int { return l.len }

// Front returns the first element of list l or nil.
func (l *{{.Name}}List) Front() *{{.Name}}Element {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last element of list l or nil.
func (l *{{.Name}}List) Back() *{{.Name}}Element {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// lazyInit lazily initializes a zero {{.Name}}List value.
func (l *{{.Name}}List) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// insert inserts e after at, increments l.len, and returns e.
func (l *{{.Name}}List) insert(e, at *{{.Name}}Element) *{{.Name}}Element {
	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e
	e.list = l
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&{{.Name}}Element{Value: v}, at).
func (l *{{.Name}}List) insertValue(v {{.Pointer}}{{.Name}}, at *{{.Name}}Element) *{{.Name}}Element {
	return l.insert(&{{.Name}}Element{Value: v}, at)
}

// remove removes e from its list, decrements l.len, and returns e.
func (l *{{.Name}}List) remove(e *{{.Name}}Element) *{{.Name}}Element {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
	return e
}

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value.
func (l *{{.Name}}List) Remove(e *{{.Name}}Element) {{.Pointer}}{{.Name}} {
	if e.list == l {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero {{.Name}}Element) and l.remove will crash
		l.remove(e)
	}
	return e.Value
}

// PushFront inserts a new element e with value v at the front of list l and returns e.
func (l *{{.Name}}List) PushFront(v {{.Pointer}}{{.Name}}) *{{.Name}}Element {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

// PushBack inserts a new element e with value v at the back of list l and returns e.
func (l *{{.Name}}List) PushBack(v {{.Pointer}}{{.Name}}) *{{.Name}}Element {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}

// InsertBefore inserts a new element e with value v immediately before mark and returns e.
// If mark is not an element of l, the list is not modified.
func (l *{{.Name}}List) InsertBefore(v {{.Pointer}}{{.Name}}, mark *{{.Name}}Element) *{{.Name}}Element {
	if mark.list != l {
		return nil
	}
	// see comment in {{.Name}}List.Remove about initialization of l
	return l.insertValue(v, mark.prev)
}

// InsertAfter inserts a new element e with value v immediately after mark and returns e.
// If mark is not an element of l, the list is not modified.
func (l *{{.Name}}List) InsertAfter(v {{.Pointer}}{{.Name}}, mark *{{.Name}}Element) *{{.Name}}Element {
	if mark.list != l {
		return nil
	}
	// see comment in {{.Name}}List.Remove about initialization of l
	return l.insertValue(v, mark)
}

// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
func (l *{{.Name}}List) MoveToFront(e *{{.Name}}Element) {
	if e.list != l || l.root.next == e {
		return
	}
	// see comment in {{.Name}}List.Remove about initialization of l
	l.insert(l.remove(e), &l.root)
}

// MoveToBack moves element e to the back of list l.
// If e is not an element of l, the list is not modified.
func (l *{{.Name}}List) MoveToBack(e *{{.Name}}Element) {
	if e.list != l || l.root.prev == e {
		return
	}
	// see comment in {{.Name}}List.Remove about initialization of l
	l.insert(l.remove(e), l.root.prev)
}

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
func (l *{{.Name}}List) MoveBefore(e, mark *{{.Name}}Element) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.insert(l.remove(e), mark.prev)
}

// MoveAfter moves element e to its new position after mark.
// If e is not an element of l, or e == mark, the list is not modified.
func (l *{{.Name}}List) MoveAfter(e, mark *{{.Name}}Element) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.insert(l.remove(e), mark)
}

// PushBackList inserts a copy of an other list at the back of list l.
// The lists l and other may be the same.
func (l *{{.Name}}List) PushBackList(other *{{.Name}}List) {
	l.lazyInit()
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Value, l.root.prev)
	}
}

// PushFrontList inserts a copy of an other list at the front of list l.
// The lists l and other may be the same.
func (l *{{.Name}}List) PushFrontList(other *{{.Name}}List) {
	l.lazyInit()
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Value, &l.root)
	}
}
`},
	"Ring": &typewriter.Template{
		Text: `
// A Ring is an element of a circular list, or ring.
// Rings do not have a beginning or end; a pointer to any ring element
// serves as reference to the entire ring. Empty rings are represented
// as nil Ring pointers. The zero value for a Ring is a one-element
// ring with a nil Value.
//
type {{.Name}}Ring struct {
	next, prev *{{.Name}}Ring
	Value      {{.Pointer}}{{.Name}} // for use by client; untouched by this library
}

func (r *{{.Name}}Ring) init() *{{.Name}}Ring {
	r.next = r
	r.prev = r
	return r
}

// Next returns the next ring element. r must not be empty.
func (r *{{.Name}}Ring) Next() *{{.Name}}Ring {
	if r.next == nil {
		return r.init()
	}
	return r.next
}

// Prev returns the previous ring element. r must not be empty.
func (r *{{.Name}}Ring) Prev() *{{.Name}}Ring {
	if r.next == nil {
		return r.init()
	}
	return r.prev
}

// Move moves n % r.Len() elements backward (n < 0) or forward (n >= 0)
// in the ring and returns that ring element. r must not be empty.
//
func (r *{{.Name}}Ring) Move(n int) *{{.Name}}Ring {
	if r.next == nil {
		return r.init()
	}
	switch {
	case n < 0:
		for ; n < 0; n++ {
			r = r.prev
		}
	case n > 0:
		for ; n > 0; n-- {
			r = r.next
		}
	}
	return r
}

// New creates a ring of n elements.
func New{{.Name}}Ring(n int) *{{.Name}}Ring {
	if n <= 0 {
		return nil
	}
	r := new({{.Name}}Ring)
	p := r
	for i := 1; i < n; i++ {
		p.next = &{{.Name}}Ring{prev: p}
		p = p.next
	}
	p.next = r
	r.prev = p
	return r
}

// Link connects ring r with ring s such that r.Next()
// becomes s and returns the original value for r.Next().
// r must not be empty.
//
// If r and s point to the same ring, linking
// them removes the elements between r and s from the ring.
// The removed elements form a subring and the result is a
// reference to that subring (if no elements were removed,
// the result is still the original value for r.Next(),
// and not nil).
//
// If r and s point to different rings, linking
// them creates a single ring with the elements of s inserted
// after r. The result points to the element following the
// last element of s after insertion.
//
func (r *{{.Name}}Ring) Link(s *{{.Name}}Ring) *{{.Name}}Ring {
	n := r.Next()
	if s != nil {
		p := s.Prev()
		// Note: Cannot use multiple assignment because
		// evaluation order of LHS is not specified.
		r.next = s
		s.prev = r
		n.prev = p
		p.next = n
	}
	return n
}

// Unlink removes n % r.Len() elements from the ring r, starting
// at r.Next(). If n % r.Len() == 0, r remains unchanged.
// The result is the removed subring. r must not be empty.
//
func (r *{{.Name}}Ring) Unlink(n int) *{{.Name}}Ring {
	if n <= 0 {
		return nil
	}
	return r.Link(r.Move(n + 1))
}

// Len computes the number of elements in ring r.
// It executes in time proportional to the number of elements.
//
func (r *{{.Name}}Ring) Len() int {
	n := 0
	if r != nil {
		n = 1
		for p := r.Next(); p != r; p = p.next {
			n++
		}
	}
	return n
}

// Do calls function f on each element of the ring, in forward order.
// The behavior of Do is undefined if f changes *r.
func (r *{{.Name}}Ring) Do(f func({{.Pointer}}{{.Name}})) {
	if r != nil {
		f(r.Value)
		for p := r.Next(); p != r; p = p.next {
			f(p.Value)
		}
	}
}
`},
	"Set": &typewriter.Template{
		Text: `
// The primary type that represents a set
type {{.Name}}Set map[{{.Pointer}}{{.Name}}]struct{}

// Creates and returns a reference to an empty set.
func New{{.Name}}Set() {{.Name}}Set {
	return make({{.Name}}Set)
}

// Creates and returns a reference to a set from an existing slice
func New{{.Name}}SetFromSlice(s []{{.Pointer}}{{.Name}}) {{.Name}}Set {
	a := New{{.Name}}Set()
	for _, item := range s {
		a.Add(item)
	}
	return a
}

// Adds an item to the current set if it doesn't already exist in the set.
func (set {{.Name}}Set) Add(i {{.Pointer}}{{.Name}}) bool {
	_, found := set[i]
	set[i] = struct{}{}
	return !found //False if it existed already
}

// Determines if a given item is already in the set.
func (set {{.Name}}Set) Contains(i {{.Pointer}}{{.Name}}) bool {
	_, found := set[i]
	return found
}

// Determines if the given items are all in the set
func (set {{.Name}}Set) ContainsAll(i ...{{.Pointer}}{{.Name}}) bool {
	allSet := New{{.Name}}SetFromSlice(i)
	if allSet.IsSubset(set) {
		return true
	}
	return false
}

// Determines if every item in the other set is in this set.
func (set {{.Name}}Set) IsSubset(other {{.Name}}Set) bool {
	for elem := range set {
		if !other.Contains(elem) {
			return false
		}
	}
	return true
}

// Determines if every item of this set is in the other set.
func (set {{.Name}}Set) IsSuperset(other {{.Name}}Set) bool {
	return other.IsSubset(set)
}

// Returns a new set with all items in both sets.
func (set {{.Name}}Set) Union(other {{.Name}}Set) {{.Name}}Set {
	unionedSet := New{{.Name}}Set()

	for elem := range set {
		unionedSet.Add(elem)
	}
	for elem := range other {
		unionedSet.Add(elem)
	}
	return unionedSet
}

// Returns a new set with items that exist only in both sets.
func (set {{.Name}}Set) Intersect(other {{.Name}}Set) {{.Name}}Set {
	intersection := New{{.Name}}Set()
	// loop over smaller set
	if set.Cardinality() < other.Cardinality() {
		for elem := range set {
			if other.Contains(elem) {
				intersection.Add(elem)
			}
		}
	} else {
		for elem := range other {
			if set.Contains(elem) {
				intersection.Add(elem)
			}
		}
	}
	return intersection
}

// Returns a new set with items in the current set but not in the other set
func (set {{.Name}}Set) Difference(other {{.Name}}Set) {{.Name}}Set {
	differencedSet := New{{.Name}}Set()
	for elem := range set {
		if !other.Contains(elem) {
			differencedSet.Add(elem)
		}
	}
	return differencedSet
}

// Returns a new set with items in the current set or the other set but not in both.
func (set {{.Name}}Set) SymmetricDifference(other {{.Name}}Set) {{.Name}}Set {
	aDiff := set.Difference(other)
	bDiff := other.Difference(set)
	return aDiff.Union(bDiff)
}

// Clears the entire set to be the empty set.
func (set *{{.Name}}Set) Clear() {
	*set = make({{.Name}}Set)
}

// Allows the removal of a single item in the set.
func (set {{.Name}}Set) Remove(i {{.Pointer}}{{.Name}}) {
	delete(set, i)
}

// Cardinality returns how many items are currently in the set.
func (set {{.Name}}Set) Cardinality() int {
	return len(set)
}

// Iter() returns a channel of type {{.Pointer}}{{.Name}} that you can range over.
func (set {{.Name}}Set) Iter() <-chan {{.Pointer}}{{.Name}} {
	ch := make(chan {{.Pointer}}{{.Name}})
	go func() {
		for elem := range set {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

// Equal determines if two sets are equal to each other.
// If they both are the same size and have the same items they are considered equal.
// Order of items is not relevent for sets to be equal.
func (set {{.Name}}Set) Equal(other {{.Name}}Set) bool {
	if set.Cardinality() != other.Cardinality() {
		return false
	}
	for elem := range set {
		if !other.Contains(elem) {
			return false
		}
	}
	return true
}

// Returns a clone of the set.
// Does NOT clone the underlying elements.
func (set {{.Name}}Set) Clone() {{.Name}}Set {
	clonedSet := New{{.Name}}Set()
	for elem := range set {
		clonedSet.Add(elem)
	}
	return clonedSet
}
`,
		RequiresComparable: true,
	},
}
