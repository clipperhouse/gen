package container

import (
	"github.com/clipperhouse/typewriter"
	"io"
)

func init() {
	err := typewriter.Register(ContainerWriter{})
	if err != nil {
		panic(err)
	}
}

type ContainerWriter struct{}

func (s ContainerWriter) Name() string {
	return "container"
}

var tagsByType = make(map[string]typewriter.Tag)

func (s ContainerWriter) Validate(t typewriter.Type) (bool, error) {
	tag, found, err := t.Tags.ByName("containers")
	if found && err == nil {
		tagsByType[t.String()] = tag
	}
	return found, err
}

func (s ContainerWriter) WriteHeader(w io.Writer, t typewriter.Type) {
	// TODO: add license
	return
}

func (s ContainerWriter) Imports(t typewriter.Type) (result []string) {
	return result
}

func (s ContainerWriter) Write(w io.Writer, t typewriter.Type) {
	tag := tagsByType[t.String()] // validated above

	for _, c := range tag.Items {
		tmpl, err1 := containerTemplates.Get(c) // validate above to avoid err check here?
		if err1 != nil {
			panic(err1)
		}

		err2 := tmpl.Execute(w, t)
		if err2 != nil {
			panic(err2)
		}
	}

	return
}
