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

func SliceName(typ typewriter.Type) string {
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

	for _, v := range tag.Values {
		// just a validation here, template is used in Write()
		_, err := templates.Get(v)

		if err != nil {
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

	// a bit frustrating that errors are still possible here (templates mainly),
	// but all we can do is panic

	tmpl, err := templates.ByName("slice")

	if err != nil {
		panic(err)
	}

	m := model{
		Type:      typ,
		SliceName: SliceName(typ),
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
			SliceName:     SliceName(typ),
			TypeParameter: tp,
			TagValue:      v,
		}

		tmpl, err := templates.Get(v) // already validated above

		if err != nil {
			panic(err)
		}

		if err := tmpl.Execute(w, m); err != nil {
			panic(err)
		}
	}

	if includeSortInterface(cache.values) {
		tmpl, _ := templates.ByName("sortInterface") // already validated above
		err := tmpl.Execute(w, m)
		if err != nil {
			panic(err)
		}
	}

	if includeSortSupport(cache.values) {
		tmpl, _ := templates.ByName("sortImplementation") // already validated above
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
