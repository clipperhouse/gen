package slice

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/clipperhouse/gen/typewriter"
)

func init() {
	err := typewriter.Register(NewSliceWriter())
	if err != nil {
		panic(err)
	}
}

type cache struct {
	values []typewriter.TagValue
}

func Plural(typ typewriter.Type) string {
	return typ.Name + "Slice"
}

type SliceWriter struct {
	// map key is type name; type is not comparable, use typ.String()
	validated map[string]bool
	caches    map[string]cache
}

func NewSliceWriter() *SliceWriter {
	return &SliceWriter{
		validated: make(map[string]bool),
		caches:    make(map[string]cache),
	}
}

func (sw *SliceWriter) Name() string {
	return "slice"
}

func (sw *SliceWriter) Validate(typ typewriter.Type) (bool, error) {
	tag, found, err := typ.Tags.ByName("slice")

	if err != nil {
		return false, err
	}

	if !found {
		return false, nil
	}

	sw.validated[typ.String()] = true

	var values []typewriter.TagValue

	// filter methods applicable to type
	for _, v := range tag.Values {
		tmpl, ok := templates[v.Name]

		if !ok {
			err = fmt.Errorf("unknown slice method %s", v.Name)
			return false, err
		}

		if !tmpl.ApplicableToType(typ) {
			// TODO better error message
			err = fmt.Errorf("type %s cannot implement %s", typ, v)
			return false, err
		}

		if !tmpl.ApplicableToValue(v) {
			err = fmt.Errorf("type %s cannot implement %s; requires %d type parameters", typ, v, tmpl.RequiresTypeParameters)
			return false, err
		}

		values = append(values, v)
	}

	// store it for later, so we don't have to look for the tag again
	sw.caches[typ.String()] = cache{values}

	return true, nil
}

func (sw *SliceWriter) ensureValidation(typ typewriter.Type) {
	if !sw.validated[typ.String()] {
		err := fmt.Errorf("Type '%s' has not been previously validated. TypeWriter.Validate() must be called on all types before using them in subsequent methods.", typ)
		panic(err)
	}
}

func (sw *SliceWriter) WriteHeader(w io.Writer, typ typewriter.Type) {
	sw.ensureValidation(typ)

	cache, exists := sw.caches[typ.String()]

	if !exists {
		return
	}

	s := `// See http://clipperhouse.github.io/gen for documentation

`
	w.Write([]byte(s))

	if includeSortSupport(cache.values) {
		s := `// Sort implementation is a modification of http://golang.org/pkg/sort/#Sort
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found at http://golang.org/LICENSE.

`
		w.Write([]byte(s))
	}
}

func (sw *SliceWriter) Imports(typ typewriter.Type) []typewriter.ImportSpec {
	sw.ensureValidation(typ)

	var result []typewriter.ImportSpec

	methodRequiresErrors := map[string]bool{
		"First":   true,
		"Single":  true,
		"Max":     true,
		"Min":     true,
		"MaxBy":   true,
		"MinBy":   true,
		"Average": true,
	}

	methodRequiresSort := map[string]bool{
		"Sort": true,
	}

	cache, exists := sw.caches[typ.String()]

	if !exists {
		return result
	}

	imports := make(map[string]bool)

	for _, v := range cache.values {
		if methodRequiresErrors[v.Name] {
			imports["errors"] = true
		}

		if methodRequiresSort[v.Name] {
			imports["sort"] = true
		}
	}

	for imp := range imports {
		result = append(result, typewriter.ImportSpec{
			Path: imp,
		})
	}

	return result
}

func (sw *SliceWriter) WriteBody(w io.Writer, typ typewriter.Type) {
	sw.ensureValidation(typ)

	cache, exists := sw.caches[typ.String()]

	if !exists {
		return
	}

	tmpl, _ := templates.Get("slice")

	m := model{
		Type:   typ,
		Plural: Plural(typ),
	}

	if err := tmpl.Execute(w, m); err != nil {
		panic(err)
	}

	for _, v := range cache.values {
		var tp typewriter.Type

		if len(v.TypeParameters) > 0 {
			tp = v.TypeParameters[0]
		}

		m := model{
			Type:          typ,
			Plural:        Plural(typ),
			TypeParameter: tp,
			TagValue:      v,
		}

		tmpl, _ := templates.Get(v.Name) // already validated above
		err := tmpl.Execute(w, m)
		if err != nil {
			panic(err)
		}
	}
}

func includeSortSupport(values []typewriter.TagValue) bool {
	for _, v := range values {
		if strings.HasPrefix(v.Name, "SortBy") {
			return true
		}
	}
	return false
}

func includeSortInterface(values []typewriter.TagValue) bool {
	reg := regexp.MustCompile(`^Sort(Desc)?$`)
	for _, v := range values {
		if reg.MatchString(v.Name) {
			return true
		}
	}
	return false
}
