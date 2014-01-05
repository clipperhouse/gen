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

func newGenSpec(ptr, pkg, name string) *genSpec {
	plural := inflect.Pluralize(name)
	if plural == name {
		plural += "s"
	}
	return &genSpec{
		Pointer:   ptr,
		Package:   pkg,
		Singular:  name,
		Plural:    inflect.Pluralize(name),
		Receiver:  "rcv",
		Loop:      "v",
		Generated: time.Now().UTC().Format(time.RFC1123),
		Command:   fmt.Sprintf("%s %s%s.%s", "gen", ptr, pkg, name),
		FileName:  strings.ToLower(name) + "_gen.go",
	}
}

func (g genSpec) Type() string {
	return g.Pointer + g.Package + "." + g.Singular
}

func (g genSpec) String() string {
	return g.Plural
}
