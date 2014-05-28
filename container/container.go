package container

import (
	"github.com/clipperhouse/typewriter"
	"io"
)

func init() {
	err := typewriter.Register(NewContainerWriter())
	if err != nil {
		panic(err)
	}
}

type ContainerWriter struct {
	tagsByType map[string]typewriter.Tag // typewriter.Type is not comparable, key by .String()
}

func NewContainerWriter() ContainerWriter {
	return ContainerWriter{
		tagsByType: make(map[string]typewriter.Tag),
	}
}

func (c ContainerWriter) Name() string {
	return "container"
}

func (c ContainerWriter) Validate(t typewriter.Type) (bool, error) {
	tag, found, err := t.Tags.ByName("containers")
	if found && err == nil {
		c.tagsByType[t.String()] = tag
	}
	return found, err
}

func (c ContainerWriter) WriteHeader(w io.Writer, t typewriter.Type) {
	// TODO: add license
	return
}

func (c ContainerWriter) Imports(t typewriter.Type) (result []string) {
	return result
}

func (c ContainerWriter) Write(w io.Writer, t typewriter.Type) {
	tag := c.tagsByType[t.String()] // validated above

	for _, s := range tag.Items {
		tmpl, err1 := containerTemplates.Get(s) // validate above to avoid err check here?
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
