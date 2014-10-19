package slice

import (
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

type SliceWriter struct{}

func NewSliceWriter() *SliceWriter {
	return &SliceWriter{}
}

func (sw *SliceWriter) Name() string {
	return "slice"
}

func (sw *SliceWriter) WriteHeader(w io.Writer, typ typewriter.Type) error {
	tag, found, err := typ.Tags.ByName("slice")

	if err != nil {
		return err
	}

	if !found {
		return nil
	}

	s := `// See http://clipperhouse.github.io/gen for documentation

`
	w.Write([]byte(s))

	if includeSortSupport(tag.Values) {
		s := `// Sort implementation is a modification of http://golang.org/pkg/sort/#Sort
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found at http://golang.org/LICENSE.

`
		w.Write([]byte(s))
	}

	return nil
}

func (sw *SliceWriter) Imports(typ typewriter.Type) (result []typewriter.ImportSpec) {
	tag, found, err := typ.Tags.ByName("slice")

	if err != nil {
		return
	}

	if !found {
		return
	}

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

	imports := make(map[string]bool)

	for _, v := range tag.Values {
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

func (sw *SliceWriter) WriteBody(w io.Writer, typ typewriter.Type) error {
	tag, found, err := typ.Tags.ByName("slice")

	if err != nil {
		return err
	}

	if !found {
		return nil
	}

	tmpl, err := templates.ByName("slice")

	if err != nil {
		return err
	}

	m := model{
		Type:      typ,
		SliceName: SliceName(typ),
	}

	if err := tmpl.Execute(w, m); err != nil {
		return err
	}

	for _, v := range tag.Values {
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
			return err
		}

		if err := tmpl.Execute(w, m); err != nil {
			return err
		}
	}

	if includeSortInterface(tag.Values) {
		tmpl, _ := templates.ByName("sortInterface") // already validated above
		err := tmpl.Execute(w, m)
		if err != nil {
			return err
		}
	}

	if includeSortSupport(tag.Values) {
		tmpl, _ := templates.ByName("sortImplementation") // already validated above
		err := tmpl.Execute(w, m)
		if err != nil {
			return err
		}
	}

	return nil
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
