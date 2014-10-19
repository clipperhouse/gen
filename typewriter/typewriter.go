package typewriter

import (
	"io"
)

type TypeWriter interface {
	Name() string
	// WriteHeader writer to the top of the generated code, before the package declaration; intended for licenses or general documentation.
	WriteHeader(w io.Writer, t Type) error
	// Imports is a slice of imports required for the type; each will be written into the imports declaration.
	Imports(t Type) []ImportSpec
	// WriteBody writes to the body of the generated code, following package declaration, headers and imports. This is the meat.
	WriteBody(w io.Writer, t Type) error
}
