package container

import "github.com/clipperhouse/gen/typewriter"

var list = &typewriter.Template{
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
`}
