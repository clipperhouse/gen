package main

import (
	"text/template"
)

type Template struct {
	Text               string
	RequiresNumeric    bool
	RequiresComparable bool
	RequiresOrdered    bool
}

func getTemplate(name string) (result *template.Template, err error) {
	if isProjectionMethod(name) {
		return getProjectionTemplate(name)
	}
	return getStandardTemplate(name)
}

func getHeaderTemplate() *template.Template {
	return template.Must(template.New("header").Parse(header))
}

const header = `// This file was auto-generated using github.com/clipperhouse/gen
// Modifying this file is not recommended as it will likely be overwritten in the future

// Sort (if included below) is a modification of http://golang.org/pkg/sort/#Sort
// List (if included below) is a modification of http://golang.org/pkg/container/list/
// Ring (if included below) is a modification of http://golang.org/pkg/container/ring/
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Set (if included below) is a modification of https://github.com/deckarep/golang-set
// The MIT License (MIT)
// Copyright (c) 2013 Ralph Caraveo (deckarep@gmail.com)

package {{.Package.Name}}
{{if gt (len .Imports) 0}}
import ({{range .Imports}}
	"{{.}}"{{end}}
)
{{end}}
// {{.Plural}} is a slice of type {{.Pointer}}{{.Name}}, for use with gen methods below. Use this type where you would use []{{.Pointer}}{{.Name}}. (This is required because slices cannot be method receivers.)
type {{.Plural}} []{{.Pointer}}{{.Name}}
`
