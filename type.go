package main

type Type struct {
	*typeArg
	SubsettedMethods []string
	ProjectedTypes   []string
}

func newType(t *typeArg) (typ *Type) {
	typ = &Type{}
	typ.Pointer = t.Pointer
	typ.Package = t.Package
	typ.Name = t.Name
	return
}
