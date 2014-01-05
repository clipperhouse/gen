package main

import (
	"strings"
)

type projection struct {
	Method string
	Type   string
	Parent *genSpec
}

func (p projection) String() string {
	return p.Method + strings.Title(strings.Replace(p.Type, "[]", "ArrayOf", -1)) // super hacky, TODO: really handle arbitrarty types
}
