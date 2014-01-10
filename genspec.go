package main

import (
	"fmt"
	"github.com/clipperhouse/gen/inflect"
	"strings"
	"time"
)

type genSpec struct {
	Pointer     string
	Package     string
	Singular    string
	Methods     []string
	Projections []*projection
	Imports     []string
	Plural      string
	Receiver    string
	Loop        string
	Generated   string
	Command     string
	FileName    string
}

func newGenSpec(t *Type) *genSpec {
	plural := inflect.Pluralize(t.Name)
	if plural == t.Name {
		plural += "s"
	}
	return &genSpec{
		Pointer:   t.Pointer,
		Package:   t.Package,
		Singular:  t.Name,
		Plural:    plural,
		Receiver:  "rcv",
		Loop:      "v",
		Generated: time.Now().UTC().Format(time.RFC1123),
		Command:   fmt.Sprintf("%s %s", "gen", t),
		FileName:  strings.ToLower(t.Name) + "_gen.go",
	}
}

func (g genSpec) Type() string {
	return g.Pointer + g.Package + "." + g.Singular
}

func (g genSpec) String() string {
	return g.Plural
}
