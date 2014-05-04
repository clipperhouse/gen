package container

import (
	"gen/typewriter"
	"io"
)

func init() {
	typewriter.Register(ContainerWriter{})
}

type ContainerWriter struct{}

func (s ContainerWriter) Name() string {
	return "container"
}

func (s ContainerWriter) Validate(t typewriter.Type) (bool, error) {
	return false, nil
}

func (s ContainerWriter) WriteHeader(w io.Writer, t typewriter.Type) {
	return
}

func (s ContainerWriter) Imports(t typewriter.Type) (result []string) {
	return result
}

func (s ContainerWriter) Write(w io.Writer, t typewriter.Type) {
	return
}
